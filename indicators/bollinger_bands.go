package indicators

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

type BollingerBands interface {
	GetUpper() float64
	GetMiddle() float64
	GetLower() float64
	StandardDeviation() float64
	Calculate(price float64) (float64, float64, float64)
	common.FinancialIndicator
}

type BBandParams struct {
	Period int64
	K      float64
}

type BollingerBandsImpl struct {
	name        string
	displayName string
	price       float64
	sma         common.MovingAverage
	params      *BBandParams
	BollingerBands
}

func NewBollingerBands(candles []common.Candlestick) BollingerBands {
	params := []string{"20", "2"}
	return CreateBollingerBands(candles, params)
}

func CreateBollingerBands(candles []common.Candlestick, params []string) BollingerBands {
	period, _ := strconv.ParseInt(params[0], 10, 64)
	k, _ := strconv.ParseFloat(params[1], 64)
	sma := NewSimpleMovingAverage(candles[:period])
	return &BollingerBandsImpl{
		name:        "BollingerBands",
		displayName: "Bollinger Bands®",
		sma:         sma,
		params: &BBandParams{
			Period: period,
			K:      k}}
}

func (b *BollingerBandsImpl) GetUpper() float64 {
	return util.RoundFloat(b.sma.GetAverage()+(b.StandardDeviation()*2), 2)
}

func (b *BollingerBandsImpl) GetMiddle() float64 {
	return util.RoundFloat(b.sma.GetAverage(), 2)
}

func (b *BollingerBandsImpl) GetLower() float64 {
	return util.RoundFloat(b.sma.GetAverage()-(b.StandardDeviation()*2), 2)
}

func (b *BollingerBandsImpl) StandardDeviation() float64 {
	return b.calculateStandardDeviation(b.sma.GetPrices(), b.sma.GetAverage())
}

func (b *BollingerBandsImpl) Calculate(price float64) (float64, float64, float64) {
	total := 0.0
	prices := b.sma.GetPrices()
	prices[0] = price
	for _, p := range prices {
		total += p
	}
	avg := total / float64(len(prices))
	stdDev := b.calculateStandardDeviation(prices, avg)
	upper := util.RoundFloat(avg+(stdDev*b.params.K), 2)
	middle := util.RoundFloat(avg, 2)
	lower := util.RoundFloat(avg-(stdDev*b.params.K), 2)
	return upper, middle, lower
}

func (b *BollingerBandsImpl) calculateStandardDeviation(prices []float64, mean float64) float64 {
	total := 0.0
	for _, price := range prices {
		total += math.Pow(price-mean, 2)
	}
	variance := total / float64(len(prices))
	return util.RoundFloat(math.Sqrt(variance), 2)
}

func (b *BollingerBandsImpl) OnPeriodChange(candle *common.Candlestick) {
	fmt.Println("[Bollinger] OnPeriodChange: ", candle.Date, candle.Close)
	b.sma.Add(candle)
}

func (b *BollingerBandsImpl) GetName() string {
	return b.name
}

func (b *BollingerBandsImpl) GetParameters() []string {
	return []string{
		fmt.Sprintf("%d", b.params.Period),
		fmt.Sprintf("%f", b.params.K)}
}

func (b *BollingerBandsImpl) GetDefaultParameters() []string {
	return []string{"20", "2"}
}

func (b *BollingerBandsImpl) GetDisplayName() string {
	return b.displayName
}
