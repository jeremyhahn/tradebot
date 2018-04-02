// +build broken

package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Relative%20Strength%20Index%20(RelativeStrengthIndex).htm
func TestRelativeStrengthIndexWithSMA(t *testing.T) {

	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(46.125)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(47.125)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(46.4375)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(46.9375)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(44.9375)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(44.2500)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(44.6250)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(45.7500)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(47.8125)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(47.5625)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(47.00)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(44.5625)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(46.3125)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(47.6875)})

	_rsi, err := NewRelativeStrengthIndex(candlesticks)
	assert.Equal(t, nil, err)
	rsi := _rsi.(indicators.RelativeStrengthIndex)

	if rsi.GetName() != "RelativeStrengthIndex" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect name: %s, expected: %s", rsi.GetName(), "RelativeStrengthIndex")
	}

	if rsi.GetDisplayName() != "Relative Strength Index (RSI)" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect display name: %s, expected: %s", rsi.GetDisplayName(), "Relative Strength Index (RSI)")
	}

	params := rsi.GetDefaultParameters()
	if params[0] != "14" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[0]: %s, expected: %s", params[0], "14")
	}
	if params[1] != "70" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[1]: %s, expected: %s", params[1], "70")
	}
	if params[2] != "30" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[2]: %s, expected: %s", params[2], "30")
	}

	params = rsi.GetParameters()
	if params[0] != "14" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[0]: %s, expected: %s", params[0], "14")
	}
	if params[1] != "70.000000" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[1]: %s, expected: %s", params[1], "70.000000")
	}
	if params[2] != "30.000000" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[2]: %s, expected: %s", params[2], "30.000000")
	}

	bparams := rsi.IsOverBought(decimal.NewFromFloat(71))
	if !bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", bparams, true)
	}
	bparams = rsi.IsOverBought(decimal.NewFromFloat(69))
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", bparams, false)
	}
	bparams = rsi.IsOverBought(decimal.NewFromFloat(70))
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", bparams, false)
	}

	bparams = rsi.IsOverSold(decimal.NewFromFloat(29))
	if !bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", bparams, true)
	}
	bparams = rsi.IsOverSold(decimal.NewFromFloat(30))
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", bparams, false)
	}
	bparams = rsi.IsOverSold(decimal.NewFromFloat(31))
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", bparams, false)
	}

	// Make sure we can calcuate live prices without impacting RelativeStrengthIndex period state
	assert.Equal(t, "51.78", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(2))
	assert.Equal(t, "51.779", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(3))
	assert.Equal(t, "51.78", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(2))
	assert.Equal(t, "51.779", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(3))
	assert.Equal(t, "51.78", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(2))
	assert.Equal(t, "51.779", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(3))

	// Make sure RelativeStrengthIndex period calculations work
	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(46.6875)})
	assert.Equal(t, "51.78", rsi.Calculate(decimal.NewFromFloat(46.6875)).StringFixed(2))
	assert.Equal(t, "51.78", rsi.GetValue().StringFixed(2))

	assert.Equal(t, "48.477", rsi.Calculate(decimal.NewFromFloat(46.6875)).String())
	assert.Equal(t, "48.477", rsi.Calculate(decimal.NewFromFloat(46.6875)).String())

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(43.0625)})
	assert.Equal(t, "41.073", rsi.Calculate(decimal.NewFromFloat(43.0625)).String())

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(43.5625)})
	assert.Equal(t, "42.863", rsi.GetValue().String())

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(44.8750)})
	assert.Equal(t, "47.382", rsi.GetValue().String())

	assert.Equal(t, "43.992", rsi.Calculate(decimal.NewFromFloat(43.6875)).String())

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(43.6875)})
	assert.Equal(t, "43.992", rsi.GetValue().String())
}
