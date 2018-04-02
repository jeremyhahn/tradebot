package common

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Candlestick struct {
	Exchange     string          `json:"exchange"`
	CurrencyPair *CurrencyPair   `json:"currency_pair"`
	Period       int             `json:"period"`
	Date         time.Time       `json:"date"`
	Open         decimal.Decimal `json:"open"`
	Close        decimal.Decimal `json:"close"`
	High         decimal.Decimal `json:"high"`
	Low          decimal.Decimal `json:"low"`
	Volume       decimal.Decimal `json:"volume"`
}

func CreateCandlestick(exchangeName string, currencyPair *CurrencyPair, period int, prices []decimal.Decimal) *Candlestick {
	var candle = &Candlestick{
		Exchange:     exchangeName,
		CurrencyPair: currencyPair,
		Period:       period,
		Date:         time.Now(),
		Open:         prices[0],
		Close:        prices[len(prices)-1],
		Volume:       decimal.NewFromFloat(float64(len(prices)))}
	for _, price := range prices {
		if candle.High.GreaterThan(price) {
			candle.High = price
		}
		if candle.Low.LessThan(price) {
			candle.Low = price
		}
	}
	return candle
}

func NewCandlestickPeriod(period int) time.Time {
	t := time.Now()
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	needleChunk := int(time.Since(startOfDay).Seconds())
	chunks := 86400 / period
	for i := 0; i < chunks; i++ {
		currentChunkSeconds := period * i
		lastChunkSeconds := period * (i - 1)
		if currentChunkSeconds > needleChunk {
			return startOfDay.Add(time.Duration(lastChunkSeconds) * time.Second)
		}
	}
	return t
}

func (candle *Candlestick) String() string {
	var base, quote string
	if candle.CurrencyPair != nil {
		base = candle.CurrencyPair.Base
		quote = candle.CurrencyPair.Quote
	}
	return fmt.Sprintf("Exchange: %s, CurrencyPair: %s-%s, Period: %d, Date: %s, Open: %s, Close: %s, High: %s, Low: %s, Volume: %s",
		candle.Exchange, base, quote, candle.Period, candle.Date, candle.Open,
		candle.Close, candle.High, candle.Low, candle.Volume)
}
