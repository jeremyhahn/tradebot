package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type BBandParams struct {
	Period int64
	K      float64
}

type BollingerBandsImpl struct {
	name        string
	displayName string
	price       decimal.Decimal
	sma         indicators.SimpleMovingAverage
	params      *BBandParams
	indicators.BollingerBands
}

func NewBollingerBands(candles []common.Candlestick) (common.FinancialIndicator, error) {
	params := []string{"20", "2"}
	return CreateBollingerBands(candles, params)
}

func CreateBollingerBands(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &BollingerBandsImpl{}
		params = temp.GetDefaultParameters()
	}
	period, _ := strconv.ParseInt(params[0], 10, 64)
	k, _ := strconv.ParseFloat(params[1], 64)
	smaIndicator, err := NewSimpleMovingAverage(candles[:period])
	sma := smaIndicator.(indicators.SimpleMovingAverage)
	if err != nil {
		return nil, err
	}
	bollinger := &BollingerBandsImpl{
		name:        "BollingerBands",
		displayName: "Bollinger Bands®",
		sma:         sma,
		params: &BBandParams{
			Period: period,
			K:      k}}
	for _, c := range candles[period:] {
		bollinger.OnPeriodChange(&c)
	}
	return bollinger, nil
}

func (b *BollingerBandsImpl) GetUpper() decimal.Decimal {
	//return util.RoundFloat(b.sma.GetAverage()+(b.StandardDeviation()*2), 2)
	return b.sma.GetAverage().Add(b.StandardDeviation().Mul(decimal.NewFromFloat(2)))
}

func (b *BollingerBandsImpl) GetMiddle() decimal.Decimal {
	//return util.RoundFloat(b.sma.GetAverage(), 2)
	return b.sma.GetAverage()
}

func (b *BollingerBandsImpl) GetLower() decimal.Decimal {
	//return util.RoundFloat(b.sma.GetAverage()-(b.StandardDeviation()*2), 2)
	return b.sma.GetAverage()
}

func (b *BollingerBandsImpl) StandardDeviation() decimal.Decimal {
	return b.calculateStandardDeviation(b.sma.GetPrices(), b.sma.GetAverage())
}

func (b *BollingerBandsImpl) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	total := decimal.NewFromFloat(0)
	prices := b.sma.GetPrices()
	prices[0] = price
	for _, p := range prices {
		total = total.Add(p)
	}
	avg := total.Div(decimal.NewFromFloat(float64(len(prices))))
	stdDev := b.calculateStandardDeviation(prices, avg)
	//upper := util.RoundFloat(avg+(stdDev*b.params.K), 2)
	//middle := util.RoundFloat(avg, 2)
	///lower := util.RoundFloat(avg-(stdDev*b.params.K), 2)
	upper := avg.Add(stdDev.Mul(decimal.NewFromFloat(b.params.K)))
	middle := avg
	lower := avg.Sub(stdDev.Mul(decimal.NewFromFloat(b.params.K)))
	return upper, middle, lower
}

func (b *BollingerBandsImpl) calculateStandardDeviation(prices []decimal.Decimal, mean decimal.Decimal) decimal.Decimal {
	total := decimal.NewFromFloat(0)
	for _, price := range prices {
		//total += math.Pow(price-mean, 2)
		total = total.Add(price.Sub(mean).Pow(decimal.NewFromFloat(2)))
	}
	variance := total.Div(decimal.NewFromFloat(float64(len(prices))))
	//return util.RoundFloat(math.Sqrt(variance), 2)
	return variance.Mul(variance)
}

func (b *BollingerBandsImpl) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[BollingerBands] OnPeriodChange: %s", candle.ToString())
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
