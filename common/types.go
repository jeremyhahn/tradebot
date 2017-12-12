package common

import (
	"time"
)

type PriceListener interface {
	OnCandlestickCreated(*Candlestick)
	OnPriceChange(price float64)
}

type Exchange interface {
	SubscribeToLiveFeed(price chan float64)
	GetPrice() float64
	GetTradeHistory(start, end time.Time, granularity int) []Candlestick
}

type Strategy interface {
	IsTimeToBuy() bool
	IsTimeToSell() bool
}

type Indicator interface {
	Calculate(price float64)
	RecommendSell() bool
	RecommendBuy() bool
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
