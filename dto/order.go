package dto

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type OrderDTO struct {
	Id            string               `json:"id"`
	Exchange      string               `json:"exchange"`
	Date          time.Time            `json:"date"`
	Type          string               `json:"type"`
	CurrencyPair  *common.CurrencyPair `json:"currency_pair"`
	Quantity      float64              `json:"quantity"`
	Price         float64              `json:"price"`
	Fee           float64              `json:"fee"`
	Total         float64              `json:"total"`
	PriceCurrency string               `json:"price_currency"`
	FeeCurrency   string               `json:"fee_currency"`
	TotalCurrency string               `json:"total_currency"`
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

func (dto *OrderDTO) GetCurrencyPair() *common.CurrencyPair {
	return dto.CurrencyPair
}

func (dto *OrderDTO) GetQuantity() float64 {
	return dto.Quantity
}

func (dto *OrderDTO) GetPrice() float64 {
	return dto.Price
}

func (dto *OrderDTO) GetFee() float64 {
	return dto.Fee
}

func (dto *OrderDTO) GetTotal() float64 {
	return dto.Total
}

func (dto *OrderDTO) GetPriceCurrency() string {
	return dto.PriceCurrency
}

func (dto *OrderDTO) GetFeeCurrency() string {
	return dto.FeeCurrency
}

func (dto *OrderDTO) GetTotalCurrency() string {
	return dto.TotalCurrency
}

func (dto *OrderDTO) String() string {
	return fmt.Sprintf("[OrderDTO] Id: %s, Exchange: %s, Date: %s, Type: %s, CurrencyPair: %s, Quantity: %f, Price: %f, Fee: %f, Total: %f, PriceCurrency: %s, FeeCurrency: %s, TotalCurrency: %s",
		dto.Id, dto.Exchange, dto.Date, dto.Type, dto.GetCurrencyPair().String(), dto.Quantity, dto.Price, dto.Fee, dto.Total, dto.PriceCurrency, dto.FeeCurrency, dto.TotalCurrency)
}
