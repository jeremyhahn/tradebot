package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
)

func TestRelativeStrengthIndexWithSMA(t *testing.T) {
	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.34})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.09})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.15})
	candlesticks = append(candlesticks, common.Candlestick{Close: 43.61})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.33})
	candlesticks = append(candlesticks, common.Candlestick{Close: 44.83})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.10})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.42})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.84})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.08})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.89})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.03})
	candlesticks = append(candlesticks, common.Candlestick{Close: 45.61})
	candlesticks = append(candlesticks, common.Candlestick{Close: 46.28})

	sma := NewSimpleMovingAverage(candlesticks)
  rsi := NewRelativeStrengthIndex(sma)

  value := rsi.Calculate(46.28)

  //if value != 70.53 {
  expected := 70.46
  if value != expected {
		t.Errorf("[RSI] Incorrect RSI (SMA) calcuation, got: %f, want: %f.", value, expected)
	}

}

func TestRelativeStrengthIndexWithEMA(t *testing.T) {
	var candlesticks []common.Candlestick
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.27})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.19})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.08})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.17})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.18})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.13})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.23})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.43})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.24})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.29})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.15})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.39})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.38})
	candlesticks = append(candlesticks, common.Candlestick{Close: 22.61})

	ema := NewExponentialMovingAverage(candlesticks)
  rsi := NewRelativeStrengthIndex(ema)

  value := rsi.Calculate(46.28)

  expected := 101.81
  if value != expected {
		t.Errorf("[RSI] Incorrect RSI (EMA) calcuation, got: %f, want: %f.", value, expected)
	}

}
