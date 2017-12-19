package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type EMA struct {
	size         int
	candlesticks []common.Candlestick
	prices       []decimal.Decimal
	count        int
	index        int
	average      decimal.Decimal
	last         decimal.Decimal
	multiplier   decimal.Decimal
	common.MovingAverage
	common.PeriodListener
}

func NewExponentialMovingAverage(candles []common.Candlestick) *EMA {
	return CreateExponentialMovingAverage(candles, len(candles))
}

func CreateExponentialMovingAverage(candles []common.Candlestick, size int) *EMA {
	fsize := float64(size)
	var prices []decimal.Decimal
	var total decimal.Decimal
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total.Add(c.Close)
	}
	ema := &EMA{
		prices:       make([]decimal.Decimal, size),
		size:         size,
		candlesticks: candles,
		index:        0,
		count:        0,
		average:      decimal.NewFromFloat(0.0),
		last:         decimal.NewFromFloat(0.0),
		multiplier:   decimal.NewFromFloat(0.0)}
	candleLen := len(candles)
	if candles[0].Close.GreaterThan(decimal.NewFromFloat(0)) {
		ema.prices = prices
		ema.count = candleLen
		ema.index = candleLen - 1
		ema.average = total.Div(decimal.NewFromFloat(fsize))
		ema.last = total.Div(decimal.NewFromFloat(fsize))
		ema.multiplier = decimal.NewFromFloat(2).Div(decimal.NewFromFloat(fsize)).Add(decimal.NewFromFloat(1))
	}
	return ema
}

func (ema *EMA) Add(candle *common.Candlestick) decimal.Decimal {
	if ema.count == ema.size {
		ema.index = (ema.index + 1) % ema.size
		ema.prices[ema.index] = candle.Close
		step1 := candle.Close.Sub(ema.last)
		step2 := step1.Mul(ema.multiplier)
		ema.average = step2.Add(ema.last)
		ema.last = ema.average
	} else if ema.count != 0 && ema.count < ema.size {
		ema.index = (ema.index + 1) % ema.size
		ema.prices[ema.index] = candle.Close
		step1 := candle.Close.Sub(ema.last)
		step2 := step1.Mul(ema.multiplier)
		ema.average = step2.Add(ema.last)
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

func (ema *EMA) GetAverage() decimal.Decimal {
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

func (ema *EMA) Sum() decimal.Decimal {
	var i decimal.Decimal
	for _, price := range ema.prices {
		i = i.Add(price)
	}
	return i
}

func (ema *EMA) GetMultiplier() decimal.Decimal {
	return ema.multiplier
}

func (ema *EMA) GetGainsAndLosses() (decimal.Decimal, decimal.Decimal) {
	var u, d decimal.Decimal
	var lastClose = ema.candlesticks[0].Close
	for _, c := range ema.candlesticks {
		difference := c.Close.Sub(lastClose)
		if difference.LessThan(decimal.NewFromFloat(0)) {
			d = d.Add(difference.Abs())
		} else {
			u = u.Add(difference)
		}
		lastClose = c.Close
	}
	return u, d
}

func (ema *EMA) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[EMA] OnPeriodChange: ", candle.Date, candle.Close)
	ema.Add(candle)
}
