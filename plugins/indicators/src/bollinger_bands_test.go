package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/stretchr/testify/assert"
)

func TestBollingerBands(t *testing.T) {
	var candlesticks []common.Candlestick

	candlesticks = append(candlesticks, common.Candlestick{Close: 86.16})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.09})
	candlesticks = append(candlesticks, common.Candlestick{Close: 88.78})
	candlesticks = append(candlesticks, common.Candlestick{Close: 90.32})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.07})
	candlesticks = append(candlesticks, common.Candlestick{Close: 91.15})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.44})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.18})
	candlesticks = append(candlesticks, common.Candlestick{Close: 86.93})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.68})
	candlesticks = append(candlesticks, common.Candlestick{Close: 86.96})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.43})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.32})
	candlesticks = append(candlesticks, common.Candlestick{Close: 88.72})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.45})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.26})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.50})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.90})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.13})
	candlesticks = append(candlesticks, common.Candlestick{Close: 90.70})

	bollingerIndicator, err := NewBollingerBands(candlesticks)
	assert.Equal(t, nil, err)
	bollinger := bollingerIndicator.(indicators.BollingerBands)

	if bollinger.GetName() != "BollingerBands" {
		t.Errorf("[Bollinger] Got incorrect name: %s, expected: %s", bollinger.GetName(), "BollingerBands")
	}

	if bollinger.GetDisplayName() != "Bollinger Bands®" {
		t.Errorf("[Bollinger] Got incorrect display name: %s, expected: %s", bollinger.GetDisplayName(), "Bollinger Bands®")
	}

	params := bollinger.GetDefaultParameters()
	if params[0] != "20" {
		t.Errorf("[Bollinger] Got incorrect default parameter[0]: %s, expected: %s", params[0], "20")
	}
	if params[1] != "2" {
		t.Errorf("[Bollinger] Got incorrect default parameter[1]: %s, expected: %s", params[1], "2")
	}

	params = bollinger.GetParameters()
	if params[0] != "20" {
		t.Errorf("[Bollinger] Got incorrect parameter[0]: %s, expected: %s", params[0], "20")
	}
	if params[1] != "2.000000" {
		t.Errorf("[Bollinger] Got incorrect parameter[1]: %s, expected: %s", params[1], "2.000000")
	}

	actual := bollinger.StandardDeviation()
	expected := 1.29
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect starting standard deviation: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetMiddle()
	expected = 88.71
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect starting middle band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetUpper()
	expected = 91.29
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect starting upper band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetLower()
	expected = 86.13
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect starting lower band: %f, expected: %f", actual, expected)
	}

	bollinger.OnPeriodChange(&common.Candlestick{Close: 92.90})
	actual = bollinger.StandardDeviation()
	expected = 1.45
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect standard deviation: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetMiddle()
	expected = 89.05
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect middle band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetUpper()
	expected = 91.95
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect upper band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetLower()
	expected = 86.14
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect lower band: %f, expected: %f", actual, expected)
	}

	bollinger.OnPeriodChange(&common.Candlestick{Close: 92.98})
	actual = bollinger.StandardDeviation()
	expected = 1.69
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect standard deviation: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetMiddle()
	expected = 89.24
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect middle band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetUpper()
	expected = 92.61
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect upper band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetLower()
	expected = 85.87
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect lower band: %f, expected: %f", actual, expected)
	}

	bollinger.OnPeriodChange(&common.Candlestick{Close: 91.80})
	actual = bollinger.StandardDeviation()
	expected = 1.77
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect standard deviation: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetMiddle()
	expected = 89.39
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect middle band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetUpper()
	expected = 92.93
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect upper band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetLower()
	expected = 85.85
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect lower band: %f, expected: %f", actual, expected)
	}

	bollinger.OnPeriodChange(&common.Candlestick{Close: 92.66})
	actual = bollinger.StandardDeviation()
	expected = 1.90
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect standard deviation: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetMiddle()
	expected = 89.51
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect middle band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetUpper()
	expected = 93.31
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect upper band: %f, expected: %f", actual, expected)
	}
	actual = bollinger.GetLower()
	expected = 85.70
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect lower band: %f, expected: %f", actual, expected)
	}
}

func TestBollingerBands_Calculate(t *testing.T) {
	var candlesticks []common.Candlestick

	candlesticks = append(candlesticks, common.Candlestick{Close: 86.16})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.09})
	candlesticks = append(candlesticks, common.Candlestick{Close: 88.78})
	candlesticks = append(candlesticks, common.Candlestick{Close: 90.32})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.07})
	candlesticks = append(candlesticks, common.Candlestick{Close: 91.15})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.44})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.18})
	candlesticks = append(candlesticks, common.Candlestick{Close: 86.93})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.68})
	candlesticks = append(candlesticks, common.Candlestick{Close: 86.96})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.43})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.32})
	candlesticks = append(candlesticks, common.Candlestick{Close: 88.72})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.45})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.26})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.50})
	candlesticks = append(candlesticks, common.Candlestick{Close: 87.90})
	candlesticks = append(candlesticks, common.Candlestick{Close: 89.13})
	candlesticks = append(candlesticks, common.Candlestick{Close: 90.70})

	bollingerIndicator, err := NewBollingerBands(candlesticks)
	assert.Equal(t, nil, err)
	bollinger := bollingerIndicator.(indicators.BollingerBands)

	upper, middle, lower := bollinger.Calculate(92.90)
	actual := 0.0
	expected := 0.0
	/*
		actual := bollinger.StandardDeviation()
		expected := 1.29
		if actual != expected {
			t.Errorf("[Bollinger] Got incorrect starting standard deviation: %f, expected: %f", actual, expected)
		}*/
	actual = middle
	expected = 89.05
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect middle band: %f, expected: %f", actual, expected)
	}
	actual = upper
	expected = 91.95
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect upper band: %f, expected: %f", actual, expected)
	}
	actual = lower
	expected = 86.14
	if actual != expected {
		t.Errorf("[Bollinger] Got incorrect lower band: %f, expected: %f", actual, expected)
	}
}
