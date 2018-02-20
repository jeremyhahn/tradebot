package viewmodel

import "github.com/jeremyhahn/tradebot/common"

type Order struct {
	Id            string               `json:"id"`
	Exchange      string               `json:"exchange"`
	Date          string               `json:"date"`
	Type          string               `json:"type"`
	CurrencyPair  *common.CurrencyPair `json:"currency_pair"`
	Quantity      float64              `json:"quantity"`
	Price         float64              `json:"price"`
	Fee           float64              `json:"fee"`
	Total         float64              `json:"total"`
	PriceCurrency string               `json:"price_currency"`
	FeeCurrency   string               `json:"fee_currency"`
	TotalCurrency string               `json:"total_currency"`
}
