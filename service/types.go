package service

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/shopspring/decimal"
)

type Middleware interface {
	CreateContext(w http.ResponseWriter, r *http.Request) (common.Context, error)
	GetContext(userID uint) common.Context
}

type JsonWebTokenService interface {
	ParseToken(r *http.Request, extractor request.Extractor) (*jwt.Token, *JsonWebTokenClaims, error)
	GenerateToken(w http.ResponseWriter, req *http.Request)
	Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Middleware
}

type AuthService interface {
	Login(username, password string) (common.UserContext, error)
	Register(username, password string) error
}

type TokenService interface {
	GetToken(walletAddress, contractAddress string) (common.EthereumToken, error)
	GetTokenTransactions(contractAddress string) ([]common.Transaction, error)
	//GetContract(address string) (common.EthereumContract, error)
	//GetTotalSupply() uint
	//GetAllowance() uint
	//Transfer() bool
	//Approve() bool
	//TransferFrom(from, to EthereumTokenAddress, tokens, uint)
}

type EthereumService interface {
	GetPrice() decimal.Decimal
	GetWallet(address string) (common.UserCryptoWallet, error)
	GetTransactions() ([]common.Transaction, error)
	GetTransactionsFor(address string) ([]common.Transaction, error)
	GetAccounts() ([]common.UserContext, error)
	AuthService
	TokenService
}

type GethService interface {
	Authenticate(address, passphrase string) error
	CreateAccount(passphrase string) (common.UserContext, error)
	DeleteAccount(passphrase string) error
	AuthService
	TokenService
}

type WalletService interface {
	CreateWallet(currency, address string) (common.Wallet, error)
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
	GetWallet(currency string) (common.UserCryptoWallet, error)
	GetWallets() []common.UserCryptoWallet
	GetWalletPlugins() ([]common.Wallet, error)
	GetTokensFor(wallet string) ([]common.EthereumToken, error)
	GetTokens() ([]common.EthereumToken, error)
	CreateToken(token common.EthereumToken) error
	CreateWallet(wallet common.UserCryptoWallet) error
	CreateExchange(userCryptoExchange common.UserCryptoExchange) (common.UserCryptoExchange, error)
	DeleteExchange(exchangeName string) error
}

type AutoTradeService interface {
	EndWorldHunger() error
}

type ChartService interface {
	GetCurrencyPair(chart common.Chart) *common.CurrencyPair
	GetExchange(chart common.Chart) (common.Exchange, error)
	Stream(chart common.Chart, candlesticks []common.Candlestick, strategyHandler func(price decimal.Decimal) error) error
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
	GetMapper() mapper.TradeMapper
	Save(dto common.Trade)
	GetLastTrade(chart common.Chart) common.Trade
	GetTradeHistory() []common.Transaction
	GetTransactionMapper() mapper.TransactionMapper
}

type ProfitService interface {
	Save(profit common.Profit)
	Find()
}

type ExchangeService interface {
	CreateExchange(exchangeName string) (common.Exchange, error)
	GetDisplayNames() ([]string, error)
	GetExchanges() ([]common.Exchange, error)
	GetExchange(name string) (common.Exchange, error)
	GetCurrencyPairs(exchangeName string) ([]common.CurrencyPair, error)
}

type TransactionService interface {
	GetMapper() mapper.TransactionMapper
	GetHistory(order string) ([]common.Transaction, error)
	GetOrderHistory() []common.Transaction
	GetDepositHistory() []common.Transaction
	GetWithdrawalHistory() []common.Transaction
	GetImportedTransactions() []common.Transaction
	UpdateCategory(id, category string) error
	ImportCSV(file, exchange string) ([]common.Transaction, error)
	Synchronize() ([]common.Transaction, error)
	//GetSourceTransaction(targetTx common.Transaction, transactions *[]common.Transaction) (common.Transaction, error)
}
