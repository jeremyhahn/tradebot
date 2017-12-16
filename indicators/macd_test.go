package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

// http://investexcel.net/how-to-calculate-macd-in-excel/
func TestMovingAverageConvergenceDivergence(t *testing.T) {
	var ema1candles, ema2candles []common.Candlestick

	// 12 day EMA
	ema1candles = append(ema1candles, common.Candlestick{Close: 459.99})
	ema1candles = append(ema1candles, common.Candlestick{Close: 448.85})
	ema1candles = append(ema1candles, common.Candlestick{Close: 446.06})
	ema1candles = append(ema1candles, common.Candlestick{Close: 450.81})
	ema1candles = append(ema1candles, common.Candlestick{Close: 442.80})
	ema1candles = append(ema1candles, common.Candlestick{Close: 448.97})
	ema1candles = append(ema1candles, common.Candlestick{Close: 444.57})
	ema1candles = append(ema1candles, common.Candlestick{Close: 441.4})
	ema1candles = append(ema1candles, common.Candlestick{Close: 430.47})
	ema1candles = append(ema1candles, common.Candlestick{Close: 420.05})
	ema1candles = append(ema1candles, common.Candlestick{Close: 431.14})
	ema1candles = append(ema1candles, common.Candlestick{Close: 425.66})

	// 26 day EMA
	ema2candles = ema1candles
	ema2candles = append(ema2candles, common.Candlestick{Close: 430.58})
	ema2candles = append(ema2candles, common.Candlestick{Close: 431.72})
	ema2candles = append(ema2candles, common.Candlestick{Close: 437.87})
	ema2candles = append(ema2candles, common.Candlestick{Close: 428.43})
	ema2candles = append(ema2candles, common.Candlestick{Close: 428.35})
	ema2candles = append(ema2candles, common.Candlestick{Close: 432.50})
	ema2candles = append(ema2candles, common.Candlestick{Close: 443.66})
	ema2candles = append(ema2candles, common.Candlestick{Close: 455.72})
	ema2candles = append(ema2candles, common.Candlestick{Close: 454.49})
	ema2candles = append(ema2candles, common.Candlestick{Close: 452.08})
	ema2candles = append(ema2candles, common.Candlestick{Close: 452.73})
	ema2candles = append(ema2candles, common.Candlestick{Close: 461.91})
	ema2candles = append(ema2candles, common.Candlestick{Close: 463.58})
	ema2candles = append(ema2candles, common.Candlestick{Close: 461.14})

	ema1 := NewExponentialMovingAverage(ema1candles)
	ema2 := NewExponentialMovingAverage(ema2candles)

	// EMA1 tests
	actual := ema1.GetAverage()
	expected := 440.8975
	if util.RoundFloat(actual, 4) != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: 430.58})
	actual = util.RoundFloat(ema1.GetAverage(), 6)
	expected = 439.310192
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: 431.72})
	actual = util.RoundFloat(ema1.GetAverage(), 6)
	expected = 438.142470
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: 437.87})
	actual = util.RoundFloat(ema1.GetAverage(), 6)
	expected = 438.100552
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: 428.43})
	actual = util.RoundFloat(ema1.GetAverage(), 6)
	expected = 436.612775
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: 428.35})
	actual = util.RoundFloat(ema1.GetAverage(), 6)
	expected = 435.341579
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	// EMA2 tests
	actual = util.RoundFloat(ema2.GetAverage(), 6)
	expected = 443.289615
	if actual != expected {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	// Add data difference between EMA2 and EMA1 to EMA1
	ema1.Add(&common.Candlestick{Close: 432.5})
	ema1.Add(&common.Candlestick{Close: 443.66})
	ema1.Add(&common.Candlestick{Close: 455.72})
	ema1.Add(&common.Candlestick{Close: 454.49})
	ema1.Add(&common.Candlestick{Close: 452.08})
	ema1.Add(&common.Candlestick{Close: 452.73})
	ema1.Add(&common.Candlestick{Close: 461.91})
	ema1.Add(&common.Candlestick{Close: 463.58})
	ema1.Add(&common.Candlestick{Close: 461.14})

	macd := NewMovingAverageConvergenceDivergence(ema1, ema2, 9)

	macd1 := util.RoundFloat(macd.GetValue(), 6)
	expected = 8.275270
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

	actual = util.RoundFloat(macd.GetSignalLine(), 6)
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
