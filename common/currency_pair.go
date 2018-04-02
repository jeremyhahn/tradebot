package common

import (
	"errors"
	"fmt"
	"strings"
)

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

func NewCurrencyPair(currencyPair, localCurrency string) (*CurrencyPair, error) {
	pieces := strings.Split(currencyPair, "-")
	if len(pieces) < 2 {
		return nil, errors.New(fmt.Sprintf("[NewCurrencyPair] Invalid currency pair format: %s", currencyPair))
	}
	return &CurrencyPair{
		Base:          pieces[0],
		Quote:         pieces[1],
		LocalCurrency: localCurrency}, nil
}

func (cp *CurrencyPair) String() string {
	return fmt.Sprintf("%s-%s", cp.Base, cp.Quote)
}

func (cp *CurrencyPair) Equals(currencyPair *CurrencyPair) bool {
	return cp.Base == currencyPair.Base &&
		cp.Quote == currencyPair.Quote
}
