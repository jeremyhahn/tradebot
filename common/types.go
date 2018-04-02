package common

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	APPNAME               = "tradebot"
	APPVERSION            = "0.0.1"
	TIME_FORMAT           = time.RFC3339
	TIME_DISPLAY_FORMAT   = "01-02-2006 15:04:05 MST"
	BUFFERED_CHANNEL_SIZE = 256
	WEBSOCKET_KEEPALIVE   = 10 * time.Second
	HTTP_CLIENT_TIMEOUT   = 10 * time.Second
	CANDLESTICK_MIN_LOAD  = 250
	INDICATOR_PLUGIN_TYPE = "indicator"
	STRATEGY_PLUGIN_TYPE  = "strategy"
	EXCHANGE_PLUGIN_TYPE  = "exchange"
	WALLET_PLUGIN_TYPE    = "wallet"
	BUY_ORDER_TYPE        = "buy"
	SELL_ORDER_TYPE       = "sell"
	DEPOSIT_ORDER_TYPE    = "deposit"
	WITHDRAWAL_ORDER_TYPE = "withdrawal"
)

type Transaction interface {
	GetId() string
	GetDate() time.Time
	GetCurrencyPair() *CurrencyPair
	GetType() string
	GetNetwork() string
	GetNetworkDisplayName() string
	GetQuantity() string
	GetQuantityCurrency() string
	GetFiatQuantity() string
	GetFiatQuantityCurrency() string
	GetPrice() string
	GetPriceCurrency() string
	GetPriceString() string
	GetFiatPrice() string
	GetFiatPriceCurrency() string
	GetFee() string
	GetFeeCurrency() string
	GetTotal() string
	GetTotalCurrency() string
	GetFiatFee() string
	GetFiatFeeCurrency() string
	GetFiatTotal() string
	GetFiatTotalCurrency() string
	String() string
}

type EthereumToken interface {
	GetName() string
	GetSymbol() string
	GetDecimals() uint8
	GetBalance() decimal.Decimal
	GetWalletAddress() string
	GetContractAddress() string
	GetValue() decimal.Decimal
}

type EthereumContract interface {
	GetAddress() string
	GetSource() string
	GetBin() string
	GetABI() string
	GetCreationDate() time.Time
}

type Plugin interface {
	GetName() string
	GetFilename() string
	GetVersion() string
	GetType() string
}

type Chart interface {
	GetId() uint
	GetBase() string
	GetQuote() string
	GetExchange() string
	GetPeriod() int
	GetPrice() float64
	GetAutoTrade() uint
	IsAutoTrade() bool
	GetIndicators() []ChartIndicator
	GetStrategies() []ChartStrategy
	GetTrades() []Trade
	ToJSON() (string, error)
}

type ChartIndicator interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
	GetFilename() string
}

type ChartStrategy interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
	GetFilename() string
}

type Trade interface {
	GetId() uint
	GetChartId() uint
	GetUserId() uint
	GetBase() string
	GetQuote() string
	GetExchange() string
	GetDate() time.Time
	GetType() string
	GetPrice() decimal.Decimal
	GetAmount() decimal.Decimal
	GetChartData() string
}

type Profit interface {
	GetUserId() uint
	GetTradeId() uint
	GetQuantity() decimal.Decimal
	GetBought() decimal.Decimal
	GetSold() decimal.Decimal
	GetFee() decimal.Decimal
	GetTax() decimal.Decimal
	GetTotal() decimal.Decimal
}

type FinancialIndicator interface {
	GetDefaultParameters() []string
	GetParameters() []string
	GetDisplayName() string
	GetName() string
	PeriodListener
}

type TradingStrategy interface {
	GetRequiredIndicators() []string
	Analyze() (bool, bool, map[string]string, error)
	CalculateFeeAndTax(price decimal.Decimal) (decimal.Decimal, decimal.Decimal)
	GetTradeAmounts() (decimal.Decimal, decimal.Decimal)
	GetParameters() *TradingStrategyParams
}

type ChartTradingStrategy interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
	GetDefaultParameters() string
	GetRequiredIndicators() string
}

type PriceListener interface {
	OnPriceChange(priceChange *PriceChange)
}

type PeriodListener interface {
	OnPeriodChange(candlestick *Candlestick)
}

type FiatPriceService interface {
	GetPriceAt(currency string, date time.Time) (*Candlestick, error)
}

type Exchange interface {
	GetName() string
	GetDisplayName() string
	GetBalances() ([]Coin, decimal.Decimal)
	GetSummary() CryptoExchangeSummary
	GetNetWorth() decimal.Decimal
	GetTradingFee() decimal.Decimal
	//GetCurrencies() []string
	SubscribeToLiveFeed(currencyPair *CurrencyPair, price chan PriceChange)
	GetPrice(currencyPair *CurrencyPair) decimal.Decimal
	GetPriceHistory(currencyPair *CurrencyPair, start, end time.Time, granularity int) ([]Candlestick, error)
	GetOrderHistory(currencyPair *CurrencyPair) []Transaction
	GetDepositHistory() ([]Transaction, error)
	GetWithdrawalHistory() ([]Transaction, error)
	GetCurrencies() (map[string]*Currency, error)
	FormattedCurrencyPair(currencyPair *CurrencyPair) string
	ParseImport(file string) ([]Transaction, error)
}

type KeyPair interface {
	GetDirectory() string
	GetPrivateKey() *rsa.PrivateKey
	GetPrivateBytes() []byte
	GetPublicKey() *rsa.PublicKey
	GetPublicBytes() []byte
}

type UserContext interface {
	GetId() uint
	GetUsername() string
	GetLocalCurrency() string
	GetFiatExchange() string
	GetEtherbase() string
	GetKeystore() string
}

type UserCryptoExchange interface {
	GetUserID() uint
	GetName() string
	GetKey() string
	GetURL() string
	GetSecret() string
	GetExtra() string
}

type Coin interface {
	GetCurrency() string
	GetPrice() decimal.Decimal
	GetExchange() string
	GetBalance() decimal.Decimal
	GetAvailable() decimal.Decimal
	GetPending() decimal.Decimal
	GetAddress() string
	GetTotal() decimal.Decimal
	GetBTC() decimal.Decimal
	GetUSD() decimal.Decimal
	IsBitcoin() bool
}

type CryptoExchangeSummary interface {
	GetName() string
	GetURL() string
	GetTotal() decimal.Decimal
	GetSatoshis() decimal.Decimal
	GetCoins() []Coin
}

type UserCryptoWallet interface {
	GetAddress() string
	GetBalance() decimal.Decimal
	GetCurrency() string
	GetValue() decimal.Decimal
}

type Portfolio interface {
	GetUser() UserContext
	GetNetWorth() decimal.Decimal
	GetExchanges() []CryptoExchangeSummary
	GetWallets() []UserCryptoWallet
	GetTokens() []EthereumToken
}

type HttpWriter interface {
	Write(w http.ResponseWriter, status int, response interface{})
}

type MarketCapService interface {
	GetMarkets() []MarketCap
	GetMarket(symbol string) MarketCap
	GetGlobalMarket(currency string) *GlobalMarketCap
	GetMarketsByPrice(order string) []MarketCap
	GetMarketsByPercentChange1H(order string) []MarketCap
	GetMarketsByPercentChange24H(order string) []MarketCap
	GetMarketsByPercentChange7D(order string) []MarketCap
	GetMarketsByTopPerformers(order string) []MarketCap
	GetTrendingMarkets(order string) []MarketCap
}

type Wallet interface {
	GetPrice() decimal.Decimal
	GetWallet() (UserCryptoWallet, error)
	GetTransactions() ([]Transaction, error)
}

type WalletParams struct {
	Context          Context
	Address          string
	WalletUser       string
	WalletSecret     string
	WalletExtra      string
	MarketCapService MarketCapService
	FiatPriceService FiatPriceService
}

type TradingStrategyParams struct {
	CurrencyPair *CurrencyPair
	Balances     []Coin
	Indicators   map[string]FinancialIndicator
	NewPrice     decimal.Decimal
	LastTrade    Trade
	TradeFee     decimal.Decimal
	Config       []string
}

type PriceChange struct {
	Exchange     string          `json:"exchange"`
	CurrencyPair *CurrencyPair   `json:"currencyPair"`
	Satoshis     decimal.Decimal `json:"satoshis"`
	Price        decimal.Decimal `json:"price"`
}

type ChartData struct {
	CurrencyPair CurrencyPair      `json:"currency"`
	Exchange     string            `json:"exchange"`
	Price        float64           `json:"price"`
	Satoshis     float64           `json:"satoshis"`
	Indicators   map[string]string `json:"indicators"`
}

type GlobalMarketCap struct {
	TotalMarketCapUSD float64 `json:"total_market_cap_usd"`
	Total24HVolumeUSD float64 `json:"total_24h_volume_usd"`
	BitcoinDominance  float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies  float64 `json:"active_currencies"`
	ActiveMarkets     float64 `json:"active_markets"`
	LastUpdated       int64   `json:"last_updated"`
}

type MarketCap struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUSD         string `json:"price_usd"`
	PriceBTC         string `json:"price_btc"`
	VolumeUSD24h     string `json:"24h_volume_usd"`
	MarketCapUSD     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}
