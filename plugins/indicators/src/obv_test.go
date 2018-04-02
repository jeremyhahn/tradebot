// +build broken

package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// https://www.investopedia.com/terms/o/onbalancevolume.asp
func TestOnBalanceVolume(t *testing.T) {

	var candlesticks []common.Candlestick

	obvIndicator, err := CreateOnBalanceVolume(candlesticks, nil)
	assert.Equal(t, nil, err)
	obv := obvIndicator.(indicators.OnBalanceVolume)

	assert.Equal(t, "0", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.00), Volume: decimal.NewFromFloat(25200)})
	assert.Equal(t, "0", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.15), Volume: decimal.NewFromFloat(30000)})
	assert.Equal(t, "30000", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.17), Volume: decimal.NewFromFloat(25600)})
	assert.Equal(t, "55600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.13), Volume: decimal.NewFromFloat(32000)})
	assert.Equal(t, "23600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.11), Volume: decimal.NewFromFloat(23000)})
	assert.Equal(t, "600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.15), Volume: decimal.NewFromFloat(40000)})
	assert.Equal(t, "40600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.20), Volume: decimal.NewFromFloat(36000)})
	assert.Equal(t, "76600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.20), Volume: decimal.NewFromFloat(20500)})
	assert.Equal(t, "76600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.22), Volume: decimal.NewFromFloat(23000)})
	assert.Equal(t, "99600", obv.GetValue().String())

	obv.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(10.21), Volume: decimal.NewFromFloat(27500)})
	assert.Equal(t, "72100", obv.GetValue().String())
}
