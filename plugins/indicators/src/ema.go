package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type ExponentialMovingAverageImpl struct {
	name         string
	displayName  string
	size         int
	candlesticks []common.Candlestick
	prices       []decimal.Decimal
	count        int
	index        int
	average      decimal.Decimal
	last         decimal.Decimal
	multiplier   decimal.Decimal
	indicators.ExponentialMovingAverage
}

func NewExponentialMovingAverage(candles []common.Candlestick) (common.FinancialIndicator, error) {
	params := []string{fmt.Sprintf("%d", len(candles))}
	return CreateExponentialMovingAverage(candles, params)
}

func CreateExponentialMovingAverage(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &ExponentialMovingAverageImpl{}
		params = temp.GetDefaultParameters()
	}
	size, _ := strconv.ParseInt(params[0], 10, 64)
	var prices []decimal.Decimal
	var total decimal.Decimal
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total.Add(c.Close)
	}
	ema := &ExponentialMovingAverageImpl{
		name:         "ExponentialMovingAverage",
		displayName:  "Exponential Moving Average (EMA)",
		prices:       make([]decimal.Decimal, size),
		size:         int(size),
		candlesticks: candles,
		index:        0,
		count:        0,
		average:      decimal.NewFromFloat(0),
		last:         decimal.NewFromFloat(0),
		multiplier:   decimal.NewFromFloat(1)}
	candleLen := len(candles)
	if candleLen > 0 && candles[0].Close.GreaterThan(decimal.NewFromFloat(0)) {
		decSize := decimal.NewFromFloat(float64(size))
		ema.prices = prices
		ema.count = candleLen
		ema.index = candleLen - 1
		ema.average = total.Div(decSize)
		ema.last = total.Div(decSize)
		ema.multiplier = decimal.NewFromFloat(2).Div(decSize).Add(decimal.NewFromFloat(1))
	}
	return ema, nil
}

func (ema *ExponentialMovingAverageImpl) Add(candle *common.Candlestick) decimal.Decimal {
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

func (ema *ExponentialMovingAverageImpl) GetCandlesticks() []common.Candlestick {
	return ema.candlesticks
}

func (ema *ExponentialMovingAverageImpl) GetPrices() []decimal.Decimal {
	return ema.prices
}

func (ema *ExponentialMovingAverageImpl) GetAverage() decimal.Decimal {
	return ema.average
}

func (ema *ExponentialMovingAverageImpl) GetSize() int {
	return ema.size
}

func (ema *ExponentialMovingAverageImpl) GetCount() int {
	return ema.count
}
func (ema *ExponentialMovingAverageImpl) GetIndex() int {
	return ema.index
}

func (ema *ExponentialMovingAverageImpl) Sum() decimal.Decimal {
	var i decimal.Decimal
	for _, price := range ema.prices {
		i = i.Add(price)
	}
	return i
}

func (ema *ExponentialMovingAverageImpl) GetMultiplier() decimal.Decimal {
	return ema.multiplier
}

func (ema *ExponentialMovingAverageImpl) GetGainsAndLosses() (decimal.Decimal, decimal.Decimal) {
	var u, d decimal.Decimal
	if len(ema.candlesticks) <= 0 {
		return decimal.NewFromFloat(0), decimal.NewFromFloat(0)
	}
	var lastClose = ema.candlesticks[0].Close
	for _, c := range ema.candlesticks {
		difference := (c.Close.Sub(lastClose))
		if difference.LessThan(decimal.NewFromFloat(0)) {
			d = d.Add(difference.Abs())
		} else {
			u = u.Add(difference)
		}
		lastClose = c.Close
	}
	return u, d
}

func (ema *ExponentialMovingAverageImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[EMA] OnPeriodChange: %s", candle.ToString())
	ema.Add(candle)
}

func (ema *ExponentialMovingAverageImpl) GetName() string {
	return ema.name
}

func (ema *ExponentialMovingAverageImpl) GetDisplayName() string {
	return ema.displayName
}

func (ema *ExponentialMovingAverageImpl) GetDefaultParameters() []string {
	return []string{"20"}
}

func (ema *ExponentialMovingAverageImpl) GetParameters() []string {
	return []string{fmt.Sprintf("%d", ema.size)}
}
