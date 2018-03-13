package dto

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type TransactionDTO struct {
	Date               time.Time            `json:"date"`
	CurrencyPair       *common.CurrencyPair `json:"currency_pair"`
	Type               string               `json:"type"`
	Source             string               `json:"source"`
	Amount             float64              `json:"amount"`
	AmountCurrency     string               `json:"amount_currency"`
	Fee                float64              `json:"fee"`
	FeeCurrency        string               `json:"fee_currency"`
	Total              float64              `json:"total"`
	TotalCurrency      string               `json:"total_currency"`
	HistoricalPrice    float64              `json:"historical_price"`
	HistoricalCurrency string               `json:"historical_currency"`
	common.Transaction
}

func NewTransactionDTO() common.Transaction {
	return &TransactionDTO{}
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

func (t *TransactionDTO) GetSource() string {
	return t.Source
}

func (t *TransactionDTO) GetAmount() float64 {
	return t.Amount
}

func (t *TransactionDTO) GetAmountCurrency() string {
	return t.AmountCurrency
}

func (t *TransactionDTO) GetFee() float64 {
	return t.Fee
}

func (t *TransactionDTO) GetFeeCurrency() string {
	return t.FeeCurrency
}

func (t *TransactionDTO) GetTotal() float64 {
	return t.Total
}

func (t *TransactionDTO) GetTotalCurrency() string {
	return t.TotalCurrency
}

func (t *TransactionDTO) GetHistoricalPrice() float64 {
	return t.HistoricalPrice
}

func (t *TransactionDTO) GetHistoricalCurrency() string {
	return t.HistoricalCurrency
}

func (t *TransactionDTO) String() string {
	return fmt.Sprintf("[TransactionDTO] Date: %s, Amount: %f, Currency: %s, Type: %s, Source: %s, Amount: %.8f, Fee: %.8f, Total: %.8f, historicalPrice: %f",
		t.Date, t.Amount, t.CurrencyPair, t.Type, t.Source, t.Amount, t.Fee, t.Total, t.HistoricalPrice)
}
