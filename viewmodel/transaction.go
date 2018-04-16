package viewmodel

import "github.com/jeremyhahn/tradebot/common"

type Transaction struct {
	Id                     string               `json:"id"`
	Network                string               `json:"network"`
	NetworkDisplayName     string               `json:"network_display_name"`
	Date                   string               `json:"date"`
	Type                   string               `json:"type"`
	Category               string               `json:"category"`
	CurrencyPair           *common.CurrencyPair `json:"currency_pair"`
	Quantity               string               `json:"quantity"`
	QuantityCurrency       string               `json:"quantity_currency"`
	FiatQuantity           string               `json:"fiat_quantity"`
	FiatQuantityCurrency   string               `json:"fiat_quantity_currency"`
	Price                  string               `json:"price"`
	PriceCurrency          string               `json:"price_currency"`
	FiatPrice              string               `json:"fiat_price"`
	FiatPriceCurrency      string               `json:"fiat_price_currency"`
	QuoteFiatPrice         string               `json:"quote_fiat_price"`
	QuoteFiatPriceCurrency string               `json:"quote_fiat_price_currency"`
	Fee                    string               `json:"fee"`
	FeeCurrency            string               `json:"fee_currency"`
	FiatFee                string               `json:"fiat_fee"`
	FiatFeeCurrency        string               `json:"fiat_fee_currency"`
	Total                  string               `json:"total"`
	TotalCurrency          string               `json:"total_currency"`
	FiatTotal              string               `json:"fiat_total"`
	FiatTotalCurrency      string               `json:"fiat_total_currency"`
}
