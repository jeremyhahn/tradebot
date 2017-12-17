package indicators

import (
	"math"

	"github.com/jeremyhahn/tradebot/common"
)

type RelativeStrengthIndex interface {
	IsOverBought()
	IsOverSold()
	common.Indicator
}

type RSI struct {
	period     int
	ma         common.MovingAverage
	oscillator float64
	overbought float64
	oversold   float64
	u          float64
	d          float64
	avgU       float64
	avgD       float64
	lastPrice  float64
	common.Indicator
	RelativeStrengthIndex
}

func NewRelativeStrengthIndex(ma common.MovingAverage) *RSI {
	return CreateRelativeStrengthIndex(ma, len(ma.GetCandlesticks()), 70, 30)
}

func CreateRelativeStrengthIndex(ma common.MovingAverage, period int, overbought, oversold float64) *RSI {
	candles := ma.GetCandlesticks()
	return &RSI{
		period:     period,
		ma:         ma,
		oscillator: 0,
		overbought: overbought,
		oversold:   oversold,
		u:          0.0,
		d:          0.0,
		avgU:       0.0,
		avgD:       0.0,
		lastPrice:  candles[len(candles)-1].Close}
}

func (rsi *RSI) Calculate(price float64) float64 {
	var oscillator float64
	curU := rsi.u
	curD := -rsi.d
	avgU := rsi.avgU
	avgD := rsi.avgD
	u, d := rsi.ma.GetGainsAndLosses()
	difference := price - rsi.lastPrice
	if difference < 0 {
		d += math.Abs(difference)
		curD = math.Abs(difference)
		curU = 0
	} else {
		u += difference
		curU = difference
		curD = 0
	}
	if avgU > 0 && avgD > 0 {
		avgU = ((avgU*float64(rsi.period-1) + curU) / float64(rsi.period))
		avgD = ((avgD*float64(rsi.period-1) + curD) / float64(rsi.period))
	} else {
		avgU = u / float64(rsi.period)
		avgD = d / float64(rsi.period)
	}
	rs := avgU / avgD
	oscillator = (100 - (100 / (1 + rs)))
	return oscillator
}

func (rsi *RSI) GetValue() float64 {
	return rsi.oscillator
}

func (rsi *RSI) IsOverBought() bool {
	return rsi.oscillator >= rsi.overbought
}

func (rsi *RSI) IsOverSold() bool {
	return rsi.oscillator <= rsi.oversold
}

func (rsi *RSI) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[RSI] OnPeriodChange: ", candle.Date, candle.Close)
	rsi.ma.Add(candle)
	u, d := rsi.ma.GetGainsAndLosses()
	difference := candle.Close - rsi.lastPrice
	if difference < 0 {
		d += math.Abs(difference)
		rsi.d = math.Abs(difference)
		rsi.u = 0
	} else {
		u += difference
		rsi.u = difference
		rsi.d = 0
	}
	if rsi.avgU > 0 && rsi.avgD > 0 {
		rsi.avgU = ((rsi.avgU*float64(rsi.period-1) + rsi.u) / float64(rsi.period))
		rsi.avgD = ((rsi.avgD*float64(rsi.period-1) + rsi.d) / float64(rsi.period))
	} else {
		rsi.avgU = u / float64(rsi.period)
		rsi.avgD = d / float64(rsi.period)
	}
	rs := rsi.avgU / rsi.avgD
	rsi.oscillator = (100 - (100 / (1 + rs)))
	rsi.lastPrice = candle.Close
}
