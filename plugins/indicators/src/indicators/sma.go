package indicators

import "github.com/jeremyhahn/tradebot/common"

type SimpleMovingAverage interface {
	Add(candle *common.Candlestick) float64
	GetCandlesticks() []common.Candlestick
	GetAverage() float64
	GetSize() int
	GetPrices() []float64
	GetGainsAndLosses() (float64, float64)
	GetCount() int
	GetIndex() int
	Sum() float64
	common.FinancialIndicator
}
