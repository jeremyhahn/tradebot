package dto

import (
	"github.com/jeremyhahn/tradebot/common"
)

type PriceHistoryDTO struct {
	Time      int64   `json:"time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	MarketCap int64   `json:"marketCap"`
	common.PriceHistory
}

func NewPriceHistoryDTO() common.PriceHistory {
	return &PriceHistoryDTO{}
}

func (ph *PriceHistoryDTO) GetTime() int64 {
	return ph.Time
}

func (ph *PriceHistoryDTO) GetOpen() float64 {
	return ph.Open
}

func (ph *PriceHistoryDTO) GetHigh() float64 {
	return ph.High
}

func (ph *PriceHistoryDTO) GetLow() float64 {
	return ph.Low
}

func (ph *PriceHistoryDTO) GetClose() float64 {
	return ph.Close
}

func (ph *PriceHistoryDTO) GetVolume() float64 {
	return ph.Volume
}

func (ph *PriceHistoryDTO) GetMarketCap() int64 {
	return ph.MarketCap
}
