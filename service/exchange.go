package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
	logging "github.com/op/go-logging"
)

type ExchangeService interface {
	GetExchanges(*common.User, *common.CurrencyPair) []common.Exchange
	NewExchange(user *common.User, exchangeName string, currencyPair *common.CurrencyPair) common.Exchange
}

type ExchangeServiceImpl struct {
	ctx         *common.Context
	dao         dao.ExchangeDAO
	exchangeMap map[string]func(*dao.UserCoinExchange, *logging.Logger, *common.CurrencyPair) common.Exchange
}

func NewExchangeService(ctx *common.Context, exchangeDAO dao.ExchangeDAO) ExchangeService {
	return &ExchangeServiceImpl{
		ctx: ctx,
		dao: exchangeDAO,
		exchangeMap: map[string]func(*dao.UserCoinExchange, *logging.Logger, *common.CurrencyPair) common.Exchange{
			"gdax":    exchange.NewGDAX,
			"bittrex": exchange.NewBittrex,
			"binance": exchange.NewBinance}}
}

func (service *ExchangeServiceImpl) NewExchange(user *common.User, exchangeName string, currencyPair *common.CurrencyPair) common.Exchange {
	userDAO := dao.NewUserDAO(service.ctx)
	ex := userDAO.GetExchange(service.ctx.User, exchangeName)
	return service.exchangeMap[exchangeName](ex, service.ctx.Logger, currencyPair)
}

func (service *ExchangeServiceImpl) GetExchanges(user *common.User, currencyPair *common.CurrencyPair) []common.Exchange {
	var exchanges []common.Exchange
	userDAO := dao.NewUserDAO(service.ctx)
	userExchanges := userDAO.GetExchanges(user)
	for _, ex := range userExchanges {
		if strings.Contains(ex.Extra, ",") {
			symbols := strings.Split(ex.Extra, ",")
			for _, s := range symbols {
				baseQuote := strings.Split(s, "-")
				currencyPair = &common.CurrencyPair{
					Base:          baseQuote[0],
					Quote:         baseQuote[1],
					LocalCurrency: user.LocalCurrency}
				exchanges = append(exchanges, service.exchangeMap[ex.Name](&ex, service.ctx.Logger, currencyPair))
			}
		} else {
			exchanges = append(exchanges, service.exchangeMap[ex.Name](&ex, service.ctx.Logger, currencyPair))
		}
	}
	return exchanges
}
