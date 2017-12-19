package indicators

import (
	"math"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
	"github.com/shopspring/decimal"
)

type BollingerBand interface {
	GetUpper() decimal.Decimal
	GetMiddle() decimal.Decimal
	GetLower() decimal.Decimal
	common.Indicator
}

type Bollinger struct {
	logger *logging.Logger
	k      int
	price  decimal.Decimal
	sma    common.MovingAverage
}

func NewBollingerBand(logger *logging.Logger, sma common.MovingAverage) *Bollinger {
	return &Bollinger{
		logger: logger,
		k:      2,
		sma:    sma}
}

func CreateBollingerBand(logger *logging.Logger, sma common.MovingAverage, k int) *Bollinger {
	return &Bollinger{
		k:   k,
		sma: sma}
}

func (b *Bollinger) GetUpper() decimal.Decimal {
	return b.sma.GetAverage().Add(b.standardDeviation(b.getPrices(), b.price))
}

func (b *Bollinger) GetMiddle() decimal.Decimal {
	return b.sma.GetAverage()
}

func (b *Bollinger) GetLower() decimal.Decimal {
	return b.sma.GetAverage().Sub(b.standardDeviation(b.getPrices(), b.price))
}

func (b *Bollinger) Calculate(price decimal.Decimal) {
	b.price = price
}

func (b *Bollinger) getPrices() []decimal.Decimal {
	var prices []decimal.Decimal
	for _, c := range b.sma.GetCandlesticks() {
		prices = append(prices, c.Close)
	}
	return prices
}

func (b *Bollinger) standardDeviation(prices []decimal.Decimal, mean decimal.Decimal) decimal.Decimal {
	total := decimal.NewFromFloat(0.0)
	for _, price := range prices {
		total = total.Add(price.Pow(decimal.NewFromFloat(float64(b.k))))
	}
	variance := total.Div(decimal.NewFromFloat(float64(len(prices) - 1)))
	f, exact := variance.Float64()
	if !exact {
		b.logger.Error("Bollinger float conversion failure: ", exact)
	}
	return decimal.NewFromFloat(math.Sqrt(f))
}

func (b *Bollinger) OnPeriodChange(candle *common.Candlestick) {
	//fmt.Println("[Bollinger] OnPeriodChange: ", candle.Date, candle.Close)
	b.sma.Add(candle)
}
