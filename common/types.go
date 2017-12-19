package common

import (
	"time"

	"github.com/shopspring/decimal"
)

type WebsocketBroadcast struct {
	Price    decimal.Decimal
	Currency string
}

type ChartData struct {
	Currency          string          `json:"currency"`
	Price             decimal.Decimal `json:"price"`
	MACDValue         decimal.Decimal `json:"macd_value"`
	MACDHistogram     decimal.Decimal `json:"macd_histogram"`
	MACDSignal        decimal.Decimal `json:"macd_signal"`
	MACDValueLive     decimal.Decimal `json:"macd_value_live"`
	MACDHistogramLive decimal.Decimal `json:"macd_histogram_live"`
	MACDSignalLive    decimal.Decimal `json:"macd_signal_live"`
	RSI               decimal.Decimal `json:"rsi"`
	RSILive           decimal.Decimal `json:"rsi_live"`
	BollingerUpper    decimal.Decimal `json:"bband_upper"`
	BollingerMiddle   decimal.Decimal `json:"bband_middle"`
	BollingerLower    decimal.Decimal `json:"bband_lower"`
}

type MovingAverage interface {
	Add(candle *Candlestick) decimal.Decimal
	GetCandlesticks() []Candlestick
	GetSize() int
	GetCount() int
	GetIndex() int
	GetAverage() decimal.Decimal
	Sum() decimal.Decimal
	GetGainsAndLosses() (decimal.Decimal, decimal.Decimal)
	PeriodListener
}

type PriceListener interface {
	OnPriceChange(price decimal.Decimal)
}

type PeriodListener interface {
	OnPeriodChange(candlestick *Candlestick)
}

type Exchange interface {
	SubscribeToLiveFeed(price chan decimal.Decimal)
	GetPrice() decimal.Decimal
	GetTradeHistory(start, end time.Time, granularity int) []Candlestick
	GetCurrency() string
}

type Indicator interface {
	Calculate(price decimal.Decimal)
	PeriodListener
}

type Account struct {
	Currency decimal.Decimal
	Balance  decimal.Decimal
}

type Trade struct {
	ID        int `gorm:"primary_key"`
	Timestamp int32
	Price     decimal.Decimal
	Size      decimal.Decimal
}
