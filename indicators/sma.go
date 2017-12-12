package indicators

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type SimpleMovingAverage interface {
	Add(candle *common.Candlestick) float64
	GetCandlesticks() []common.Candlestick
	GetSize() int
	GetAverage() float64
	common.PriceListener
}

type SMA struct {
	size         int
	candlesticks []common.Candlestick
	prices       []float64
	count        int
	index        int
	average      float64
	SimpleMovingAverage
}

func NewSimpleMovingAverage(candles []common.Candlestick) *SMA {
	return CreateSimpleMovingAverage(candles, len(candles))
}

func CreateSimpleMovingAverage(candles []common.Candlestick, size int) *SMA {
	var prices []float64
	var total float64
	for _, c := range candles {
		fmt.Println("Adding candlestick: ", c.Date, c.Close)
		prices = append(prices, c.Close)
		total = total + c.Close
	}
	return &SMA{
		size:         size,
		candlesticks: candles,
		prices:       prices,
		count:        len(candles),
		index:        len(candles) - 1,
		average:      total / float64(len(candles))}
}

func (sma *SMA) Add(candle *common.Candlestick) float64 {
	if sma.count == sma.size {
		sma.index = (sma.index + 1) % sma.size
		sma.average += (candle.Close - sma.prices[sma.index]) / float64(sma.count)
		sma.prices[sma.index] = candle.Close
	} else if sma.count != 0 {
		sma.index = (sma.index + 1) % sma.size
		sma.average = (candle.Close + float64(sma.count)*sma.average) / float64(sma.count+1)
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

func (sma *SMA) GetAverage() float64 {
	return sma.average
}

func (sma *SMA) GetSize() int {
	return sma.size
}

func (sma *SMA) OnCandlestickCreated(candle *common.Candlestick) {
	fmt.Println("SMA OnCandlestickCreated: ", candle.Date, candle.Close)
	sma.Add(candle)
}

func (sma *SMA) OnPriceChange(price float64) {
	//fmt.Println("SMA OnPriceChange: ", price)
}
