package service

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type Middleware interface {
	CreateContext(w http.ResponseWriter, r *http.Request) (common.Context, error)
	GetContext(userID uint) common.Context
}

type JsonWebTokenService interface {
	ParseToken(r *http.Request) (*jwt.Token, *JsonWebTokenClaims, error)
	GenerateToken(w http.ResponseWriter, req *http.Request)
	Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Middleware
}

type WalletService interface {
	GetBalances() float64
}

type PriceHistoryService interface {
	GetPriceHistory(currency string) []dto.PriceHistoryDTO
}

type PortfolioService interface {
	Build(user common.UserContext, currencyPair *common.CurrencyPair) (common.Portfolio, error)
	Queue(user common.UserContext) (<-chan common.Portfolio, error)
	Stream(user common.UserContext, currencyPair *common.CurrencyPair) (<-chan common.Portfolio, error)
	Stop(user common.UserContext)
	IsStreaming(user common.UserContext) bool
}

type UserService interface {
	CreateUser(user common.UserContext)
	GetCurrentUser() (common.UserContext, error)
	GetUserById(userId uint) (common.UserContext, error)
	GetUserByName(username string) (common.UserContext, error)
	GetExchange(user common.UserContext, name string, currencyPair *common.CurrencyPair) (common.Exchange, error)
	GetConfiguredExchanges() []common.UserCryptoExchange
	GetExchangeSummary(currencyPair *common.CurrencyPair) ([]common.CryptoExchangeSummary, error)
	GetWallets(user common.UserContext) []common.UserCryptoWallet
	GetWallet(user common.UserContext, currency string) common.UserCryptoWallet
	GetTokens(user common.UserContext, wallet string) ([]common.EthereumToken, error)
	GetAllTokens(user common.UserContext) ([]common.EthereumToken, error)
}

type AuthService interface {
	Login(username, password string) (common.UserContext, error)
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
	CreateExchange(exchangeName string) (common.Exchange, error)
	GetDisplayNames() []string
	GetExchanges() []common.Exchange
	GetExchange(name string) common.Exchange
	GetCurrencyPairs(exchangeName string) ([]common.CurrencyPair, error)
}

type OrderService interface {
	GetMapper() mapper.OrderMapper
	GetOrderHistory() []common.Order
	ImportCSV(file, exchange string) ([]common.Order, error)
}
