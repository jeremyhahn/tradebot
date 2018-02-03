package indicators

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
)

type RelativeStrengthIndex interface {
	IsOverSold(rsiValue float64) bool
	IsOverBought(rsiValue float64) bool
	GetValue() float64
	Calculate(price float64) float64
	common.FinancialIndicator
}

type RelativeStrengthIndexParams struct {
	Period     int64
	OverBought float64
	OverSold   float64
}

type RelativeStrengthIndexImpl struct {
	params      *RelativeStrengthIndexParams
	name        string
	displayName string
	ma          common.MovingAverage
	oscillator  float64
	u           float64
	d           float64
	avgU        float64
	avgD        float64
	lastPrice   float64
	RelativeStrengthIndex
}

func NewRelativeStrengthIndex(candles []common.Candlestick) RelativeStrengthIndex {
	params := []string{fmt.Sprintf("%d", len(candles)), "70", "30"}
	return CreateRelativeStrengthIndex(candles, params)
}

func CreateRelativeStrengthIndex(candles []common.Candlestick, params []string) RelativeStrengthIndex {
	period, _ := strconv.ParseInt(params[0], 10, 64)
	overbought, _ := strconv.ParseFloat(params[1], 64)
	oversold, _ := strconv.ParseFloat(params[2], 64)
	ma := CreateSimpleMovingAverage(candles, params)
	candleLen := len(candles)
	lastPrice := 0.0
	if candleLen > 0 {
		lastPrice = candles[candleLen-1].Close
	}
	return &RelativeStrengthIndexImpl{
		name:        "RelativeStrengthIndex",
		displayName: "Relative Strength Index (RSI)",
		ma:          ma,
		oscillator:  0,
		u:           0.0,
		d:           0.0,
		avgU:        0.0,
		avgD:        0.0,
		lastPrice:   lastPrice,
		params: &RelativeStrengthIndexParams{
			Period:     period,
			OverBought: overbought,
			OverSold:   oversold}}
}

func (rsi *RelativeStrengthIndexImpl) Calculate(price float64) float64 {
	var oscillator float64
	curU := rsi.u
	curD := rsi.d
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
		avgU = ((avgU*float64(rsi.params.Period-1) + curU) / float64(rsi.params.Period))
		avgD = ((avgD*float64(rsi.params.Period-1) + curD) / float64(rsi.params.Period))
	} else {
		avgU = u / float64(rsi.params.Period)
		avgD = d / float64(rsi.params.Period)
	}
	rs := avgU / avgD
	oscillator = (100 - (100 / (1 + rs)))
	return oscillator
}

func (rsi *RelativeStrengthIndexImpl) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[RSI] OnPeriodChange: ", candle.Date, candle.Close)
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
		rsi.avgU = ((rsi.avgU*float64(rsi.params.Period-1) + rsi.u) / float64(rsi.params.Period))
		rsi.avgD = ((rsi.avgD*float64(rsi.params.Period-1) + rsi.d) / float64(rsi.params.Period))
	} else {
		rsi.avgU = u / float64(rsi.params.Period)
		rsi.avgD = d / float64(rsi.params.Period)
	}
	rs := rsi.avgU / rsi.avgD
	rsi.oscillator = (100 - (100 / (1 + rs)))
	rsi.lastPrice = candle.Close
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

func (rsi *RelativeStrengthIndexImpl) GetValue() float64 {
	return rsi.oscillator
}

func (rsi *RelativeStrengthIndexImpl) IsOverSold(rsiValue float64) bool {
	return rsiValue < rsi.params.OverSold
}

func (rsi *RelativeStrengthIndexImpl) IsOverBought(rsiValue float64) bool {
	return rsiValue > rsi.params.OverBought
}
