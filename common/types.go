package common

import (
	"net/http"
	"time"
)

const (
	APPNAME               = "tradebot"
	APPVERSION            = "0.0.1"
	TIME_FORMAT           = time.RFC3339
	BUFFERED_CHANNEL_SIZE = 256
	WEBSOCKET_KEEPALIVE   = 10 * time.Second
	HTTP_CLIENT_TIMEOUT   = 10 * time.Second
	CANDLESTICK_MIN_LOAD  = 250
)

type Indicator interface {
	GetName() string
	GetFilename() string
	GetVersion() string
}

type Strategy interface {
	GetName() string
	GetFilename() string
	GetVersion() string
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
	GetPrice() float64
	GetAmount() float64
	GetChartData() string
}

type Order interface {
	GetId() string
	GetExchange() string
	GetDate() time.Time
	GetType() string
	GetCurrency() string
	GetQuantity() float64
	GetPrice() float64
}

type Profit interface {
	GetUserId() uint
	GetTradeId() uint
	GetQuantity() float64
	GetBought() float64
	GetSold() float64
	GetFee() float64
	GetTax() float64
	GetTotal() float64
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
	CalculateFeeAndTax(price float64) (float64, float64)
	GetTradeAmounts() (float64, float64)
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

type Exchange interface {
	GetName() string
	GetBalances() ([]Coin, float64)
	GetExchange() CryptoExchange
	GetNetWorth() float64
	GetTradingFee() float64
	SubscribeToLiveFeed(currencyPair *CurrencyPair, price chan PriceChange)
	GetPrice(currencyPair *CurrencyPair) float64
	GetPriceHistory(currencyPair *CurrencyPair, start, end time.Time, granularity int) []Candlestick
	GetOrderHistory(currencyPair *CurrencyPair) []Order
	FormattedCurrencyPair(currencyPair *CurrencyPair) string
}

type User interface {
	GetId() uint
	GetUsername() string
	GetLocalCurrency() string
	GetEtherbase() string
	GetKeystore() string
	SetEtherbase(etherbase string)
	SetKeystore(keystore string)
}

type Wallet interface {
	GetBalance() CryptoWallet
}

type UserWallet interface {
	GetUserId() uint
	GetCurrency() string
	GetAddress() string
}

type Coin interface {
	GetCurrency() string
	GetBalance() float64
	GetAvailable() float64
	GetPending() float64
	GetPrice() float64
	GetAddress() string
	GetTotal() float64
	GetBTC() float64
	IsBitcoin() bool
}

type CryptoExchange interface {
	GetName() string
	GetURL() string
	GetTotal() float64
	GetSatoshis() float64
	GetCoins() []Coin
}

type CryptoExchangeList interface {
	GetExchanges() []CryptoExchange
	GetNetWorth() float64
}

type CryptoWallet interface {
	GetAddress() string
	GetBalance() float64
	GetCurrency() string
	GetNetWorth() float64
}

type Portfolio interface {
	GetUser() User
	GetNetWorth() float64
	GetExchanges() []CryptoExchange
	GetWallets() []CryptoWallet
}

type HttpWriter interface {
	Write(w http.ResponseWriter, status int, response interface{})
}

type TradingStrategyParams struct {
	CurrencyPair *CurrencyPair
	Balances     []Coin
	Indicators   map[string]FinancialIndicator
	NewPrice     float64
	LastTrade    Trade
	TradeFee     float64
	Config       []string
}

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

type PriceChange struct {
	Exchange     string        `json:"exchange"`
	CurrencyPair *CurrencyPair `json:"currencyPair"`
	Satoshis     float64       `json:"satoshis"`
	Price        float64       `json:"price"`
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
