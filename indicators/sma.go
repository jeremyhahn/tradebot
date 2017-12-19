package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type SMA struct {
	size         int
	candlesticks []common.Candlestick
	prices       []decimal.Decimal
	count        int
	index        int
	average      decimal.Decimal
	sum          decimal.Decimal
	common.MovingAverage
	common.PeriodListener
}

func NewSimpleMovingAverage(candles []common.Candlestick) *SMA {
	return CreateSimpleMovingAverage(candles, len(candles))
}

func CreateSimpleMovingAverage(candles []common.Candlestick, size int) *SMA {
	var prices []decimal.Decimal
	var total decimal.Decimal
	for _, c := range candles {
		prices = append(prices, c.Close)
		total = total.Add(c.Close)
	}
	zero := decimal.NewFromFloat(float64(0))
	sma := &SMA{
		size:         size,
		candlesticks: candles,
		prices:       make([]decimal.Decimal, size),
		count:        0,
		index:        0,
		average:      zero,
		sum:          zero}
	candleLen := len(candles)
	if candleLen > 0 {
		sma.prices = prices
		sma.count = candleLen
		sma.index = candleLen - 1
		sma.average = total.Div(decimal.NewFromFloat(float64(size)))
	}
	return sma
}

func (sma *SMA) Add(candle *common.Candlestick) decimal.Decimal {
	if sma.count == sma.size {
		sma.index = (sma.index + 1) % sma.size
		sma.average = sma.average.Add((candle.Close.Sub(sma.prices[sma.index])).Div(decimal.NewFromFloat(float64(sma.count))))
		sma.prices[sma.index] = candle.Close
	} else if sma.count != 0 && sma.count < sma.size {
		sma.index = (sma.index + 1) % sma.size
		sma.average = sma.average.Add(candle.Close.Add(decimal.NewFromFloat(float64(sma.count))).Mul(sma.average)).Div(decimal.NewFromFloat(float64(sma.count + 1)))
		sma.prices[sma.index] = candle.Close
		sma.count++
	} else {
		sma.average = candle.Close
		sma.prices[0] = candle.Close
		sma.count = 1
	}
	return sma.average
}

func (sma *SMA) GetCandlesticks() []common.Candlestick {
	return sma.candlesticks
}

func (sma *SMA) GetAverage() decimal.Decimal {
	return sma.average
}

func (sma *SMA) GetSize() int {
	return sma.size
}

func (sma *SMA) GetGainsAndLosses() (decimal.Decimal, decimal.Decimal) {
	var gains, losses decimal.Decimal
	var lastClose = sma.candlesticks[0].Close
	for _, c := range sma.candlesticks {
		difference := c.Close.Sub(lastClose)
		if difference.LessThan(decimal.NewFromFloat(float64(0))) {
			losses = losses.Add(difference.Abs())
		} else {
			gains = gains.Add(difference)
		}
		lastClose = c.Close
	}
	return gains, losses
}

func (sma *SMA) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[SMA] OnPeriodChange: ", candle.Date, candle.Close)
	sma.Add(candle)
}
