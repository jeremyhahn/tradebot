package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/util"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Moving%20Averages.htm
func TestSimpleMovingAverage(t *testing.T) {
	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: 67.50})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.50})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.44})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.44})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.25})
	candlesticks = append(candlesticks, common.Candlestick{Close: 65.88})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.63})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.56})
	candlesticks = append(candlesticks, common.Candlestick{Close: 65.63})
	candlesticks = append(candlesticks, common.Candlestick{Close: 66.06})

	sma := NewSimpleMovingAverage(candlesticks).(indicators.SimpleMovingAverage)

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

	actual := util.RoundFloat(sma.GetAverage(), 2)
	expected := 66.39
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average, got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 63.94})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 66.03
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 64.13})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 65.80
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 64.50})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 65.60
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}
}

func TestSimpleMovingAverageUsingAdd(t *testing.T) {
	var candlesticks []common.Candlestick

	params := []string{"10"}
	sma := CreateSimpleMovingAverage(candlesticks, params).(indicators.SimpleMovingAverage)

	sma.Add(&common.Candlestick{Close: 67.50})
	sma.Add(&common.Candlestick{Close: 66.50})
	sma.Add(&common.Candlestick{Close: 66.44})
	sma.Add(&common.Candlestick{Close: 66.44})
	sma.Add(&common.Candlestick{Close: 66.25})
	sma.Add(&common.Candlestick{Close: 65.88})
	sma.Add(&common.Candlestick{Close: 66.63})
	sma.Add(&common.Candlestick{Close: 66.56})
	sma.Add(&common.Candlestick{Close: 65.63})
	sma.Add(&common.Candlestick{Close: 66.06})

	actual := util.RoundFloat(sma.GetAverage(), 2)
	expected := 66.39
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average, got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 63.94})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 66.03
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 64.13})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 65.80
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}

	sma.Add(&common.Candlestick{Close: 64.50})
	actual = util.RoundFloat(sma.GetAverage(), 2)
	expected = 65.60
	if actual != expected {
		t.Errorf("[SMA] Got incorrect average after Add(), got %f, expected :%f", actual, expected)
	}
}
