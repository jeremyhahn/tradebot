package common

import "time"

type Candlestick struct {
	Period int
	Date   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}

func CreateCandlestick(period int, prices []float64) *Candlestick {
	var candle = &Candlestick{
		Period: period,
		Date:   time.Now(),
		Open:   prices[0],
		Close:  prices[len(prices)-1],
		Volume: float64(len(prices))}
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
