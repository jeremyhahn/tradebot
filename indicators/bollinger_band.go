package indicators

import (
	"fmt"
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
	return util.RoundFloat(b.sma.GetAverage()+(b.StandardDeviation()*2), 2)
}

func (b *Bollinger) GetMiddle() float64 {
	return util.RoundFloat(b.sma.GetAverage(), 2)
}

func (b *Bollinger) GetLower() float64 {
	return util.RoundFloat(b.sma.GetAverage()-(b.StandardDeviation()*2), 2)
}

func (b *Bollinger) StandardDeviation() float64 {
	return b.CalculateStandardDeviation(b.sma.GetPrices(), b.sma.GetAverage())
}

func (b *Bollinger) Calculate(price float64) (float64, float64, float64) {
	total := 0.0
	prices := b.sma.GetPrices()
	prices[0] = price
	for _, p := range prices {
		total += p
	}
	avg := total / float64(len(prices))
	stdDev := b.CalculateStandardDeviation(prices, avg)
	upper := util.RoundFloat(avg+(stdDev*2), 2)
	middle := util.RoundFloat(avg, 2)
	lower := util.RoundFloat(avg-(stdDev*2), 2)
	return upper, middle, lower
}

func (b *Bollinger) CalculateStandardDeviation(prices []float64, mean float64) float64 {
	total := 0.0
	for _, price := range prices {
		total += math.Pow(price-mean, 2)
	}
	variance := total / float64(len(prices))
	return util.RoundFloat(math.Sqrt(variance), 2)
}

func (b *Bollinger) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[Bollinger] OnPeriodChange: ", candle.Date, candle.Close)
	b.sma.Add(candle)
}
