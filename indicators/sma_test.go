package indicators

import (
	"fmt"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
)

func TestSimpleMovingAverage(t *testing.T) {
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
	fmt.Println(sma.GetAverage())
	u, d := sma.GetGainsAndLosses()

	if u != 3.3399999999999963 {
		t.Errorf("[SMA] Average gains was incorrect")
	}

	if d != 1.3999999999999986 {
		t.Errorf("[SMA] Average losses was incorrect")
	}
}
