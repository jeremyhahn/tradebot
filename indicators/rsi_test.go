package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

func TestRelativeStrengthIndexWithSMA(t *testing.T) {

	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(46.125))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(47.125))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(46.4375))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(46.9375))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(44.9375))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(44.2500))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(44.6250))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(45.7500))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(47.8125))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(47.5625))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(47.00))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(44.5625))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(46.3125))})
	candlesticks = append(candlesticks, common.Candlestick{Close: decimal.NewFromFloat(float64(47.6875))})

	sma := NewSimpleMovingAverage(candlesticks)
	rsi := NewRelativeStrengthIndex(sma)

	// Make sure we can calcuate live prices without impacting RSI period state
	actual := rsi.Calculate(decimal.NewFromFloat(float64(46.6875)))
	expected := decimal.NewFromFloat(float64(51.779))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = rsi.Calculate(decimal.NewFromFloat(float64(46.6875)))
	expected = decimal.NewFromFloat(float64(51.779))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	actual = rsi.Calculate(decimal.NewFromFloat(float64(46.6875)))
	expected = decimal.NewFromFloat(float64(51.779))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	// Make sure RSI period calculations work
	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(46.6875))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(51.779))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(45.6875))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(48.477))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(43.0625))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(41.073))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(43.5625))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(42.863))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(44.8750))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(47.382))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

	rsi.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(43.6875))})
	actual = rsi.GetValue()
	expected = decimal.NewFromFloat(float64(43.992))
	if !actual.Equals(expected) {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", actual, expected)
	}

}
