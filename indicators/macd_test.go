package indicators

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

// http://investexcel.net/how-to-calculate-macd-in-excel/
func TestMovingAverageConvergenceDivergence(t *testing.T) {
	var ema1candles, ema2candles []common.Candlestick

	// 12 day EMA
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(459.99))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(448.85))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(446.06))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(450.81))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(442.80))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(448.97))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(444.57))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(441.4))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(430.47))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(420.05))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(431.14))})
	ema1candles = append(ema1candles, common.Candlestick{Close: decimal.NewFromFloat(float64(425.66))})

	// 26 day EMA
	ema2candles = ema1candles
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(430.58))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(431.72))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(437.87))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(428.43))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(428.35))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(432.50))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(443.66))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(455.72))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(454.49))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(452.08))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(452.73))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(461.91))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(463.58))})
	ema2candles = append(ema2candles, common.Candlestick{Close: decimal.NewFromFloat(float64(461.14))})

	ema1 := NewExponentialMovingAverage(ema1candles)
	ema2 := NewExponentialMovingAverage(ema2candles)

	// EMA1 tests
	actual := ema1.GetAverage()
	expected := decimal.NewFromFloat(float64(440.8975))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(430.58))})
	actual = ema1.GetAverage()
	expected = decimal.NewFromFloat(float64(439.310192))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(431.72))})
	actual = ema1.GetAverage()
	expected = decimal.NewFromFloat(float64(438.142470))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(437.87))})
	actual = ema1.GetAverage()
	expected = decimal.NewFromFloat(float64(438.100552))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(428.43))})
	actual = ema1.GetAverage()
	expected = decimal.NewFromFloat(float64(436.612775))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(428.35))})
	actual = ema1.GetAverage()
	expected = decimal.NewFromFloat(float64(435.341579))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	// EMA2 tests
	actual = ema2.GetAverage()
	expected = decimal.NewFromFloat(float64(443.289615))
	if !actual.Equals(expected) {
		t.Errorf("[MACD-EMA] Got incorrect SMA average: %f, expected: %f", actual, expected)
	}

	// Add data difference between EMA2 and EMA1 to EMA1
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(432.5))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(443.66))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(455.72))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(454.49))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(452.08))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(452.73))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(461.91))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(463.58))})
	ema1.Add(&common.Candlestick{Close: decimal.NewFromFloat(float64(461.14))})

	macd := NewMovingAverageConvergenceDivergence(ema1, ema2, 9)

	macd1 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(8.275270))
	if !macd1.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd1, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(452.08))})
	macd2 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(7.703378))
	if !macd2.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd2, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(442.66))})
	macd3 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(6.416075))
	if !macd3.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd3, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(428.91))})
	macd4 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(4.23752))
	if !macd4.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd4, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(429.79))})
	macd5 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(2.552583))
	if !macd5.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd5, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(431.99))})
	macd6 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(1.378886))
	if !macd6.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd6, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(427.72))})
	macd7 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(0.102981))
	if !macd7.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd7, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(423.2))})
	macd8 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(-1.2584))
	if !macd8.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", macd8, expected)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(426.21))})
	macd9 := macd.GetValue()
	expected = decimal.NewFromFloat(float64(-2.070558))
	if !macd9.Equals(expected) {
		t.Errorf("[MACD] Got incorrect value: %f, expected: %f", macd9, expected)
	}

	actual = macd.GetSignalLine()
	expected = decimal.NewFromFloat(float64(3.037526))
	actualHistogram := macd.GetHistogram()
	expectedHistogram := decimal.NewFromFloat(float64(-5.108084))
	if !actual.Equals(expected) {
		t.Errorf("[MACD] Got incorrect signal line: %f, expected: %f", actual, expected)
	}
	if actualHistogram != expectedHistogram {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(426.98))})
	actual = macd.GetSignalLine()
	expected = decimal.NewFromFloat(float64(1.905652))
	actualHistogram = macd.GetHistogram()
	expectedHistogram = decimal.NewFromFloat(float64(-4.527495))
	if !actual.Equals(expected) {
		t.Errorf("[MACD] Got incorrect average: %f, expected: %f", actual, expected)
	}
	if !actualHistogram.Equals(expectedHistogram) {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

	macd.OnPeriodChange(&common.Candlestick{Close: decimal.NewFromFloat(float64(435.69))})
	actual = macd.GetSignalLine()
	expected = decimal.NewFromFloat(float64(1.058708))
	actualHistogram = macd.GetHistogram()
	expectedHistogram = decimal.NewFromFloat(float64(-3.387775))
	if actual != expected {
		t.Errorf("[MACD] Got incorrect signal line: %f, expected: %f", actual, expected)
	}
	if actualHistogram != expectedHistogram {
		t.Errorf("[MACD] Got incorrect histogram: %f, expected: %f", actualHistogram, expectedHistogram)
	}

}
