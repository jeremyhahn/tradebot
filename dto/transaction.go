package dto

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type TransactionDTO struct {
	Id                   string               `json:"id"`
	Date                 time.Time            `json:"date"`
	CurrencyPair         *common.CurrencyPair `json:"currency_pair"`
	Type                 string               `json:"type"`
	Network              string               `json:"network"`
	NetworkDisplayName   string               `json:"network_display_name"`
	Quantity             string               `json:"quantity"`
	QuantityCurrency     string               `json:"quantity_currency"`
	FiatQuantity         string               `json:"fiat_quantity"`
	FiatQuantityCurrency string               `json:"fiat_quantity_currency"`
	Price                string               `json:"price"`
	PriceCurrency        string               `json:"price_currency"`
	FiatPrice            string               `json:"fiat_price"`
	FiatPriceCurrency    string               `json:"fiat_price_currency"`
	Fee                  string               `json:"fee"`
	FeeCurrency          string               `json:"fee_currency"`
	FiatFee              string               `json:"fiat_fee"`
	FiatFeeCurrency      string               `json:"fiat_fee_currency"`
	Total                string               `json:"total"`
	TotalCurrency        string               `json:"total_currency"`
	FiatTotal            string               `json:"fiat_total"`
	FiatTotalCurrency    string               `json:"fiat_total_currency"`
	common.Transaction   `json:"-"`
}

func NewTransactionDTO() common.Transaction {
	return &TransactionDTO{}
}

func (t *TransactionDTO) GetId() string {
	return t.Id
}

func (t *TransactionDTO) GetDate() time.Time {
	return t.Date
}

func (t *TransactionDTO) GetCurrencyPair() *common.CurrencyPair {
	return t.CurrencyPair
}

func (t *TransactionDTO) GetType() string {
	return t.Type
}

func (t *TransactionDTO) GetNetwork() string {
	return t.Network
}

func (t *TransactionDTO) GetNetworkDisplayName() string {
	return t.NetworkDisplayName
}

func (t *TransactionDTO) GetQuantity() string {
	return t.Quantity
}

func (t *TransactionDTO) GetQuantityCurrency() string {
	return t.QuantityCurrency
}

func (t *TransactionDTO) GetFiatQuantity() string {
	return t.FiatQuantity
}

func (t *TransactionDTO) GetFiatQuantityCurrency() string {
	return t.FiatQuantityCurrency
}

func (t *TransactionDTO) GetPrice() string {
	return t.Price
}

func (t *TransactionDTO) GetPriceCurrency() string {
	return t.PriceCurrency
}

func (t *TransactionDTO) GetFiatPrice() string {
	return t.FiatPrice
}

func (t *TransactionDTO) GetFiatPriceCurrency() string {
	return t.FiatPriceCurrency
}

func (t *TransactionDTO) GetFee() string {
	return t.Fee
}

func (t *TransactionDTO) GetFeeCurrency() string {
	return t.FeeCurrency
}

func (t *TransactionDTO) GetFiatFee() string {
	return t.FiatFee
}

func (t *TransactionDTO) GetFiatFeeCurrency() string {
	return t.FiatFeeCurrency
}

func (t *TransactionDTO) GetTotal() string {
	return t.Total
}

func (t *TransactionDTO) GetTotalCurrency() string {
	return t.TotalCurrency
}

func (t *TransactionDTO) GetFiatTotal() string {
	return t.FiatTotal
}

func (t *TransactionDTO) GetFiatTotalCurrency() string {
	return t.FiatTotalCurrency
}

func (t *TransactionDTO) String() string {
	return fmt.Sprintf("[TransactionDTO] Id: %s, Date: %s, CurrencyPair: %s, Type: %s, Network: %s, Quantity: %s, QuantityCurrency: %s, FiatQuantity: %s,  FiatQuantityCurrency: %s, Price: %s, PriceCurrency: %s, FiatPrice: %s, FiatPriceCurrency: %s, Fee: %s, FeeCurrency: %s, FiatFee: %s, FiatFeeCurrency: %s, Total: %s, TotalCurrency: %s, FiatTotal: %s, FiatTotalCurrency: %s",
		t.Id, t.Date, t.CurrencyPair, t.Type, t.Network, t.Quantity, t.QuantityCurrency, t.FiatQuantity, t.FiatQuantityCurrency, t.Price, t.PriceCurrency, t.FiatPrice, t.FiatPriceCurrency, t.Fee, t.FeeCurrency, t.FiatFee, t.FiatFeeCurrency, t.Total, t.TotalCurrency, t.FiatTotal, t.FiatTotalCurrency)
}
