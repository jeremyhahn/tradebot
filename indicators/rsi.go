package indicators

import (
	"math"
	"reflect"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

type RelativeStrengthIndex interface {
	IsOverBought()
	IsOverSold()
	common.Indicator
	common.PriceListener
}

type RSI struct {
	period     int
	ma         common.MovingAverage
	oscillator float64
	overbought float64
	oversold   float64
	common.Indicator
	RelativeStrengthIndex
}

func NewRelativeStrengthIndex(ma common.MovingAverage) *RSI {
	return &RSI{
		period:     14,
		ma:         ma,
		oscillator: 0,
		overbought: 70,
		oversold:   30}
}

func CreateRelativeStrengthIndex(ma common.MovingAverage, period int, overbought, oversold float64) *RSI {
	return &RSI{
		period:     period,
		ma:         ma,
		oscillator: 0,
		overbought: overbought,
		oversold:   oversold}
}

func (rsi *RSI) Calculate(price float64) float64 {
	var u float64
	var d float64
	var rs float64
	candles := rsi.ma.GetCandlesticks()
	candles = append(candles, common.Candlestick{Close: price})
	lastClose := candles[0].Close
	for _, c := range candles {
		difference := (c.Close - lastClose)
		if difference < 0 {
			d = d + math.Abs(difference)
		} else {
			u = u + difference
		}
		lastClose = c.Close
	}
	if reflect.TypeOf(rsi.ma).String() == "*indicators.SMA" {
		avgU := u / float64(rsi.period)
		avgD := d / float64(rsi.period)
		rs = avgU / avgD
		rsi.oscillator = (100 - (100 / (1 + rs)))
	} else if reflect.TypeOf(rsi.ma).String() == "*indicators.EMA" {
		a := float64(2 / (rsi.period + 1))
		avgUt := a*u + (1-a)*u - 1
		avgDt := a*d + (1-a)*d - 1
		rs = (avgUt / avgDt)
		rsi.oscillator = (100 - (100 / (1 + rs)))
	}
	return util.FloatPrecision(rsi.oscillator, 2)
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
