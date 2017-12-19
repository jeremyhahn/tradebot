package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type MACD struct {
	ema1       common.MovingAverage
	ema2       common.MovingAverage
	ema3       common.MovingAverage
	signalSize int
	signal     decimal.Decimal
	value      decimal.Decimal
	histogram  decimal.Decimal
	lastSignal decimal.Decimal
	common.Indicator
}

func NewMovingAverageConvergenceDivergence(ema1, ema2 common.MovingAverage, signalSize int) *MACD {
	ema1Avg := ema1.GetAverage()
	ema2Avg := ema2.GetAverage()

	ema3Candles := make([]common.Candlestick, signalSize)
	ema3 := NewExponentialMovingAverage(ema3Candles)
	ema3.Add(&common.Candlestick{Close: ema1.GetAverage().Sub(ema2.GetAverage())})

	return &MACD{
		ema1:       ema1, // 12-period
		ema2:       ema2, // 26-period
		ema3:       ema3, // 9-period
		value:      ema1Avg.Sub(ema2Avg),
		signal:     decimal.NewFromFloat(0),
		signalSize: signalSize,
		histogram:  (ema1Avg.Sub(ema2Avg)).Sub(ema1Avg),
		lastSignal: decimal.NewFromFloat(0.0)}
}

func (macd *MACD) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	var value, signal, histogram decimal.Decimal
	zero := decimal.NewFromFloat(float64(0))
	if (macd.ema3.GetIndex()+1) == macd.signalSize || macd.lastSignal.GreaterThan(zero) {
		value = macd.ema1.GetAverage().Sub(macd.ema2.GetAverage())
		tmpEma := macd.ema3
		tmpEma.Add(&common.Candlestick{Close: price})
		if macd.signal.Equals(zero) {
			signal = macd.ema3.Sum().Div(decimal.NewFromFloat(float64(macd.ema3.GetSize())))
		} else {
			signal = macd.value.Mul((decimal.NewFromFloat(float64(2)).Div((decimal.NewFromFloat(float64((macd.ema3.GetSize())) + 1))))).Add((macd.lastSignal.Mul((decimal.NewFromFloat(1).Sub((decimal.NewFromFloat(float64(2))).Div((decimal.NewFromFloat(float64(macd.ema3.GetSize()) + 1))))))))
		}
		histogram = value.Sub(signal)
	}
	return value, signal, histogram
}

func (macd *MACD) GetValue() decimal.Decimal {
	return macd.value
}

func (macd *MACD) GetSignalLine() decimal.Decimal {
	return macd.signal
}

func (macd *MACD) GetHistogram() decimal.Decimal {
	return macd.histogram
}

func (macd *MACD) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[MACD] OnPeriodChange: ", candle.Date, candle.Close)
	zero := decimal.NewFromFloat(float64(0))
	one := decimal.NewFromFloat(float64(1))
	two := decimal.NewFromFloat(float64(2))
	macd.ema1.Add(candle)
	macd.ema2.Add(candle)
	macd.value = macd.ema1.GetAverage().Sub(macd.ema2.GetAverage())
	macd.ema3.Add(&common.Candlestick{Close: macd.value})
	if (macd.ema3.GetIndex()+1) == macd.signalSize || macd.lastSignal.GreaterThan(zero) {
		if macd.signal.Equals(zero) {
			macd.signal = macd.ema3.Sum().Div(decimal.NewFromFloat(float64(macd.ema3.GetSize())))
		} else {
			macd.signal = macd.value.Mul((two.Div((decimal.NewFromFloat(float64(macd.ema3.GetSize() + 1))))).Add((macd.lastSignal.Mul((one.Sub(two.Div((decimal.NewFromFloat(float64(macd.ema3.GetSize()) + 1)))))))))
		}
		macd.lastSignal = macd.signal
		macd.histogram = macd.value.Sub(macd.signal)
	}
}
