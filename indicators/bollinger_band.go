package indicators

import (
	"math"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
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
	sma   common.MovingAverage
}

func NewBollingerBand(sma common.MovingAverage) *Bollinger {
	return &Bollinger{
		k:   2,
		sma: sma}
}

func CreateBollingerBand(sma common.MovingAverage, k int) *Bollinger {
	return &Bollinger{
		k:   k,
		sma: sma}
}

func (b *Bollinger) GetUpper() float64 {
	return util.RoundFloat(
		b.sma.GetAverage()+b.standardDeviation(b.getPrices(), b.price),
		2)
}

func (b *Bollinger) GetMiddle() float64 {
	return util.RoundFloat(b.sma.GetAverage(), 2)
}

func (b *Bollinger) GetLower() float64 {
	return util.RoundFloat(
		b.sma.GetAverage()-b.standardDeviation(b.getPrices(), b.price),
		2)
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
		total += math.Pow(price-mean, float64(b.k))
	}
	variance := total / float64(len(prices)-1)
	return math.Sqrt(variance)
}

func (b *Bollinger) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[Bollinger] OnPeriodChange: ", candle.Date, candle.Close)
	b.sma.Add(candle)
}
