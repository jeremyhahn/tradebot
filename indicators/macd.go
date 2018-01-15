package indicators

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type MACD struct {
	ema1       common.MovingAverage
	ema2       common.MovingAverage
	ema3       common.MovingAverage
	signalSize int
	signal     float64
	value      float64
	histogram  float64
	lastSignal float64
	common.Indicator
}

func NewMovingAverageConvergenceDivergence(ema1, ema2 common.MovingAverage, signalSize int) *MACD {
	ema1Avg := ema1.GetAverage()
	ema2Avg := ema2.GetAverage()

	ema3Candles := make([]common.Candlestick, signalSize)
	ema3 := NewExponentialMovingAverage(ema3Candles)
	ema3.Add(&common.Candlestick{Close: ema1.GetAverage() - ema2.GetAverage()})

	return &MACD{
		ema1:       ema1, // 12-period
		ema2:       ema2, // 26-period
		ema3:       ema3, // 9-period
		value:      ema1Avg - ema2Avg,
		signal:     0,
		signalSize: signalSize,
		histogram:  (ema1Avg - ema2Avg) - ema1Avg,
		lastSignal: 0.0}
}

func (macd *MACD) Calculate(price float64) (float64, float64, float64) {
	var value, signal, histogram float64

	if macd.ema3 != nil {
		prices := macd.ema3.GetPrices()
		///priceLen := len(prices)
		sum := 0.0
		for _, p := range prices {
			sum += p
		}
		size := macd.ema3.GetSize()
		value = macd.ema1.GetAverage() - macd.ema2.GetAverage()
		if macd.signal == 0 {
			signal = sum / float64(size)
		} else {
			signal = macd.value*(2/(float64(size+1))) + (macd.lastSignal * (1 - (2 / (float64(size) + 1))))
		}
		histogram = value - signal
	}
	return value, signal, histogram
}

func (macd *MACD) GetValue() float64 {
	return macd.value
}

func (macd *MACD) GetSignalLine() float64 {
	return macd.signal
}

func (macd *MACD) GetHistogram() float64 {
	return macd.histogram
}

func (macd *MACD) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[MACD] OnPeriodChange: ", candle.Date, candle.Close)
	macd.ema1.Add(candle)
	macd.ema2.Add(candle)
	macd.value = macd.ema1.GetAverage() - macd.ema2.GetAverage()
	macd.ema3.Add(&common.Candlestick{Close: macd.value})
	if (macd.ema3.GetIndex()+1) == macd.signalSize || macd.lastSignal > 0 {
		if macd.signal == 0 {
			macd.signal = macd.ema3.Sum() / float64(macd.ema3.GetSize())
		} else {
			macd.signal = macd.value*(2/(float64(macd.ema3.GetSize()+1))) + (macd.lastSignal * (1 - (2 / (float64(macd.ema3.GetSize()) + 1))))
		}
		macd.lastSignal = macd.signal
		macd.histogram = macd.value - macd.signal
	}
}
