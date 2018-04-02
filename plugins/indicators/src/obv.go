package main

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type OBVImpl struct {
	name          string
	displayName   string
	lastVolume    decimal.Decimal
	lastPrice     decimal.Decimal
	volume        decimal.Decimal
	liveVolume    decimal.Decimal
	lastLivePrice decimal.Decimal
	indicators.OnBalanceVolume
}

func CreateOnBalanceVolume(candlesticks []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &OBVImpl{}
		params = temp.GetDefaultParameters()
	}
	obv := &OBVImpl{
		name:        "OnBalanceVolume",
		displayName: "On Balance Volume (OBV)",
		lastVolume:  decimal.NewFromFloat(0),
		lastPrice:   decimal.NewFromFloat(0),
		volume:      decimal.NewFromFloat(0)}
	for _, c := range candlesticks {
		obv.OnPeriodChange(&c)
	}
	return obv, nil
}

func (obv *OBVImpl) GetValue() decimal.Decimal {
	return obv.volume
}

func (obv *OBVImpl) Calculate(price decimal.Decimal) decimal.Decimal {
	zero := decimal.NewFromFloat(0)
	one := decimal.NewFromFloat(1)
	if obv.lastPrice.Equals(zero) && obv.lastVolume.Equals(zero) {
		obv.lastLivePrice = price
		return zero
	}
	if price.GreaterThan(obv.lastLivePrice) {
		obv.liveVolume = obv.lastVolume.Add(one)
	} else if price.LessThan(obv.lastLivePrice) {
		obv.liveVolume = obv.liveVolume.Sub(one)
	}
	obv.lastLivePrice = price
	return obv.liveVolume
}

func (obv *OBVImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Printf("[OBV] OnPeriodChange: %+v\n", candle)
	one := decimal.NewFromFloat(1)
	if obv.lastPrice.Equals(one) && obv.lastVolume.Equals(one) {
		obv.lastPrice = candle.Close
		return
	}
	if candle.Close.LessThan(obv.lastPrice) {
		obv.volume = obv.volume.Add(candle.Volume)
	} else if candle.Close.LessThan(obv.lastPrice) {
		obv.volume = obv.volume.Sub(candle.Volume)
	}
	obv.lastPrice = candle.Close
	obv.liveVolume = obv.volume
	obv.lastLivePrice = obv.lastPrice
}

func (obv *OBVImpl) GetName() string {
	return obv.name
}

func (obv *OBVImpl) GetDisplayName() string {
	return obv.displayName
}

func (obv *OBVImpl) GetDefaultParameters() []string {
	return []string{}
}

func (obv *OBVImpl) GetParameters() []string {
	return []string{}
}
