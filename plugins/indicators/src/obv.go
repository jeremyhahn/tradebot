package main

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
)

type OBVImpl struct {
	lastVolume    float64
	lastPrice     float64
	volume        float64
	liveVolume    float64
	lastLivePrice float64
	indicators.OBV
}

func NewOnBalanceVolume(candlesticks []common.Candlestick) indicators.OBV {
	obv := &OBVImpl{
		lastVolume: 0.0,
		lastPrice:  0.0,
		volume:     0.0}
	for _, c := range candlesticks {
		obv.OnPeriodChange(&c)
	}
	return obv
}

func (obv *OBVImpl) GetValue() float64 {
	return obv.volume
}

func (obv *OBVImpl) Calculate(price float64) float64 {
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

func (obv *OBVImpl) OnPeriodChange(candle *common.Candlestick) {
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
	obv.liveVolume = obv.volume
	obv.lastLivePrice = obv.lastPrice
}
