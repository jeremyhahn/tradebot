package indicators

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type OBV struct {
	lastVolume    float64
	lastPrice     float64
	volume        float64
	liveVolume    float64
	lastLivePrice float64
}

func NewOnBalanceVolume(candlesticks []common.Candlestick) *OBV {
	obv := &OBV{
		lastVolume: 0.0,
		lastPrice:  0.0,
		volume:     0.0}
	for _, c := range candlesticks {
		obv.OnPeriodChange(&c)
	}
	return obv
}

func (obv *OBV) GetValue() float64 {
	return obv.volume
}

func (obv *OBV) Calculate(price float64) float64 {
	if obv.lastPrice == 0 && obv.lastVolume == 0 {
		obv.lastLivePrice = price
		return 0.0
	}
	if price > obv.lastLivePrice {
		obv.liveVolume += 1
	} else if price < obv.lastLivePrice {
		obv.liveVolume -= 1
	}
	obv.lastLivePrice = price
	return obv.liveVolume
}

func (obv *OBV) OnPeriodChange(candle *common.Candlestick) {
	fmt.Printf("[OBV] OnPeriodChange: %+v\n", candle)
	if obv.lastPrice == 0 && obv.lastVolume == 0 {
		obv.lastPrice = candle.Close
		return
	}
	if candle.Close > obv.lastPrice {
		obv.volume += candle.Volume
	} else if candle.Close < obv.lastPrice {
		obv.volume -= candle.Volume
	}
	obv.lastPrice = candle.Close
	obv.liveVolume = 0
	obv.lastLivePrice = 0
}
