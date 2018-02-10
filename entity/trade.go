package entity

import "time"

type Trade struct {
	Id        uint   `gorm:"primary_key"`
	ChartId   uint   `gorm:"foreign_key"`
	UserId    uint   `gorm:"foreign_key;index"`
	Base      string `gorm:"index"`
	Quote     string `gorm:"index"`
	Exchange  string `gorm:"index"`
	Date      time.Time
	Type      string
	Price     float64
	Amount    float64
	ChartData string
	TradeEntity
}

func (trade *Trade) GetId() uint {
	return trade.Id
}

func (trade *Trade) GetChartId() uint {
	return trade.ChartId
}

func (trade *Trade) GetUserId() uint {
	return trade.UserId
}

func (trade *Trade) GetBase() string {
	return trade.Base
}

func (trade *Trade) GetQuote() string {
	return trade.Quote
}

func (trade *Trade) GetExchangeName() string {
	return trade.Exchange
}

func (trade *Trade) GetDate() time.Time {
	return trade.Date
}

func (trade *Trade) GetType() string {
	return trade.Type
}

func (trade *Trade) GetPrice() float64 {
	return trade.Price
}

func (trade *Trade) GetAmount() float64 {
	return trade.Amount
}

func (trade *Trade) GetChartData() string {
	return trade.ChartData
}
