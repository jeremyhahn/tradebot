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
)

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
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

type Portfolio struct {
	User      *User          `json:"user"`
	NetWorth  float64        `json:"net_worth"`
	Exchanges []CoinExchange `json:"exchanges"`
	Wallets   []CryptoWallet `json:"wallets"`
}

type CoinExchange struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Total    float64 `json:"total"`
	Satoshis float64 `json:"satoshis"`
	Coins    []Coin  `json:"coins"`
}

type CoinExchangeList struct {
	Exchanges []CoinExchange `json:"exchange"`
	NetWorth  float64        `json:"net_worth"`
}

type PriceChange struct {
	Exchange     string        `json:"exchange"`
	CurrencyPair *CurrencyPair `json:"currencyPair"`
	Satoshis     float64       `json:"satoshis"`
	Price        float64       `json:"price"`
}

type ChartData struct {
	CurrencyPair      CurrencyPair `json:"currency"`
	Exchange          string       `json:"exchange"`
	Price             float64      `json:"price"`
	Satoshis          float64      `json:"satoshis"`
	MACDValue         float64      `json:"macd_value"`
	MACDHistogram     float64      `json:"macd_histogram"`
	MACDSignal        float64      `json:"macd_signal"`
	MACDValueLive     float64      `json:"macd_value_live"`
	MACDHistogramLive float64      `json:"macd_histogram_live"`
	MACDSignalLive    float64      `json:"macd_signal_live"`
	RSI               float64      `json:"rsi"`
	RSILive           float64      `json:"rsi_live"`
	BollingerUpper    float64      `json:"bband_upper"`
	BollingerMiddle   float64      `json:"bband_middle"`
	BollingerLower    float64      `json:"bband_lower"`
}

type MovingAverage interface {
	Add(candle *Candlestick) float64
	GetCandlesticks() []Candlestick
	GetSize() int
	GetCount() int
	GetIndex() int
	GetAverage() float64
	Sum() float64
	GetGainsAndLosses() (float64, float64)
	PeriodListener
}

type PriceListener interface {
	OnPriceChange(priceChange *PriceChange)
}

type PeriodListener interface {
	OnPeriodChange(candlestick *Candlestick)
}

type Exchange interface {
	SubscribeToLiveFeed(price chan PriceChange)
	GetPriceUSD() float64
	GetPrice() float64
	GetSatoshis() float64
	GetTradeHistory(start, end time.Time, granularity int) []Candlestick
	GetCurrencyPair() CurrencyPair
	FormattedCurrencyPair() string
	GetBalances() ([]Coin, float64)
	GetName() string
	GetExchangeAsync(*chan CoinExchange)
	GetExchange() CoinExchange
	GetNetWorth() float64
}

type Indicator interface {
	Calculate(price float64)
	PeriodListener
}
