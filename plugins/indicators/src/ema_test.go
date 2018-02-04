package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/util"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Moving%20Averages.htm
func TestExponentialMovingAverage(t *testing.T) {

	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: 64.75})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.79})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.73})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.73})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.55})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.19})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.91})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.85})
	candlesticks = append(candlesticks, common.Candlestick{Close: 62.95})
	candlesticks = append(candlesticks, common.Candlestick{Close: 63.37})

	ema := NewExponentialMovingAverage(candlesticks).(indicators.ExponentialMovingAverage)

	if ema.GetName() != "ExponentialMovingAverage" {
		t.Errorf("[ExponentialMovingAverage] Got incorrect name: %s, expected: %s", ema.GetName(), "ExponentialMovingAverage")
	}

	if ema.GetDisplayName() != "Exponential Moving Average (EMA)" {
		t.Errorf("[ExponentialMovingAverage] Got incorrect display name: %s, expected: %s", ema.GetDisplayName(), "Exponential Moving Average (EMA)")
	}

	params := ema.GetDefaultParameters()
	if params[0] != "20" {
		t.Errorf("[ExponentialMovingAverage] Got incorrect default parameter[0]: %s, expected: %s", params[0], "20")
	}

	params = ema.GetParameters()
	if params[0] != "10" {
		t.Errorf("[ExponentialMovingAverage] Got incorrect parameter[0]: %s, expected: %s", params[0], "10")
	}

	/*
		actual := ema.GetMultiplier()
		expected := 0.181818
		if actual != expected {
			t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
		}*/

	actual := ema.GetAverage()
	expected := 63.682
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: 61.33})
	actual = util.RoundFloat(ema.GetAverage(), 3)
	expected = 63.254
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: 61.51})
	actual = util.RoundFloat(ema.GetAverage(), 3)
	expected = 62.937
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: 61.87})
	actual = util.RoundFloat(ema.GetAverage(), 3)
	expected = 62.743
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: 60.25})
	actual = util.RoundFloat(ema.GetAverage(), 3)
	expected = 62.290
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

	ema.Add(&common.Candlestick{Close: 59.35})
	actual = util.RoundFloat(ema.GetAverage(), 3)
	expected = 61.755
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}

}

/*
func TestExponentialMovingAverageUsingAdd(t *testing.T) {
	var candlesticks []common.Candlestick

	ema := CreateExponentialMovingAverage(candlesticks, 10)
	ema.Add(&common.Candlestick{Close: 64.75})
	ema.Add(&common.Candlestick{Close: 63.79})
	ema.Add(&common.Candlestick{Close: 63.73})
	ema.Add(&common.Candlestick{Close: 63.73})
	ema.Add(&common.Candlestick{Close: 63.55})
	ema.Add(&common.Candlestick{Close: 63.19})
	ema.Add(&common.Candlestick{Close: 63.91})
	ema.Add(&common.Candlestick{Close: 63.85})
	ema.Add(&common.Candlestick{Close: 62.95})
	ema.Add(&common.Candlestick{Close: 63.37})

	actual := ema.GetAverage()
	expected := 63.682
	if actual != expected {
		t.Errorf("[EMA] Got incorrect average: %f, expected: %f", actual, expected)
	}
}
*/
