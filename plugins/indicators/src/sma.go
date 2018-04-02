package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type SimpleMovingAverageImpl struct {
	name         string
	displayName  string
	size         int
	candlesticks []common.Candlestick
	prices       []decimal.Decimal
	count        int
	index        int
	average      decimal.Decimal
	sum          decimal.Decimal
	indicators.SimpleMovingAverage
}

func NewSimpleMovingAverage(candles []common.Candlestick) (common.FinancialIndicator, error) {
	var params []string
	params = append(params, fmt.Sprintf("%d", len(candles)))
	return CreateSimpleMovingAverage(candles, params)
}

func CreateSimpleMovingAverage(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &SimpleMovingAverageImpl{}
		params = temp.GetDefaultParameters()
	}
	size, _ := strconv.ParseInt(params[0], 10, 64)
	var prices []decimal.Decimal
	var total decimal.Decimal
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total.Add(c.Close)
	}
	sma := &SimpleMovingAverageImpl{
		name:         "SimpleMovingAverage",
		displayName:  "Simple Moving Average (SMA)",
		size:         int(size),
		candlesticks: candles,
		prices:       make([]decimal.Decimal, size),
		count:        0,
		index:        0,
		average:      decimal.NewFromFloat(0),
		sum:          decimal.NewFromFloat(0)}
	candleLen := len(candles)
	if candleLen > 0 {
		sma.prices = prices
		sma.count = candleLen
		sma.index = candleLen - 1
		sma.average = total.Div(decimal.NewFromFloat(float64(size)))
	}
	return sma, nil
}

func (sma *SimpleMovingAverageImpl) Add(candle *common.Candlestick) decimal.Decimal {
	if sma.count == sma.size {
		sma.index = (sma.index + 1) % sma.size
		//sma.average += (candle.Close - sma.prices[sma.index]) / float64(sma.count)
		sma.average = sma.average.Add(candle.Close.Sub(sma.prices[sma.index]).Div(decimal.NewFromFloat(float64(sma.count))))
		sma.prices[sma.index] = candle.Close
	} else if sma.count != 0 && sma.count < sma.size {
		sma.index = (sma.index + 1) % sma.size
		//sma.average = (candle.Close + float64(sma.count)*sma.average) / float64(sma.count+1)
		sma.average = sma.average.Add(candle.Close.Add(decimal.NewFromFloat(float64(sma.count))).Mul(sma.average).Div(decimal.NewFromFloat(float64(sma.count + 1))))
		sma.prices[sma.index] = candle.Close
		sma.count++
	} else {
		sma.average = candle.Close
		sma.prices[0] = candle.Close
		sma.count = 1
	}
	return sma.average
}

func (sma *SimpleMovingAverageImpl) GetCandlesticks() []common.Candlestick {
	return sma.candlesticks
}

func (sma *SimpleMovingAverageImpl) GetAverage() decimal.Decimal {
	return sma.average
}

func (sma *SimpleMovingAverageImpl) GetSize() int {
	return sma.size
}

func (sma *SimpleMovingAverageImpl) GetCount() int {
	return sma.count
}

func (sma *SimpleMovingAverageImpl) GetPrices() []decimal.Decimal {
	return sma.prices
}

func (sma *SimpleMovingAverageImpl) GetGainsAndLosses() (decimal.Decimal, decimal.Decimal) {
	var gains, losses decimal.Decimal
	if len(sma.candlesticks) <= 0 {
		return decimal.NewFromFloat(0), decimal.NewFromFloat(0)
	}
	var lastClose = sma.candlesticks[0].Close
	zero := decimal.NewFromFloat(0)
	for _, c := range sma.candlesticks {
		difference := c.Close.Sub(lastClose)
		if difference.LessThan(zero) {
			losses = losses.Add(difference.Abs())
		} else {
			gains = gains.Add(difference)
		}
		lastClose = c.Close
	}
	return gains, losses
}

func (sma *SimpleMovingAverageImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[SMA] OnPeriodChange: %s", candle.ToString())
	sma.Add(candle)
}

func (sma *SimpleMovingAverageImpl) GetName() string {
	return sma.name
}

func (sma *SimpleMovingAverageImpl) GetDisplayName() string {
	return sma.displayName
}

func (sma *SimpleMovingAverageImpl) GetDefaultParameters() []string {
	return []string{"20"}
}

func (sma *SimpleMovingAverageImpl) GetParameters() []string {
	return []string{fmt.Sprintf("%d", sma.size)}
}
