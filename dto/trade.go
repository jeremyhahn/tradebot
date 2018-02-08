package dto

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type TradeDTO struct {
	Id        uint      `json:"id"`
	ChartId   uint      `json:"chart_id"`
	UserId    uint      `json:"user_id"`
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	Exchange  string    `json:"exchange"`
	Date      time.Time `json:"date"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	ChartData string    `json:"chart_data"`
	common.Trade
}

func NewTradeDTO() common.Trade {
	return &TradeDTO{}
}

func (dto *TradeDTO) GetId() uint {
	return dto.Id
}

func (dto *TradeDTO) GetChartId() uint {
	return dto.ChartId
}

func (dto *TradeDTO) GetUserId() uint {
	return dto.UserId
}

func (dto *TradeDTO) GetBase() string {
	return dto.Base
}

func (dto *TradeDTO) GetQuote() string {
	return dto.Quote
}

func (dto *TradeDTO) GetExchange() string {
	return dto.Exchange
}

func (dto *TradeDTO) GetDate() time.Time {
	return dto.Date
}

func (dto *TradeDTO) GetType() string {
	return dto.Type
}

func (dto *TradeDTO) GetPrice() float64 {
	return dto.Price
}

func (dto *TradeDTO) GetAmount() float64 {
	return dto.Amount
}

func (dto *TradeDTO) GetChartData() string {
	return dto.ChartData
}
