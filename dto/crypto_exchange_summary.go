package dto

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type CryptoExchangeSummaryDTO struct {
	Name     string          `json:"name"`
	URL      string          `json:"url"`
	Total    decimal.Decimal `json:"total"`
	Satoshis decimal.Decimal `json:"satoshis"`
	Coins    []common.Coin   `json:"coins"`
	common.CryptoExchangeSummary
}

func NewCryptExchange() common.CryptoExchangeSummary {
	return &CryptoExchangeSummaryDTO{}
}

func (ce *CryptoExchangeSummaryDTO) GetName() string {
	return ce.Name
}

func (ce *CryptoExchangeSummaryDTO) GetURL() string {
	return ce.URL
}

func (ce *CryptoExchangeSummaryDTO) GetTotal() decimal.Decimal {
	return ce.Total
}

func (ce *CryptoExchangeSummaryDTO) GetSatoshis() decimal.Decimal {
	return ce.Satoshis
}

func (ce *CryptoExchangeSummaryDTO) GetCoins() []common.Coin {
	return ce.Coins
}
