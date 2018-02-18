package common

import "fmt"

type CurrencyPair struct {
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	LocalCurrency string `json:"local_currency"`
}

func (cp *CurrencyPair) String() string {
	return fmt.Sprintf("%s-%s", cp.Base, cp.Quote)
}
