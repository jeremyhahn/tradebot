package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type PriceHistoryService interface {
	GetPriceHistory(currency string) []dto.PriceHistoryDTO
}

type PortfolioService interface {
	Build(user common.User, currencyPair *common.CurrencyPair) common.Portfolio
	Queue(user common.User) <-chan common.Portfolio
	Stream(user common.User, currencyPair *common.CurrencyPair) <-chan common.Portfolio
	Stop(user common.User)
	IsStreaming(user common.User) bool
}

type UserService interface {
	CreateUser(user common.User)
	GetCurrentUser() (common.User, error)
	GetUserById(userId uint) (common.User, error)
	GetUserByName(username string) (common.User, error)
	GetExchange(user common.User, name string, currencyPair *common.CurrencyPair) common.Exchange
	GetExchanges(user common.User, currencyPair *common.CurrencyPair) []common.CryptoExchange
	GetWallets(user common.User) []common.CryptoWallet
	GetWallet(user common.User, currency string) common.CryptoWallet
}

type AuthService interface {
	Login(username, password string) (common.User, error)
	Register(username, password string) error
}

type AutoTradeService interface {
	EndWorldHunger() error
}

type ChartService interface {
	GetCurrencyPair(chart common.Chart) *common.CurrencyPair
	GetExchange(chart common.Chart) common.Exchange
	Stream(chart common.Chart, candlesticks []common.Candlestick, strategyHandler func(price float64) error) error
	StopStream(chart common.Chart)
	GetChart(id uint) (common.Chart, error)
	GetCharts(autoTradeOnly bool) ([]common.Chart, error)
	GetTrades(chart common.Chart) ([]common.Trade, error)
	GetLastTrade(chart common.Chart) (common.Trade, error)
	GetIndicator(chart common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error)
	GetIndicators(chart common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error)
	CreateIndicator(dao entity.ChartIndicator) common.FinancialIndicator
	LoadCandlesticks(chart common.Chart, exchange common.Exchange) []common.Candlestick
}

type TradeService interface {
	Save(dto common.Trade)
	GetLastTrade(chart common.Chart) common.Trade
	GetMapper() mapper.TradeMapper
}

type ProfitService interface {
	Save(profit common.Profit)
	Find()
}

type ExchangeService interface {
	CreateExchange(user common.User, exchangeName string) common.Exchange
	GetDisplayNames(user common.User) []string
	GetExchanges(common.User) []common.Exchange
	GetExchange(user common.User, name string) common.Exchange
	GetCurrencyPairs(user common.User, exchangeName string) ([]common.CurrencyPair, error)
}

type OrderService interface {
	GetMapper() mapper.OrderMapper
	GetOrderHistory() []common.Order
	ImportCSV(file, exchange string) ([]common.Order, error)
}
