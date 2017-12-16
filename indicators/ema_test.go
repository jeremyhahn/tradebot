package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
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

	ema := NewExponentialMovingAverage(candlesticks)

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

}
