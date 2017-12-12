package indicators

import (
	"fmt"
	"math"

	"github.com/jeremyhahn/tradebot/common"
)

type BollingerBand interface {
	GetUpper() float64
	GetMiddle() float64
	GetLower() float64
	common.Indicator
}

type Bollinger struct {
	k     int
	price float64
	sma   SimpleMovingAverage
}

func NewBollingerBand(sma SimpleMovingAverage) *Bollinger {
	return &Bollinger{
		k:   2,
		sma: sma}
}

func CreateBollingerBand(sma SimpleMovingAverage, k int) *Bollinger {
	return &Bollinger{
		k:   k,
		sma: sma}
}

func (b *Bollinger) GetUpper() float64 {
	fmt.Println("GetUpper() price: ", b.price)
	for _, c := range b.sma.GetCandlesticks() {
		fmt.Printf("%+v\n", c)
	}
	return b.sma.GetAverage() + b.standardDeviation(b.getPrices(), b.price)
}

func (b *Bollinger) GetMiddle() float64 {
	return b.sma.GetAverage()
}

func (b *Bollinger) GetLower() float64 {
	return b.sma.GetAverage() - b.standardDeviation(b.getPrices(), b.price)
}

func (b *Bollinger) Calculate(price float64) {
	b.price = price
}

func (b *Bollinger) getPrices() []float64 {
	var prices []float64
	for _, c := range b.sma.GetCandlesticks() {
		prices = append(prices, c.Close)
	}
	return prices
}

func (b *Bollinger) standardDeviation(prices []float64, mean float64) float64 {
	total := 0.0
	for _, price := range prices {
		total += math.Pow(price-mean, 2)
	}
	variance := total / float64(len(prices)-1)
	return math.Sqrt(variance)
}
