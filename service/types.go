package service

import (
	"github.com/jeremyhahn/tradebot/common"
)

type AutoTradeService interface {
	EndWorldHunger()
}

type ChartService interface {
	GetExchange(chart *common.Chart) common.Exchange
	//GetCurrencyPair(chart *common.Chart) common.CurrencyPair
	GetIndicator(*common.Chart, string) common.FinancialIndicator
	Stream(chart *common.Chart, strategyHandler func(price float64))
	StopStream(chart *common.Chart)
	ToJson(chart *common.Chart) string

	GetIndicators(chart *common.Chart) map[string]common.FinancialIndicator
	GetChart(id uint) *common.Chart
	GetCharts() []common.Chart
	GetChartsByUser(user *common.User) []common.Chart
}

type TradeService interface {
	Save(trade *common.Trade)
	GetLastTrade(chart *common.Chart) *common.Trade
}

type ProfitService interface {
	Save(profit *common.Profit)
	Find()
}

type ExchangeService interface {
	CreateExchange(user *common.User, exchangeName string) common.Exchange
	GetExchanges(*common.User) []common.Exchange
	GetExchange(user *common.User, name string) common.Exchange
}

type OrderService interface {
	GetOrderHistory() []common.Order
}
