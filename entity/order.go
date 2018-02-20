package entity

import "time"

type Order struct {
	Id       uint      `gorm:"primary_key"`
	UserId   uint      `gorm:"unique_index:idx_orderhistory"`
	Date     time.Time `gorm:"unique_index:idx_orderhistory"`
	Exchange string    `gorm:"unique_index:idx_orderhistory"`
	Type     string
	Currency string `gorm:"unique_index:idx_orderhistory"`
	Quantity float64
	Price    float64
	Fee      float64
	Total    float64 `gorm:"unique_index:idx_orderhistory"`
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

func (order *Order) GetFee() float64 {
	return order.Fee
}

func (order *Order) GetTotal() float64 {
	return order.Total
}
