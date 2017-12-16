package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

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

	sma := NewSimpleMovingAverage(candlesticks)
	rsi := NewRelativeStrengthIndex(sma)

	// Make sure we can calcuate live prices without impacting RSI period state
	actual := util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected := 51.779
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = util.RoundFloat(rsi.Calculate(46.6875), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	// Make sure RSI period calculations work
	rsi.OnPeriodChange(&common.Candlestick{Close: 46.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 51.779
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 45.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 48.477
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.0625})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 41.073
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.5625})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 42.863
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 44.8750})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 47.382
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: 43.6875})
	actual = util.RoundFloat(rsi.GetValue(), 3)
	expected = 43.992
	if actual != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

}
