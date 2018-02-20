package common

import (
	"fmt"
	"strings"
)

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

func NewCurrencyPair(currencyPair, localCurrency string) *CurrencyPair {
	pieces := strings.Split(currencyPair, "-")
	return &CurrencyPair{
		Base:          pieces[0],
		Quote:         pieces[1],
		LocalCurrency: localCurrency}
}

func (cp *CurrencyPair) String() string {
	return fmt.Sprintf("%s-%s", cp.Base, cp.Quote)
}
