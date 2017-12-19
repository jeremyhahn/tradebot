package common

import (
	"time"

	logging "github.com/op/go-logging"
	"github.com/shopspring/decimal"
)

type Candlestick struct {
	Period int
	Date   time.Time
	Open   decimal.Decimal
	Close  decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Volume decimal.Decimal
}

func CreateCandlestick(logger *logging.Logger, period int, prices []decimal.Decimal) *Candlestick {
	volume, err := decimal.NewFromString(string(len(prices)))
	if err != nil {
		logger.Error(err)
	}
	var candle = &Candlestick{
		Period: period,
		Date:   time.Now(),
		Open:   prices[0],
		Close:  prices[len(prices)-1],
		Volume: volume}
	for _, price := range prices {
		if price.GreaterThan(candle.High) {
			candle.High = price
		}
		if price.LessThan(candle.Low) {
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
