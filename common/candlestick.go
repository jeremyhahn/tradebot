package common

import (
	"fmt"
	"time"
)

type Candlestick struct {
	Exchange     string        `json:"exchange"`
	CurrencyPair *CurrencyPair `json:"currency_pair"`
	Period       int           `json:"period"`
	Date         time.Time     `json:"date"`
	Open         float64       `json:"open"`
	Close        float64       `json:"close"`
	High         float64       `json:"high"`
	Low          float64       `json:"low"`
	Volume       float64       `json:"volume"`
}

func CreateCandlestick(exchangeName string, currencyPair *CurrencyPair, period int, prices []float64) *Candlestick {
	var candle = &Candlestick{
		Exchange:     exchangeName,
		CurrencyPair: currencyPair,
		Period:       period,
		Date:         time.Now(),
		Open:         prices[0],
		Close:        prices[len(prices)-1],
		Volume:       float64(len(prices))}
	for _, price := range prices {
		if price > candle.High {
			candle.High = price
		}
		if price < candle.Low {
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
	return fmt.Sprintf("Exchange: %s, CurrencyPair: %s-%s, Period: %d, Date: %s, Open: %.2f, Close: %.2f, High: %.2f, Low: %.2f, Volume: %.2f",
		candle.Exchange, base, quote, candle.Period, candle.Date, candle.Open,
		candle.Close, candle.High, candle.Low, candle.Volume)
}
