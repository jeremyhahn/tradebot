package viewmodel

import "github.com/jeremyhahn/tradebot/common"

type Order struct {
	Id                 string               `json:"id"`
	Exchange           string               `json:"exchange"`
	Date               string               `json:"date"`
	Type               string               `json:"type"`
	CurrencyPair       *common.CurrencyPair `json:"currency_pair"`
	Quantity           float64              `json:"quantity"`
	QuantityCurrency   string               `json:"quantity_currency"`
	Price              float64              `json:"price"`
	PriceCurrency      string               `json:"price_currency"`
	Fee                float64              `json:"fee"`
	FeeCurrency        string               `json:"fee_currency"`
	Total              float64              `json:"total"`
	TotalCurrency      string               `json:"total_currency"`
	HistoricalPrice    float64              `json:"historical_price"`
	HistoricalCurrency string               `json:"historical_currency"`
}
