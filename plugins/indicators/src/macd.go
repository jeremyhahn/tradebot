package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type MovingAverageConvergenceDivergenceParams struct {
	EMA1Period int64
	EMA2Period int64
	SignalSize int64
}

type MovingAverageConvergenceDivergenceImpl struct {
	name        string
	displayName string
	params      *MovingAverageConvergenceDivergenceParams
	ema1        indicators.ExponentialMovingAverage
	ema2        indicators.ExponentialMovingAverage
	ema3        indicators.ExponentialMovingAverage
	signal      decimal.Decimal
	value       decimal.Decimal
	histogram   decimal.Decimal
	lastSignal  decimal.Decimal
	common.FinancialIndicator
}

func NewMovingAverageConvergenceDivergence(candles []common.Candlestick) (common.FinancialIndicator, error) {
	params := []string{"12", "26", "9"}
	return CreateMovingAverageConvergenceDivergence(candles, params)
}

func CreateMovingAverageConvergenceDivergence(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &MovingAverageConvergenceDivergenceImpl{}
		params = temp.GetDefaultParameters()
	}
	ema1Period, _ := strconv.ParseInt(params[0], 10, 32)
	ema2Period, _ := strconv.ParseInt(params[1], 10, 32)
	signalSize, _ := strconv.ParseInt(params[2], 10, 32)
	ema1Indicator, err := NewExponentialMovingAverage(candles[:ema1Period])
	if err != nil {
		return nil, err
	}
	ema1 := ema1Indicator.(indicators.ExponentialMovingAverage)
	ema2Indicator, err2 := NewExponentialMovingAverage(candles[:ema2Period])
	if err2 != nil {
		return nil, err2
	}
	ema2 := ema2Indicator.(indicators.ExponentialMovingAverage)
	for _, c := range candles[ema1Period:ema2Period] {
		ema1.OnPeriodChange(&c)
	}
	ema3Candles := make([]common.Candlestick, signalSize)
	ema3Indicator, err3 := NewExponentialMovingAverage(ema3Candles)
	if err3 != nil {
		return nil, err3
	}
	ema3 := ema3Indicator.(indicators.ExponentialMovingAverage)
	ema3.Add(&common.Candlestick{Close: ema1.GetAverage().Sub(ema2.GetAverage())})
	for _, c := range candles[ema2Period:] {
		ema3.OnPeriodChange(&c)
	}
	ema1Avg := ema1.GetAverage()
	ema2Avg := ema2.GetAverage()
	macd := &MovingAverageConvergenceDivergenceImpl{
		name:        "MovingAverageConvergenceDivergence",
		displayName: "Moving Average Convergence Divergence (MACD)",
		params:      &MovingAverageConvergenceDivergenceParams{EMA1Period: ema1Period, EMA2Period: ema2Period, SignalSize: signalSize},
		ema1:        ema1, // 12-period
		ema2:        ema2, // 26-period
		ema3:        ema3, // 9-period
		value:       ema1Avg.Sub(ema2Avg),
		signal:      decimal.NewFromFloat(0),
		histogram:   ema1Avg.Sub(ema2Avg).Sub(ema1Avg),
		lastSignal:  decimal.NewFromFloat(0)}

	for _, c := range candles[ema2Period:] {
		macd.OnPeriodChange(&c)
	}

	return macd, nil
}

func (macd *MovingAverageConvergenceDivergenceImpl) Calculate(price decimal.Decimal) (decimal.Decimal,
	decimal.Decimal, decimal.Decimal) {
	var value, signal, histogram decimal.Decimal
	if macd.ema3 != nil {
		prices := macd.ema3.GetPrices()
		sum := decimal.NewFromFloat(0)
		for _, p := range prices {
			sum = sum.Add(p)
		}
		size := macd.ema3.GetSize()
		value = macd.ema1.GetAverage().Sub(macd.ema2.GetAverage())
		if macd.signal.Equals(decimal.NewFromFloat(0)) {
			signal = sum.Div(decimal.NewFromFloat(float64(size)))
		} else {
			//signal = macd.value*(2/(float64(size+1))) + (macd.lastSignal * (1 - (2 / (float64(size) + 1))))
			one := decimal.NewFromFloat(1)
			two := decimal.NewFromFloat(2)
			decSize := decimal.NewFromFloat(float64(size + 1))
			signal = macd.value.Mul(two.Div(decSize).Add(macd.lastSignal.Mul(one.Sub(two.Div(decSize.Add(one))))))
		}
		histogram = value.Sub(signal)
	}
	return value, signal, histogram
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetValue() decimal.Decimal {
	return macd.value
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetSignalLine() decimal.Decimal {
	return macd.signal
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetHistogram() decimal.Decimal {
	return macd.histogram
}

func (macd *MovingAverageConvergenceDivergenceImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[MACD] OnPeriodChange: %s", candle.ToString())
	macd.ema1.Add(candle)
	macd.ema2.Add(candle)
	macd.value = macd.ema1.GetAverage().Sub(macd.ema2.GetAverage())
	macd.ema3.Add(&common.Candlestick{Close: macd.value})
	zero := decimal.NewFromFloat(0)
	if (macd.ema3.GetIndex()+1) == int(macd.params.SignalSize) || macd.lastSignal.GreaterThan(zero) {
		decSize := decimal.NewFromFloat(float64(macd.ema3.GetSize()))
		if macd.signal.Equals(zero) {
			macd.signal = macd.ema3.Sum().Div(decSize)
		} else {
			//macd.signal = macd.value*(2/(float64(macd.ema3.GetSize()+1))) + (macd.lastSignal * (1 - (2 / (float64(macd.ema3.GetSize()) + 1))))
			one := decimal.NewFromFloat(1)
			two := decimal.NewFromFloat(2)

			macd.signal = macd.value.Mul(two.Div(decSize.Add(one))).Add(macd.lastSignal.Mul(one.Sub(two.Div(decSize.Add(one)))))
		}
		macd.lastSignal = macd.signal
		macd.histogram = macd.value.Sub(macd.signal)
	}
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetName() string {
	return macd.name
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetDisplayName() string {
	return macd.displayName
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetDefaultParameters() []string {
	return []string{"12", "26", "9"}
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetParameters() []string {
	return []string{
		fmt.Sprintf("%d", macd.params.EMA1Period),
		fmt.Sprintf("%d", macd.params.EMA2Period),
		fmt.Sprintf("%d", macd.params.SignalSize)}
}
