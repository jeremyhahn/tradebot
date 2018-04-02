package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type RelativeStrengthIndexParams struct {
	Period     int64
	OverBought float64
	OverSold   float64
}

type RelativeStrengthIndexImpl struct {
	params      *RelativeStrengthIndexParams
	name        string
	displayName string
	sma         indicators.SimpleMovingAverage
	oscillator  decimal.Decimal
	u           decimal.Decimal
	d           decimal.Decimal
	avgU        decimal.Decimal
	avgD        decimal.Decimal
	lastPrice   decimal.Decimal
	indicators.RelativeStrengthIndex
}

func NewRelativeStrengthIndex(candles []common.Candlestick) (common.FinancialIndicator, error) {
	params := []string{fmt.Sprintf("%d", len(candles)), "70", "30"}
	return CreateRelativeStrengthIndex(candles, params)
}

func CreateRelativeStrengthIndex(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &RelativeStrengthIndexImpl{}
		params = temp.GetDefaultParameters()
	}
	period, _ := strconv.ParseInt(params[0], 10, 64)
	overbought, _ := strconv.ParseFloat(params[1], 64)
	oversold, _ := strconv.ParseFloat(params[2], 64)
	candleLen := len(candles)
	if candleLen < int(period) {
		return nil, errors.New(fmt.Sprintf(
			"RelativeStrengthIndex requires candlestick length of %d, received %d", int(period), candleLen))
	}
	smaIndicator, err := CreateSimpleMovingAverage(candles[:period], []string{params[0]})
	if err != nil {
		return nil, err
	}
	sma := smaIndicator.(indicators.SimpleMovingAverage)
	var lastPrice decimal.Decimal
	if candleLen > 0 {
		lastPrice = candles[candleLen-1].Close
	}
	rsi := &RelativeStrengthIndexImpl{
		name:        "RelativeStrengthIndex",
		displayName: "Relative Strength Index (RSI)",
		sma:         sma,
		oscillator:  decimal.NewFromFloat(0),
		u:           decimal.NewFromFloat(0),
		d:           decimal.NewFromFloat(0),
		avgU:        decimal.NewFromFloat(0),
		avgD:        decimal.NewFromFloat(0),
		lastPrice:   lastPrice,
		params: &RelativeStrengthIndexParams{
			Period:     period,
			OverBought: overbought,
			OverSold:   oversold}}

	for _, c := range candles[period:] {
		rsi.OnPeriodChange(&c)
	}

	return rsi, nil
}

func (rsi *RelativeStrengthIndexImpl) Calculate(price decimal.Decimal) decimal.Decimal {
	var oscillator decimal.Decimal
	curU := rsi.u
	curD := rsi.d
	avgU := rsi.avgU
	avgD := rsi.avgD
	u, d := rsi.sma.GetGainsAndLosses()
	zero := decimal.NewFromFloat(0)
	difference := price.Sub(rsi.lastPrice)
	if difference.LessThan(zero) {
		d = d.Add(difference.Abs())
		curD = difference.Abs()
		curU = decimal.NewFromFloat(0)
	} else {
		u = u.Add(difference)
		curU = difference
		curD = decimal.NewFromFloat(0)
	}
	decPeriod := decimal.NewFromFloat(float64(rsi.params.Period))
	one := decimal.NewFromFloat(1)
	if avgU.GreaterThan(zero) && avgD.GreaterThan(zero) {
		//avgU = ((avgU*float64(rsi.params.Period-1) + curU) / float64(rsi.params.Period))
		//avgD = ((avgD*float64(rsi.params.Period-1) + curD) / float64(rsi.params.Period))
		avgU = avgU.Mul(decPeriod.Sub(one).Add(curU)).Div(decPeriod)
		avgD = avgD.Mul(decPeriod.Sub(one).Add(curD)).Div(decPeriod)
	} else {
		//avgU = u / float64(rsi.params.Period)
		//avgD = d / float64(rsi.params.Period)
		avgU = u.Div(decPeriod)
		avgD = d.Div(decPeriod)
	}
	rs := avgU.Div(avgD)
	//oscillator = (100 - (100 / (1 + rs)))
	hundred := decimal.NewFromFloat(100)
	oscillator = hundred.Sub(hundred.Div(one.Add(rs)))
	return oscillator
}

func (rsi *RelativeStrengthIndexImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[RSI] OnPeriodChange: %s", candle.ToString())
	rsi.sma.Add(candle)
	u, d := rsi.sma.GetGainsAndLosses()
	difference := candle.Close.Sub(rsi.lastPrice)
	zero := decimal.NewFromFloat(0)
	one := decimal.NewFromFloat(1)
	if difference.LessThan(zero) {
		d = d.Add(difference.Abs())
		rsi.d = difference.Abs()
		rsi.u = decimal.NewFromFloat(0)
	} else {
		u = u.Add(difference)
		rsi.u = difference
		rsi.d = decimal.NewFromFloat(0)
	}
	decPeriod := decimal.NewFromFloat(float64(rsi.params.Period))
	if rsi.avgU.GreaterThan(zero) && rsi.avgD.GreaterThan(zero) {
		//rsi.avgU = ((rsi.avgU*float64(rsi.params.Period-1) + rsi.u) / float64(rsi.params.Period))
		//rsi.avgD = ((rsi.avgD*float64(rsi.params.Period-1) + rsi.d) / float64(rsi.params.Period))
		rsi.avgU = rsi.avgU.Mul(decPeriod.Sub(one).Add(rsi.u)).Div(decPeriod)
		rsi.avgD = rsi.avgD.Mul(decPeriod.Sub(one).Add(rsi.d)).Div(decPeriod)
	} else {
		//rsi.avgU = u / float64(rsi.params.Period)
		//rsi.avgD = d / float64(rsi.params.Period)
		rsi.avgU = u.Div(decPeriod)
		rsi.avgD = d.Div(decPeriod)
	}
	//rs := rsi.avgU / rsi.avgD
	rs := rsi.avgU.Div(rsi.avgD)
	rsi.lastPrice = candle.Close
	//rsi.oscillator = (100 - (100 / (1 + rs)))
	hundred := decimal.NewFromFloat(100)
	rsi.oscillator = hundred.Sub(hundred.Div(one.Add(rs)))
}

func (rsi *RelativeStrengthIndexImpl) GetName() string {
	return rsi.name
}

func (rsi *RelativeStrengthIndexImpl) GetDisplayName() string {
	return rsi.displayName
}

func (rsi *RelativeStrengthIndexImpl) GetDefaultParameters() []string {
	return []string{"14", "70", "30"}
}

func (rsi *RelativeStrengthIndexImpl) GetParameters() []string {
	return []string{
		fmt.Sprintf("%d", rsi.params.Period),
		fmt.Sprintf("%f", rsi.params.OverBought),
		fmt.Sprintf("%f", rsi.params.OverSold)}
}

func (rsi *RelativeStrengthIndexImpl) GetValue() decimal.Decimal {
	return rsi.oscillator
}

func (rsi *RelativeStrengthIndexImpl) IsOverSold(rsiValue decimal.Decimal) bool {
	return rsiValue.LessThan(decimal.NewFromFloat(rsi.params.OverSold))
}

func (rsi *RelativeStrengthIndexImpl) IsOverBought(rsiValue decimal.Decimal) bool {
	return rsiValue.GreaterThan(decimal.NewFromFloat(rsi.params.OverBought))
}
