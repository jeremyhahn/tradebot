package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
)

type ExchangeServiceImpl struct {
	ctx         *common.Context
	dao         dao.ExchangeDAO
	exchangeMap map[string]func(*common.Context, *dao.UserCryptoExchange) common.Exchange
	ExchangeService
}

func NewExchangeService(ctx *common.Context, exchangeDAO dao.ExchangeDAO) ExchangeService {
	return &ExchangeServiceImpl{
		ctx: ctx,
		dao: exchangeDAO,
		exchangeMap: map[string]func(ctx *common.Context, exchange *dao.UserCryptoExchange) common.Exchange{
			"gdax":    exchange.NewGDAX,
			"bittrex": exchange.NewBittrex,
			"binance": exchange.NewBinance}}
}

func (service *ExchangeServiceImpl) CreateExchange(user *common.User, exchangeName string) common.Exchange {
	userDAO := dao.NewUserDAO(service.ctx)
	exchange := userDAO.GetExchange(service.ctx.User, exchangeName)
	return service.exchangeMap[exchangeName](service.ctx, exchange)
}

func (service *ExchangeServiceImpl) GetExchanges(user *common.User) []common.Exchange {
	var exchanges []common.Exchange
	userDAO := dao.NewUserDAO(service.ctx)
	userExchanges := userDAO.GetExchanges(user)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, service.exchangeMap[ex.Name](service.ctx, &ex))
	}
	return exchanges
}

func (service *ExchangeServiceImpl) GetExchange(user *common.User, name string) common.Exchange {
	for _, ex := range service.GetExchanges(user) {
		if ex.GetName() == name {
			return ex
		}
	}
	return nil
}
