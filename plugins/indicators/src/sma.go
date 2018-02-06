package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
)

type SimpleMovingAverageImpl struct {
	name         string
	displayName  string
	size         int
	candlesticks []common.Candlestick
	prices       []float64
	count        int
	index        int
	average      float64
	sum          float64
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
	var prices []float64
	var total float64
	for _, c := range candles {
		prices = append(prices, c.Close)
		total += c.Close
	}
	sma := &SimpleMovingAverageImpl{
		name:         "SimpleMovingAverage",
		displayName:  "Simple Moving Average (SMA)",
		size:         int(size),
		candlesticks: candles,
		prices:       make([]float64, size),
		count:        0,
		index:        0,
		average:      0,
		sum:          0}
	candleLen := len(candles)
	if candleLen > 0 {
		sma.prices = prices
		sma.count = candleLen
		sma.index = candleLen - 1
		sma.average = total / float64(size)
	}
	return sma, nil
}

func (sma *SimpleMovingAverageImpl) Add(candle *common.Candlestick) float64 {
	if sma.count == sma.size {
		sma.index = (sma.index + 1) % sma.size
		sma.average += (candle.Close - sma.prices[sma.index]) / float64(sma.count)
		sma.prices[sma.index] = candle.Close
	} else if sma.count != 0 && sma.count < sma.size {
		sma.index = (sma.index + 1) % sma.size
		sma.average = (candle.Close + float64(sma.count)*sma.average) / float64(sma.count+1)
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

func (sma *SimpleMovingAverageImpl) GetAverage() float64 {
	return sma.average
}

func (sma *SimpleMovingAverageImpl) GetSize() int {
	return sma.size
}

func (sma *SimpleMovingAverageImpl) GetCount() int {
	return sma.count
}

func (sma *SimpleMovingAverageImpl) GetPrices() []float64 {
	return sma.prices
}

func (sma *SimpleMovingAverageImpl) GetGainsAndLosses() (float64, float64) {
	var gains, losses float64
	if len(sma.candlesticks) <= 0 {
		return 0.0, 0.0
	}
	var lastClose = sma.candlesticks[0].Close
	for _, c := range sma.candlesticks {
		difference := (c.Close - lastClose)
		if difference < 0 {
			losses += math.Abs(difference)
		} else {
			gains += difference
		}
		lastClose = c.Close
	}
	return gains, losses
}

func (sma *SimpleMovingAverageImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[SimpleMovingAverage] OnPeriodChange: ", candle.Date, candle.Close)
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
