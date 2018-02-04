package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
)

type ExponentialMovingAverage struct {
	name         string
	displayName  string
	size         int
	candlesticks []common.Candlestick
	prices       []float64
	count        int
	index        int
	average      float64
	last         float64
	multiplier   float64
	indicators.ExponentialMovingAverage
}

func NewExponentialMovingAverage(candles []common.Candlestick) indicators.ExponentialMovingAverage {
	params := []string{fmt.Sprintf("%d", len(candles))}
	return CreateExponentialMovingAverage(candles, params)
}

func CreateExponentialMovingAverage(candles []common.Candlestick, params []string) indicators.ExponentialMovingAverage {
	size, _ := strconv.ParseInt(params[0], 10, 64)
	var prices []float64
	var total float64
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total + c.Close
	}
	ema := &ExponentialMovingAverage{
		name:         "ExponentialMovingAverage",
		displayName:  "Exponential Moving Average (EMA)",
		prices:       make([]float64, size),
		size:         int(size),
		candlesticks: candles,
		index:        0,
		count:        0,
		average:      0.0,
		last:         0.0,
		multiplier:   1}
	candleLen := len(candles)
	if candleLen > 0 && candles[0].Close > 0 {
		ema.prices = prices
		ema.count = candleLen
		ema.index = candleLen - 1
		ema.average = total / float64(size)
		ema.last = total / float64(size)
		ema.multiplier = 2 / (float64(size) + 1)
	}
	return ema
}

func (ema *ExponentialMovingAverage) Add(candle *common.Candlestick) float64 {
	if ema.count == ema.size {
		ema.index = (ema.index + 1) % ema.size
		ema.prices[ema.index] = candle.Close
		step1 := candle.Close - ema.last
		step2 := step1 * ema.multiplier
		ema.average = step2 + ema.last
		ema.last = ema.average
	} else if ema.count != 0 && ema.count < ema.size {
		ema.index = (ema.index + 1) % ema.size
		ema.prices[ema.index] = candle.Close
		step1 := candle.Close - ema.last
		step2 := step1 * ema.multiplier
		ema.average = step2 + ema.last
		ema.last = ema.average
		ema.count++
	} else {
		ema.average = candle.Close
		ema.prices[0] = candle.Close
		ema.last = candle.Close
		ema.count = 1
	}
	return ema.average
}

func (ema *ExponentialMovingAverage) GetCandlesticks() []common.Candlestick {
	return ema.candlesticks
}

func (ema *ExponentialMovingAverage) GetPrices() []float64 {
	return ema.prices
}

func (ema *ExponentialMovingAverage) GetAverage() float64 {
	return ema.average
}

func (ema *ExponentialMovingAverage) GetSize() int {
	return ema.size
}

func (ema *ExponentialMovingAverage) GetCount() int {
	return ema.count
}
func (ema *ExponentialMovingAverage) GetIndex() int {
	return ema.index
}

func (ema *ExponentialMovingAverage) Sum() float64 {
	var i float64
	for _, price := range ema.prices {
		i += price
	}
	return i
}

func (ema *ExponentialMovingAverage) GetMultiplier() float64 {
	return ema.multiplier
}

func (ema *ExponentialMovingAverage) GetGainsAndLosses() (float64, float64) {
	var u, d float64
	if len(ema.candlesticks) <= 0 {
		return 0.0, 0.0
	}
	var lastClose = ema.candlesticks[0].Close
	for _, c := range ema.candlesticks {
		difference := (c.Close - lastClose)
		if difference < 0 {
			d += math.Abs(difference)
		} else {
			u += difference
		}
		lastClose = c.Close
	}
	return u, d
}

func (ema *ExponentialMovingAverage) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[ExponentialMovingAverage] OnPeriodChange: ", candle.Date, candle.Close)
	ema.Add(candle)
}

func (ema *ExponentialMovingAverage) GetName() string {
	return ema.name
}

func (ema *ExponentialMovingAverage) GetDisplayName() string {
	return ema.displayName
}

func (ema *ExponentialMovingAverage) GetDefaultParameters() []string {
	return []string{"20"}
}

func (ema *ExponentialMovingAverage) GetParameters() []string {
	return []string{fmt.Sprintf("%d", ema.size)}
}
