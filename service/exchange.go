package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/exchange"
)

type DefaultExchangeService struct {
	ctx         *common.Context
	dao         dao.ExchangeDAO
	exchangeMap map[string]func(*common.Context, entity.UserExchangeEntity) common.Exchange
	ExchangeService
}

func NewExchangeService(ctx *common.Context, exchangeDAO dao.ExchangeDAO) ExchangeService {
	return &DefaultExchangeService{
		ctx: ctx,
		dao: exchangeDAO,
		exchangeMap: map[string]func(ctx *common.Context, exchange entity.UserExchangeEntity) common.Exchange{
			"gdax":    exchange.NewGDAX,
			"bittrex": exchange.NewBittrex,
			"binance": exchange.NewBinance}}
}

func (service *DefaultExchangeService) CreateExchange(user common.User, exchangeName string) common.Exchange {
	userDAO := dao.NewUserDAO(service.ctx)
	userEntity := &entity.User{Id: user.GetId()}
	exchange := userDAO.GetExchange(userEntity, exchangeName)
	return service.exchangeMap[exchangeName](service.ctx, exchange)
}

func (service *DefaultExchangeService) GetExchanges(user common.User) []common.Exchange {
	var exchanges []common.Exchange
	userDAO := dao.NewUserDAO(service.ctx)
	userEntity := &entity.User{Id: user.GetId()}
	userExchanges := userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, service.exchangeMap[ex.Name](service.ctx, &ex))
	}
	return exchanges
}

func (service *DefaultExchangeService) GetExchange(user common.User, name string) common.Exchange {
	for _, ex := range service.GetExchanges(user) {
		if ex.GetName() == name {
			return ex
		}
	}
	return nil
}
