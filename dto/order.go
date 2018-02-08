package dto

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type OrderDTO struct {
	Id       string    `json:"id"`
	Exchange string    `json:"exchange"`
	Date     time.Time `json:"date"`
	Type     string    `json:"type"`
	Currency string    `json:"currency"`
	Quantity float64   `json:"quantity"`
	Price    float64   `json:"price"`
	common.Order
}

func NewOrderDTO() common.Order {
	return &OrderDTO{}
}

func (dto *OrderDTO) GetId() string {
	return dto.Id
}

func (dto *OrderDTO) GetExchange() string {
	return dto.Exchange
}

func (dto *OrderDTO) GetDate() time.Time {
	return dto.Date
}

func (dto *OrderDTO) GetType() string {
	return dto.Type
}

func (dto *OrderDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *OrderDTO) GetQuantity() float64 {
	return dto.Quantity
}

func (dto *OrderDTO) GetPrice() float64 {
	return dto.Price
}
