package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type AutoTradeService interface {
	EndWorldHunger() error
}

type ChartService interface {
	GetCurrencyPair(chart *common.Chart) *common.CurrencyPair
	GetExchange(chart *common.Chart) common.Exchange
	Stream(chart *common.Chart, strategyHandler func(price float64) error) error
	StopStream(chart *common.Chart)
	GetCharts() ([]common.Chart, error)
	GetTrades(chart *common.Chart) ([]common.Trade, error)
	GetLastTrade(chart common.Chart) (*common.Trade, error)
	GetChart(id uint) (*common.Chart, error)
	GetIndicator(chart *common.Chart, name string) (common.FinancialIndicator, error)
	GetIndicators(chart *common.Chart) (map[string]common.FinancialIndicator, error)
	CreateIndicator(dao *dao.ChartIndicator) common.FinancialIndicator
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
