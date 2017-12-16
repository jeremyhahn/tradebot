package common

import (
	"time"
)

type Chart struct {
	Currency        string  `json:"currency"`
	Price           float64 `json:"Price"`
	MACDValue       float64 `json:"macd_value"`
	MACDHistogram   float64 `json:"macd_histogram"`
	MACDSignal      float64 `json:"macd_signal"`
	RSI             float64 `json:"rsi"`
	BollingerUpper  float64 `json:"bband_upper"`
	BollingerMiddle float64 `json:"bband_middle"`
	BollingerLower  float64 `json:"bband_lower"`
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
	OnPeriodChange(*Candlestick)
}

type Exchange interface {
	SubscribeToLiveFeed(price chan float64)
	GetPrice() float64
	GetTradeHistory(start, end time.Time, granularity int) []Candlestick
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
