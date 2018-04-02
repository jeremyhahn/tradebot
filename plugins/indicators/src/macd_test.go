// +build broken

package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// http://investexcel.net/how-to-calculate-macd-in-excel/
func TestMovingAverageConvergenceDivergence(t *testing.T) {
	var candles []common.Candlestick
	// 12 day EMA
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(459.99)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(448.85)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(446.06)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(450.81)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(442.80)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(448.97)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(444.57)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(441.4)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(430.47)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(420.05)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(431.14)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(425.66)})
	// 26 day EMA
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(430.58)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(431.72)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(437.87)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(428.43)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(428.35)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(432.50)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(443.66)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(455.72)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(454.49)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(452.08)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(452.73)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(461.91)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(463.58)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(461.14)})

	macdIndicator, err := NewMovingAverageConvergenceDivergence(candles)
	assert.Equal(t, nil, err)
	macd := macdIndicator.(indicators.MovingAverageConvergenceDivergence)

	if macd.GetName() != "MovingAverageConvergenceDivergence" {
		t.Errorf("[MovingAverageConvergenceDivergence] Got incorrect name: %s, expected: %s", macd.GetName(), "MovingAverageConvergenceDivergence")
	}

	if macd.GetDisplayName() != "Moving Average Convergence Divergence (MACD)" {
		t.Errorf("[MovingAverageConvergenceDivergence] Got incorrect display name: %s, expected: %s", macd.GetDisplayName(), "Moving Average Convergence Divergence (MACD)")
	}

	params := macd.GetDefaultParameters()
	if params[0] != "12" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[0]: %s, expected: %s", params[0], "12")
	}
	if params[1] != "26" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[1]: %s, expected: %s", params[1], "26")
	}
	if params[2] != "9" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[2]: %s, expected: %s", params[2], "9")
	}

	params = macd.GetParameters()
	if params[0] != "12" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[0]: %s, expected: %s", params[0], "12")
	}
	if params[1] != "26" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[1]: %s, expected: %s", params[1], "26")
	}
	if params[2] != "9" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[2]: %s, expected: %s", params[2], "9")
	}

	assert.Equal(t, decimal.NewFromFloat(8.275270).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(452.08)})
	assert.Equal(t, decimal.NewFromFloat(7.703378).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(442.66)})
	assert.Equal(t, decimal.NewFromFloat(6.416075).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(428.91)})
	assert.Equal(t, decimal.NewFromFloat(4.23752).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(429.79)})
	assert.Equal(t, decimal.NewFromFloat(2.552583).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(431.99)})
	assert.Equal(t, decimal.NewFromFloat(1.378886).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(427.72)})
	assert.Equal(t, decimal.NewFromFloat(0.102981).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(423.2)})
	assert.Equal(t, decimal.NewFromFloat(-1.2584).String(), macd.GetValue().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(426.21)})
	assert.Equal(t, decimal.NewFromFloat(-2.070558).String(), macd.GetValue().String())

	// Signal line & histogram
	assert.Equal(t, decimal.NewFromFloat(3.037526).String(), macd.GetSignalLine().String())
	assert.Equal(t, decimal.NewFromFloat(-5.108084).String(), macd.GetHistogram().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(426.98)})
	assert.Equal(t, decimal.NewFromFloat(1.905652).String(), macd.GetSignalLine().String())
	assert.Equal(t, decimal.NewFromFloat(-4.527495).String(), macd.GetHistogram().String())

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(435.69)})
	assert.Equal(t, decimal.NewFromFloat(1.058708).String(), macd.GetSignalLine().String())
	assert.Equal(t, decimal.NewFromFloat(-3.387775).String(), macd.GetHistogram().String())
}
