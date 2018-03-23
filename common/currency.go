package common

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Currency struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Symbol       string          `json:"symbol"`
	BaseUnit     int32           `json:"base_unit"`
	DecimalPlace int32           `json:"decimal_place"`
	TxFee        decimal.Decimal `json:"tx_fee"`
}

func (c *Currency) GetID() string {
	return c.ID
}

func (c *Currency) GetName() string {
	return c.Name
}

func (c *Currency) GetSymbol() string {
	return c.Symbol
}

func (c *Currency) GetBaseUnit() int32 {
	return c.BaseUnit
}

func (c *Currency) GetDecimalPlace() int32 {
	return c.DecimalPlace
}

func (c *Currency) GetTransactionFee() decimal.Decimal {
	return c.TxFee
}

func (c *Currency) IsFiat() bool {
	_, found := FiatCurrencies[c.ID]
	return found
}

func (c *Currency) IsCrypto() bool {
	_, found := FiatCurrencies[c.ID]
	return !found
}

func (c *Currency) String() string {
	return fmt.Sprintf("[Currency] Id: %s, Symbol: %s, Name: %s, BaseUnit: %d, DecimalPlace: %d, TxFee: %s",
		c.ID, c.Symbol, c.Name, c.BaseUnit, c.DecimalPlace, c.TxFee.String())
}
