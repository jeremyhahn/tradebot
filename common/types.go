package common

import (
	"time"

	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
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

type Profit struct {
	UserId   uint    `json:"id"`
	TradeId  uint    `json:"trade_id"`
	Quantity float64 `json:"quantity"`
	Bought   float64 `json:"bought"`
	Sold     float64 `json:"sold"`
	Fee      float64 `json:"fee"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
}

type TradingStrategy interface {
	GetRequiredIndicators() []string
	GetBuySellSignals() (bool, bool, error)
	CalculateFeeAndTax(price float64) (float64, float64)
	GetTradeAmounts() (float64, float64)
}

type Trade struct {
	Id        uint      `json:"id"`
	ChartId   uint      `json:"chart_id"`
	UserId    uint      `json:"user_id"`
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	Exchange  string    `json:"exchange"`
	Date      time.Time `json:"date"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	ChartData string    `json:"chart_data"`
}

type Order struct {
	Id       string    `json:"id"`
	Exchange string    `json:"exchange"`
	Date     time.Time `json:"date"`
	Type     string    `json:"type"`
	Currency string    `json:"currency"`
	Quantity float64   `json:"quantity"`
	Price    float64   `json:"price"`
}

type Context struct {
	Logger *logging.Logger
	DB     *gorm.DB
	User   *User
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

type GlobalMarketCap struct {
	TotalMarketCapUSD float64 `json:"total_market_cap_usd"`
	Total24HVolumeUSD float64 `json:"total_24h_volume_usd"`
	BitcoinDominance  float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies  float64 `json:"active_currencies"`
	ActiveMarkets     float64 `json:"active_markets"`
	LastUpdated       int64   `json:"last_updated"`
}

type Wallet interface {
	GetBalance() CryptoWallet
}

type CryptoWallet struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	NetWorth float64 `json:"net_worth"`
}

type User struct {
	Id            uint   `json:"id"`
	Username      string `json:"username"`
	LocalCurrency string `json:"local_currency"`
}

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

type Portfolio struct {
	User      *User            `json:"user"`
	NetWorth  float64          `json:"net_worth"`
	Exchanges []CryptoExchange `json:"exchanges"`
	Wallets   []CryptoWallet   `json:"wallets"`
}

type PriceChange struct {
	Exchange     string        `json:"exchange"`
	CurrencyPair *CurrencyPair `json:"currencyPair"`
	Satoshis     float64       `json:"satoshis"`
	Price        float64       `json:"price"`
}

type Chart struct {
	Id         uint        `json:"id"`
	Base       string      `json:"base"`
	Quote      string      `json:"quote"`
	Exchange   string      `json:"exchange"`
	Period     int         `json:"period"`
	Price      float64     `json:"price"`
	AutoTrade  uint        `json:"autotrade"`
	Indicators []Indicator `json:"indicators"`
	Trades     []Trade     `json:"trades"`
}

type CryptoExchange struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Total    float64 `json:"total"`
	Satoshis float64 `json:"satoshis"`
	Coins    []Coin  `json:"coins"`
}

type CryptoExchangeList struct {
	Exchanges []CryptoExchange `json:"exchange"`
	NetWorth  float64          `json:"net_worth"`
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

type ChartData struct {
	CurrencyPair        CurrencyPair `json:"currency"`
	Exchange            string       `json:"exchange"`
	Price               float64      `json:"price"`
	Satoshis            float64      `json:"satoshis"`
	MACDValue           float64      `json:"macd_value"`
	MACDHistogram       float64      `json:"macd_histogram"`
	MACDSignal          float64      `json:"macd_signal"`
	MACDValueLive       float64      `json:"macd_value_live"`
	MACDHistogramLive   float64      `json:"macd_histogram_live"`
	MACDSignalLive      float64      `json:"macd_signal_live"`
	RSI                 float64      `json:"rsi"`
	RSILive             float64      `json:"rsi_live"`
	BollingerUpper      float64      `json:"bband_upper"`
	BollingerMiddle     float64      `json:"bband_middle"`
	BollingerLower      float64      `json:"bband_lower"`
	BollingerUpperLive  float64      `json:"bband_upper_live"`
	BollingerMiddleLive float64      `json:"bband_middle_live"`
	BollingerLowerLive  float64      `json:"bband_lower_live"`
	OnBalanceVolume     float64      `json:"on_balance_volume"`
	OnBalanceVolumeLive float64      `json:"on_balance_volume_live"`
}

type PriceListener interface {
	OnPriceChange(priceChange *PriceChange)
}

type PeriodListener interface {
	OnPeriodChange(candlestick *Candlestick)
}

type FinancialIndicator interface {
	GetDefaultParameters() []string
	GetParameters() []string
	GetDisplayName() string
	GetName() string
	PeriodListener
}

type Indicator struct {
	Id         uint   `json:"id"`
	ChartId    uint   `json:"chart_id"`
	Name       string `json:"name"`
	Parameters string `json:"parameters"`
	FinancialIndicator
}
