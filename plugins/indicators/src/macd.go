package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
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
	signal      float64
	value       float64
	histogram   float64
	lastSignal  float64
	common.FinancialIndicator
}

func NewMovingAverageConvergenceDivergence(candles []common.Candlestick) common.FinancialIndicator {
	params := []string{"12", "26", "9"}
	return CreateMovingAverageConvergenceDivergence(candles, params)
}

func CreateMovingAverageConvergenceDivergence(candles []common.Candlestick, params []string) common.FinancialIndicator {
	if params == nil {
		temp := &MovingAverageConvergenceDivergenceImpl{}
		params = temp.GetDefaultParameters()
	}
	ema1Period, _ := strconv.ParseInt(params[0], 10, 32)
	ema2Period, _ := strconv.ParseInt(params[1], 10, 32)
	signalSize, _ := strconv.ParseInt(params[2], 10, 32)
	ema1 := NewExponentialMovingAverage(candles[:ema1Period]).(indicators.ExponentialMovingAverage)
	ema2 := NewExponentialMovingAverage(candles[:ema2Period]).(indicators.ExponentialMovingAverage)
	for _, c := range candles[ema1Period:ema2Period] {
		ema1.OnPeriodChange(&c)
	}
	ema3Candles := make([]common.Candlestick, signalSize)
	ema3 := NewExponentialMovingAverage(ema3Candles).(indicators.ExponentialMovingAverage)
	ema3.Add(&common.Candlestick{Close: ema1.GetAverage() - ema2.GetAverage()})
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
		value:       ema1Avg - ema2Avg,
		signal:      0,
		histogram:   (ema1Avg - ema2Avg) - ema1Avg,
		lastSignal:  0.0}

	for _, c := range candles[ema2Period:] {
		macd.OnPeriodChange(&c)
	}

	return macd
}

func (macd *MovingAverageConvergenceDivergenceImpl) Calculate(price float64) (float64, float64, float64) {
	var value, signal, histogram float64

	if macd.ema3 != nil {
		prices := macd.ema3.GetPrices()
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

func (macd *MovingAverageConvergenceDivergenceImpl) GetValue() float64 {
	return macd.value
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetSignalLine() float64 {
	return macd.signal
}

func (macd *MovingAverageConvergenceDivergenceImpl) GetHistogram() float64 {
	return macd.histogram
}

func (macd *MovingAverageConvergenceDivergenceImpl) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[MovingAverageConvergenceDivergence] OnPeriodChange: ", candle.Date, candle.Close)
	macd.ema1.Add(candle)
	macd.ema2.Add(candle)
	macd.value = macd.ema1.GetAverage() - macd.ema2.GetAverage()
	macd.ema3.Add(&common.Candlestick{Close: macd.value})
	if (macd.ema3.GetIndex()+1) == int(macd.params.SignalSize) || macd.lastSignal > 0 {
		if macd.signal == 0 {
			macd.signal = macd.ema3.Sum() / float64(macd.ema3.GetSize())
		} else {
			macd.signal = macd.value*(2/(float64(macd.ema3.GetSize()+1))) + (macd.lastSignal * (1 - (2 / (float64(macd.ema3.GetSize()) + 1))))
		}
		macd.lastSignal = macd.signal
		macd.histogram = macd.value - macd.signal
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
