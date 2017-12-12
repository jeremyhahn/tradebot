package indicators

import (
	"math"

	"github.com/jeremyhahn/tradebot/common"
)

type RelativeStrengthIndex interface {
	IsOverBought()
	IsOverSold()
	common.Indicator
	common.PriceListener
}

type RSI struct {
	period     int
	sma        SimpleMovingAverage
	oscillator float64
	overbought float64
	oversold   float64
	common.Indicator
	RelativeStrengthIndex
}

func NewRelativeStrengthIndex(sma SimpleMovingAverage) *RSI {
	return &RSI{
		period:     sma.GetSize(),
		sma:        sma,
		oscillator: 0,
		overbought: 70,
		oversold:   30}
}

func CreateRelativeStrengthIndex(sma SimpleMovingAverage, period int, overbought, oversold float64) *RSI {
	return &RSI{
		period:     period,
		sma:        sma,
		oscillator: 0,
		overbought: overbought,
		oversold:   oversold}
}

func (rsi *RSI) Calculate(price float64) float64 {
	var lastClose float64
	var gains float64
	var losses float64
	var avgGain float64
	var avgLoss float64
	var rs float64
	candles := rsi.sma.GetCandlesticks()
	candles = candles[:rsi.period-1]
	candles = append(candles, common.Candlestick{Close: price})
	for _, c := range candles {
		if lastClose == 0 {
			lastClose = c.Close
			continue
		}
		difference := (c.Close - lastClose)
		if difference < 0 {
			losses = losses + math.Abs(difference)
		} else {
			gains = gains + difference
		}
		lastClose = c.Close
	}
	avgGain = gains / float64(rsi.period)
	avgLoss = losses / float64(rsi.period)
	rs = (avgGain / avgLoss)
	rsi.oscillator = (100 - (100 / (1 + rs)))
	return rsi.oscillator
}

func (rsi *RSI) RecommendBuy() bool {
	return rsi.IsOverSold()
}

func (rsi *RSI) RecommendSell() bool {
	return rsi.IsOverBought()
}

func (rsi *RSI) IsOverBought() bool {
	return rsi.oscillator >= rsi.overbought
}

func (rsi *RSI) IsOverSold() bool {
	return rsi.oscillator <= rsi.oversold
}
