package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
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

	sma := NewSimpleMovingAverage(candlesticks)

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

	sma := CreateSimpleMovingAverage(candlesticks, 10)

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
