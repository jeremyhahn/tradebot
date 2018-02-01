package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

// http://cns.bu.edu/~gsc/CN710/fincast/Technical%20_indicators/Relative%20Strength%20Index%20(RelativeStrengthIndex).htm
func TestRelativeStrengthIndexWithSMA(t *testing.T) {

	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.125})
	candlesticks = append(candlesticks, common.Candlestick{Close: 47.125})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.4375})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.9375})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.9375})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.2500})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.6250})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.7500})
	candlesticks = append(candlesticks, common.Candlestick{Close: 47.8125})
	candlesticks = append(candlesticks, common.Candlestick{Close: 47.5625})
	candlesticks = append(candlesticks, common.Candlestick{Close: 47.00})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.5625})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.3125})
	candlesticks = append(candlesticks, common.Candlestick{Close: 47.6875})

	rsi := NewRelativeStrengthIndex(candlesticks)

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

	bparams := rsi.IsOverBought(71)
	if !bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", params[0], true)
	}
	bparams = rsi.IsOverBought(69)
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", params[0], false)
	}
	bparams = rsi.IsOverBought(70)
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over bought signal[0]: %t, expected: %t", params[0], false)
	}

	bparams = rsi.IsOverSold(29)
	if !bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", params[0], true)
	}
	bparams = rsi.IsOverSold(30)
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", params[0], false)
	}
	bparams = rsi.IsOverSold(31)
	if bparams {
		t.Errorf("[RelativeStrengthIndex] Got incorrect over sold signal[0]: %t, expected: %t", params[0], false)
	}

	// Make sure we can calcuate live prices without impacting RelativeStrengthIndex period state
	actual := util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected := 51.779
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	// Make sure RelativeStrengthIndex period calculations work
	rsi.OnPeriodChange(&common.Candlestick{Close: 46.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(45.6875), 3)
	expected = 48.477
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 45.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 48.477
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.0625})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 41.073
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.5625})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 42.863
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 44.8750})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 47.382
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(43.6875), 3)
	expected = 43.992
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 43.992
	if actual != expected {
		t.Errorf("[RelativeStrengthIndex] Incorrect RelativeStrengthIndex (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

}
