package entity

type PriceHistory struct {
	Time      int64   `json:"time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	MarketCap int64   `json:"marketCap"`
	PriceHistoryEntity
}

func (ph *PriceHistory) GetTime() int64 {
	return ph.Time
}

func (ph *PriceHistory) GetOpen() float64 {
	return ph.Open
}

func (ph *PriceHistory) GetHigh() float64 {
	return ph.High
}

func (ph *PriceHistory) GetLow() float64 {
	return ph.Low
}

func (ph *PriceHistory) GetClose() float64 {
	return ph.Close
}

func (ph *PriceHistory) GetVolume() float64 {
	return ph.Volume
}

func (ph *PriceHistory) GetMarketCap() int64 {
	return ph.MarketCap
}
