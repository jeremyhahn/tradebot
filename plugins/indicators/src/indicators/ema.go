package indicators

import "github.com/jeremyhahn/tradebot/common"

type ExponentialMovingAverage interface {
	Add(candle *common.Candlestick) float64
	GetCandlesticks() []common.Candlestick
	GetSize() int
	GetCount() int
	GetIndex() int
	GetAverage() float64
	GetPrices() []float64
	Sum() float64
	GetMultiplier() float64
	GetGainsAndLosses() (float64, float64)
	common.FinancialIndicator
}
