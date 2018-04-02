package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Moving%20Averages.htm
func TestSimpleMovingAverage(t *testing.T) {
	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(67.50)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.50)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.44)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.44)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.25)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(65.88)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.63)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.56)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(65.63)})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(66.06)})

	smaIndicator, err := NewSimpleMovingAverage(candlesticks)
	assert.Equal(t, nil, err)
	sma := smaIndicator.(indicators.SimpleMovingAverage)

	if sma.GetName() != "SimpleMovingAverage" {
		t.Errorf("[SimpleMovingAverage] Got incorrect name: %s, expected: %s", sma.GetName(), "SimpleMovingAverage")
	}

	if sma.GetDisplayName() != "Simple Moving Average (SMA)" {
		t.Errorf("[SimpleMovingAverage] Got incorrect display name: %s, expected: %s", sma.GetDisplayName(), "Simple Moving Average (SMA)")
	}

	params := sma.GetDefaultParameters()
	if params[0] != "20" {
		t.Errorf("[SimpleMovingAverage] Got incorrect default parameter[0]: %s, expected: %s", params[0], "20")
	}

	params = sma.GetParameters()
	if params[0] != "10" {
		t.Errorf("[SimpleMovingAverage] Got incorrect parameter[0]: %s, expected: %s", params[0], "10")
	}

	assert.Equal(t, decimal.NewFromFloat(66.389).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(63.94)})
	assert.Equal(t, decimal.NewFromFloat(66.033).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(64.13)})
	assert.Equal(t, decimal.NewFromFloat(65.796).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(64.50)})
	assert.Equal(t, decimal.NewFromFloat(65.602).String(), sma.GetAverage().String())
}

/*
func TestSimpleMovingAverageUsingAdd(t *testing.T) {
	var candlesticks []common.Candlestick

	smaIndicator, err := CreateSimpleMovingAverage(candlesticks, []string{"10"})
	assert.Equal(t, nil, err)
	sma := smaIndicator.(indicators.SimpleMovingAverage)

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(67.50)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.50)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.44)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.44)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.25)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(65.88)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.63)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.56)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(65.63)})
	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(66.06)})

	assert.Equal(t, decimal.NewFromFloat(66.39).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(63.94)})
	assert.Equal(t, decimal.NewFromFloat(66.03).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(64.13)})
	assert.Equal(t, decimal.NewFromFloat(65.80).String(), sma.GetAverage().String())

	sma.Add(&common.Candlestick{Close: decimal.NewFromFloat(64.50)})
	assert.Equal(t, decimal.NewFromFloat(65.60).String(), sma.GetAverage().String())
}
*/
