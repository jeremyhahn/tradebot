package indicators

import (
	"fmt"
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
}

func NewExponentialMovingAverage(candles []common.Candlestick) *EMA {
	return CreateExponentiaMovingAverage(candles, len(candles))
}

func CreateExponentiaMovingAverage(candles []common.Candlestick, size int) *EMA {
	var prices []float64
	var total float64
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total + c.Close
	}
	return &EMA{
		size:         size,
		candlesticks: candles,
		prices:       prices,
		count:        len(candles),
		index:        len(candles) - 1,
		average:      total / float64(size),
		last:         total / float64(size),
		multiplier:   2 / (float64(size) + 1)}
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
		ema.average = (candle.Close + float64(ema.count)*ema.average) / float64(ema.count+1)
		ema.prices[ema.index] = candle.Close
		ema.count++
	} else {
		ema.average = candle.Close
		ema.prices[0] = candle.Close
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

func (ema *EMA) OnCandlestickCreated(candle *common.Candlestick) {
	fmt.Println("EMA OnCandlestickCreated: ", candle.Date, candle.Close)
	ema.Add(candle)
}

func (ema *EMA) OnPriceChange(price float64) {
	//fmt.Println("SMA OnPriceChange: ", price)
}
