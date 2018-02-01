package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

// http://investexcel.net/how-to-calculate-macd-in-excel/
func TestMovingAverageConvergenceDivergence(t *testing.T) {
	var candles []common.Candlestick
	// 12 day EMA
	candles = append(candles, common.Candlestick{Close: 459.99})
	candles = append(candles, common.Candlestick{Close: 448.85})
	candles = append(candles, common.Candlestick{Close: 446.06})
	candles = append(candles, common.Candlestick{Close: 450.81})
	candles = append(candles, common.Candlestick{Close: 442.80})
	candles = append(candles, common.Candlestick{Close: 448.97})
	candles = append(candles, common.Candlestick{Close: 444.57})
	candles = append(candles, common.Candlestick{Close: 441.4})
	candles = append(candles, common.Candlestick{Close: 430.47})
	candles = append(candles, common.Candlestick{Close: 420.05})
	candles = append(candles, common.Candlestick{Close: 431.14})
	candles = append(candles, common.Candlestick{Close: 425.66})
	// 26 day EMA
	candles = append(candles, common.Candlestick{Close: 430.58})
	candles = append(candles, common.Candlestick{Close: 431.72})
	candles = append(candles, common.Candlestick{Close: 437.87})
	candles = append(candles, common.Candlestick{Close: 428.43})
	candles = append(candles, common.Candlestick{Close: 428.35})
	candles = append(candles, common.Candlestick{Close: 432.50})
	candles = append(candles, common.Candlestick{Close: 443.66})
	candles = append(candles, common.Candlestick{Close: 455.72})
	candles = append(candles, common.Candlestick{Close: 454.49})
	candles = append(candles, common.Candlestick{Close: 452.08})
	candles = append(candles, common.Candlestick{Close: 452.73})
	candles = append(candles, common.Candlestick{Close: 461.91})
	candles = append(candles, common.Candlestick{Close: 463.58})
	candles = append(candles, common.Candlestick{Close: 461.14})

	macd := NewMovingAverageConvergenceDivergence(candles)

	if macd.GetName() != "MovingAverageConvergenceDivergence" {
		t.Errorf("[MovingAverageConvergenceDivergence] Got incorrect name: %s, expected: %s", macd.GetName(), "MovingAverageConvergenceDivergence")
	}

	if macd.GetDisplayName() != "Moving Average Convergence Divergence (MACD)" {
		t.Errorf("[MovingAverageConvergenceDivergence] Got incorrect display name: %s, expected: %s", macd.GetDisplayName(), "Moving Average Convergence Divergence (MACD)")
	}

	params := macd.GetDefaultParameters()
	if params[0] != "12" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[0]: %s, expected: %s", params[0], "12")
	}
	if params[1] != "26" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[1]: %s, expected: %s", params[1], "26")
	}
	if params[2] != "9" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect default parameter[2]: %s, expected: %s", params[2], "9")
	}

	params = macd.GetParameters()
	if params[0] != "12" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[0]: %s, expected: %s", params[0], "12")
	}
	if params[1] != "26" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[1]: %s, expected: %s", params[1], "26")
	}
	if params[2] != "9" {
		t.Errorf("[RelativeStrengthIndex] Got incorrect parameter[2]: %s, expected: %s", params[2], "9")
	}

	macd1 := util.RoundFloat(macd.GetValue(), 6)
	expected := 8.275270
	if macd1 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd1, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 452.08})
	macd2 := util.RoundFloat(macd.GetValue(), 6)
	expected = 7.703378
	if macd2 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd2, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 442.66})
	macd3 := util.RoundFloat(macd.GetValue(), 6)
	expected = 6.416075
	if macd3 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd3, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 428.91})
	macd4 := util.RoundFloat(macd.GetValue(), 6)
	expected = 4.23752
	if macd4 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd4, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 429.79})
	macd5 := util.RoundFloat(macd.GetValue(), 6)
	expected = 2.552583
	if macd5 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd5, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 431.99})
	macd6 := util.RoundFloat(macd.GetValue(), 6)
	expected = 1.378886
	if macd6 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd6, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 427.72})
	macd7 := util.RoundFloat(macd.GetValue(), 6)
	expected = 0.102981
	if macd7 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd7, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 423.2})
	macd8 := util.RoundFloat(macd.GetValue(), 4)
	expected = -1.2584
	if macd8 != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd8, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 426.21})
	macd9 := util.RoundFloat(macd.GetValue(), 6)
	expected = -2.070558
	if macd9 != expected {
		t.Errorf("[MACD] Got incorrect value: %f, expected: %f", macd9, expected)
	}

	actual := util.RoundFloat(macd.GetSignalLine(), 6)
	expected = 3.037526
	actualHistogram := util.RoundFloat(macd.GetHistogram(), 6)
	expectedHistogram := -5.108084
	if actual != expected {
		t.Errorf("[MACD] Got incorrect signal line: %f, expected: %f", actual, expected)
	}
	if actualHistogram != expectedHistogram {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 426.98})
	actual = util.RoundFloat(macd.GetSignalLine(), 6)
	expected = 1.905652
	actualHistogram = util.RoundFloat(macd.GetHistogram(), 6)
	expectedHistogram = -4.527495
	if actual != expected {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", actual, expected)
	}
	if actualHistogram != expectedHistogram {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: 435.69})
	actual = util.RoundFloat(macd.GetSignalLine(), 6)
	expected = 1.058708
	actualHistogram = util.RoundFloat(macd.GetHistogram(), 6)
	expectedHistogram = -3.387775
	if actual != expected {
		t.Errorf("[MACD] Got incorrect signal line: %f, expected: %f", actual, expected)
	}
	if actualHistogram != expectedHistogram {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

}
