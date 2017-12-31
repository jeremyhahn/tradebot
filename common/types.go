package common

import (
	"time"
)

const (
	APPNAME     = "tradebot"
	APPVERSION  = "0.0.1"
	TIME_FORMAT = time.RFC3339
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

type Portfolio struct {
	User      *User
	NetWorth  float64        `json:"netWorth"`
	Exchanges []CoinExchange `json:"exchanges"`
}

type CoinExchange struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Total    float64 `json:"total"`
	Satoshis float64 `json:"satoshis"`
	Coins    []Coin  `json:"coins"`
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
	GetExchange() (CoinExchange, float64)
	GetNetWorth() float64
}

type Indicator interface {
	Calculate(price float64)
	PeriodListener
}
