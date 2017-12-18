package indicators

import (
	"math"

	"github.com/jeremyhahn/tradebot/common"
)

type EMA struct {
	size         int
	candlesticks []common.Candlestick
	prices       []float64
	count        int
	index        int
	average      float64
	last         float64
	multiplier   float64
	common.MovingAverage
	common.PeriodListener
}

func NewExponentialMovingAverage(candles []common.Candlestick) *EMA {
	return CreateExponentialMovingAverage(candles, len(candles))
}

func CreateExponentialMovingAverage(candles []common.Candlestick, size int) *EMA {
	var prices []float64
	var total float64
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total + c.Close
	}
	ema := &EMA{
		prices:       make([]float64, size),
		size:         size,
		candlesticks: candles,
		index:        0,
		count:        0,
		average:      0.0,
		last:         0.0,
		multiplier:   1}
	candleLen := len(candles)
	if candles[0].Close > 0 {
		ema.prices = prices
		ema.count = candleLen
		ema.index = candleLen - 1
		ema.average = total / float64(size)
		ema.last = total / float64(size)
		ema.multiplier = 2 / (float64(size) + 1)
	}
	return ema
}

func (ema *EMA) Add(candle *common.Candlestick) float64 {
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

func (ema *EMA) GetCandlesticks() []common.Candlestick {
	return ema.candlesticks
}

func (ema *EMA) GetAverage() float64 {
	return ema.average
}

func (ema *EMA) GetSize() int {
	return ema.size
}

func (ema *EMA) GetCount() int {
	return ema.count
}
func (ema *EMA) GetIndex() int {
	return ema.index
}

func (ema *EMA) Sum() float64 {
	var i float64
	for _, price := range ema.prices {
		i += price
	}
	return i
}

func (ema *EMA) GetMultiplier() float64 {
	return ema.multiplier
}

func (ema *EMA) GetGainsAndLosses() (float64, float64) {
	var u, d float64
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

func (ema *EMA) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[EMA] OnPeriodChange: ", candle.Date, candle.Close)
	ema.Add(candle)
}
