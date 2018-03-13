package entity

import "time"

type Order struct {
	Id                 uint      `gorm:"primary_key"`
	UserId             uint      `gorm:"unique_index:idx_orderhistory"`
	Date               time.Time `gorm:"unique_index:idx_orderhistory"`
	Exchange           string    `gorm:"unique_index:idx_orderhistory"`
	Type               string
	Currency           string `gorm:"unique_index:idx_orderhistory"`
	Quantity           float64
	QuantityCurrency   string `gorm:"varchar(6)"`
	Price              float64
	PriceCurrency      string `gorm:"varchar(6)"`
	Fee                float64
	FeeCurrency        string  `gorm:"varchar(6)"`
	Total              float64 `gorm:"unique_index:idx_orderhistory"`
	TotalCurrency      string  `gorm:"varchar(6)"`
	HistoricalPrice    float64
	HistoricalCurrency string `gorm:"varchar(6)"`
	OrderEntity
}

func (order *Order) GetId() uint {
	return order.Id
}

func (order *Order) GetUserId() uint {
	return order.UserId
}

func (order *Order) GetDate() time.Time {
	return order.Date
}

func (order *Order) GetExchange() string {
	return order.Exchange
}

func (order *Order) GetType() string {
	return order.Type
}

func (order *Order) GetCurrency() string {
	return order.Currency
}

func (order *Order) GetQuantity() float64 {
	return order.Quantity
}

func (order *Order) GetPrice() float64 {
	return order.Price
}

func (order *Order) GetPriceCurrency() string {
	return order.PriceCurrency
}

func (order *Order) GetFee() float64 {
	return order.Fee
}

func (order *Order) GetFeeCurrency() string {
	return order.FeeCurrency
}

func (order *Order) GetTotal() float64 {
	return order.Total
}

func (order *Order) GetTotalCurrency() string {
	return order.TotalCurrency
}

func (order *Order) GetQuantityCurrency() string {
	return order.QuantityCurrency
}

func (order *Order) GetHistoricalPrice() float64 {
	return order.HistoricalPrice
}

func (order *Order) GetHistoricalCurrency() string {
	return order.HistoricalCurrency
}
