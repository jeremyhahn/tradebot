package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
)

// https://www.investopedia.com/terms/o/onbalancevolume.asp
func TestOnBalanceVolume(t *testing.T) {

	var candlesticks []common.Candlestick
	obv := CreateOnBalanceVolume(candlesticks, nil).(indicators.OnBalanceVolume)

	actual := obv.GetValue()
	expected := 0.0
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.00, Volume: 25200})
	actual = obv.GetValue()
	expected = 0.0
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.15, Volume: 30000})
	actual = obv.GetValue()
	expected = 30000
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.17, Volume: 25600})
	actual = obv.GetValue()
	expected = 55600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.13, Volume: 32000})
	actual = obv.GetValue()
	expected = 23600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.11, Volume: 23000})
	actual = obv.GetValue()
	expected = 600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.15, Volume: 40000})
	actual = obv.GetValue()
	expected = 40600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.20, Volume: 36000})
	actual = obv.GetValue()
	expected = 76600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.20, Volume: 20500})
	actual = obv.GetValue()
	expected = 76600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.22, Volume: 23000})
	actual = obv.GetValue()
	expected = 99600
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}

	obv.OnPeriodChange(&common.Candlestick{Close: 10.21, Volume: 27500})
	actual = obv.GetValue()
	expected = 72100
	if actual != expected {
		t.Errorf("[OBV] Incorrect OBV calcuation, got: %f, want: %f.", actual, expected)
	}
}
