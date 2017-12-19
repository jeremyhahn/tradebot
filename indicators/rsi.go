package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type RelativeStrengthIndex interface {
	IsOverBought()
	IsOverSold()
	common.Indicator
}

type RSI struct {
	period     int
	ma         common.MovingAverage
	oscillator decimal.Decimal
	overbought decimal.Decimal
	oversold   decimal.Decimal
	u          decimal.Decimal
	d          decimal.Decimal
	avgU       decimal.Decimal
	avgD       decimal.Decimal
	lastPrice  decimal.Decimal
	common.Indicator
	RelativeStrengthIndex
}

func NewRelativeStrengthIndex(ma common.MovingAverage) *RSI {
	seventy := decimal.NewFromFloat(float64(70.00))
	thirty := decimal.NewFromFloat(float64(30.00))
	return CreateRelativeStrengthIndex(ma, len(ma.GetCandlesticks()), seventy, thirty)
}

func CreateRelativeStrengthIndex(ma common.MovingAverage, period int, overbought, oversold decimal.Decimal) *RSI {
	zero := decimal.NewFromFloat(float64(0.0))
	candles := ma.GetCandlesticks()
	return &RSI{
		period:     period,
		ma:         ma,
		oscillator: zero,
		overbought: overbought,
		oversold:   oversold,
		u:          zero,
		d:          zero,
		avgU:       zero,
		avgD:       zero,
		lastPrice:  candles[len(candles)-1].Close}
}

func (rsi *RSI) Calculate(price decimal.Decimal) decimal.Decimal {
	var oscillator decimal.Decimal
	zero := decimal.NewFromFloat(float64(0.0))
	curU := rsi.u
	curD := rsi.d
	avgU := rsi.avgU
	avgD := rsi.avgD
	one := decimal.NewFromFloat(float64(1))
	oneHundred := decimal.NewFromFloat(float64(1))
	u, d := rsi.ma.GetGainsAndLosses()
	difference := price.Sub(rsi.lastPrice)
	if difference.LessThan(zero) {
		dabs := difference.Abs()
		d = d.Add(dabs)
		curD = dabs
		curU = zero
	} else {
		u = u.Add(difference)
		curU = difference
		curD = zero
	}
	if avgU.GreaterThan(zero) && avgD.GreaterThan(zero) {
		avgU = ((avgU.Mul(decimal.NewFromFloat(float64(rsi.period - 1))).Add(curU)).Div(decimal.NewFromFloat(float64(rsi.period))))
		avgD = ((avgD.Mul(decimal.NewFromFloat(float64(rsi.period - 1))).Add(curD)).Div(decimal.NewFromFloat(float64(rsi.period))))
	} else {
		avgU = u.Div(decimal.NewFromFloat(float64(rsi.period)))
		avgD = d.Div(decimal.NewFromFloat(float64(rsi.period)))
	}
	rs := avgU.Div(avgD)
	oscillator = (oneHundred.Sub(oneHundred.Div(one.Add(rs))))
	return oscillator
}

func (rsi *RSI) GetValue() decimal.Decimal {
	return rsi.oscillator
}

func (rsi *RSI) IsOverBought() bool {
	return rsi.oscillator.GreaterThanOrEqual(rsi.overbought)
}

func (rsi *RSI) IsOverSold() bool {
	return rsi.oscillator.LessThanOrEqual(rsi.oversold)
}

func (rsi *RSI) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[RSI] OnPeriodChange: ", candle.Date, candle.Close)
	zero := decimal.NewFromFloat(float64(0))
	one := decimal.NewFromFloat(float64(1))
	hundred := decimal.NewFromFloat(float64(100))
	rsi.ma.Add(candle)
	u, d := rsi.ma.GetGainsAndLosses()
	difference := candle.Close.Sub(rsi.lastPrice)
	if difference.LessThan(zero) {
		d = d.Add(difference.Abs())
		rsi.d = difference.Abs()
		rsi.u = zero
	} else {
		u = u.Add(difference)
		rsi.u = difference
		rsi.d = zero
	}
	if rsi.avgU.GreaterThan(zero) && rsi.avgD.GreaterThan(zero) {
		rsi.avgU = ((rsi.avgU.Mul(decimal.NewFromFloat(float64(rsi.period - 1))).Add(rsi.u)).Div(decimal.NewFromFloat(float64(rsi.period))))
		rsi.avgD = ((rsi.avgD.Mul(decimal.NewFromFloat(float64(rsi.period - 1))).Add(rsi.d)).Div(decimal.NewFromFloat(float64(rsi.period))))
	} else {
		rsi.avgU = u.Div(decimal.NewFromFloat(float64(rsi.period)))
		rsi.avgD = d.Div(decimal.NewFromFloat(float64(rsi.period)))
	}
	rs := rsi.avgU.Div(rsi.avgD)
	rsi.oscillator = (hundred.Sub((hundred.Div((one.Add(rs))))))
	rsi.lastPrice = candle.Close
}
