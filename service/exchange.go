package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
	logging "github.com/op/go-logging"
)

type ExchangeService struct {
	ctx *common.Context
	dao *dao.ExchangeDAO
}

func NewExchangeService(ctx *common.Context, dao *dao.ExchangeDAO) *ExchangeService {
	return &ExchangeService{
		ctx: ctx,
		dao: dao}
}

func (service *ExchangeService) NewExchange(user *common.User, exchangeName string, currencyPair *common.CurrencyPair) common.Exchange {
	exchangeMap := map[string]func(*dao.UserCoinExchange, *logging.Logger, *common.CurrencyPair) common.Exchange{
		"gdax":    exchange.NewGDAX,
		"bittrex": exchange.NewBittrex,
		"binance": exchange.NewBinance}
	userDAO := dao.NewUserDAO(service.ctx)
	ex := userDAO.GetExchange(service.ctx.User, exchangeName)
	return exchangeMap[exchangeName](ex, service.ctx.Logger, currencyPair)
}
