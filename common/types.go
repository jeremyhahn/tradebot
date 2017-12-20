package common

import (
	"time"
)

type PriceChannel struct {
	Currency string
	Satoshis float64
	Price    float64
}

type ChartData struct {
	Currency          string  `json:"currency"`
	Price             float64 `json:"price"`
	Satoshis          float64 `json:"satoshis"`
	MACDValue         float64 `json:"macd_value"`
	MACDHistogram     float64 `json:"macd_histogram"`
	MACDSignal        float64 `json:"macd_signal"`
	MACDValueLive     float64 `json:"macd_value_live"`
	MACDHistogramLive float64 `json:"macd_histogram_live"`
	MACDSignalLive    float64 `json:"macd_signal_live"`
	RSI               float64 `json:"rsi"`
	RSILive           float64 `json:"rsi_live"`
	BollingerUpper    float64 `json:"bband_upper"`
	BollingerMiddle   float64 `json:"bband_middle"`
	BollingerLower    float64 `json:"bband_lower"`
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
	OnPriceChange(price float64)
}

type PeriodListener interface {
	OnPeriodChange(candlestick *Candlestick)
}

type Exchange interface {
	SubscribeToLiveFeed(price chan PriceChannel)
	GetPrice() float64
	GetSatoshis() float64
	GetTradeHistory(start, end time.Time, granularity int) []Candlestick
	GetCurrency() string
}

type Indicator interface {
	Calculate(price float64)
	PeriodListener
}

type Account struct {
	Currency float64
	Balance  float64
}

type Trade struct {
	ID        int `gorm:"primary_key"`
	Timestamp int32
	Price     float64
	Size      float64
}
