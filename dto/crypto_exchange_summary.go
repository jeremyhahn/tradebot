package dto

import "github.com/jeremyhahn/tradebot/common"

type CryptoExchangeSummaryDTO struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Total    float64       `json:"total"`
	Satoshis float64       `json:"satoshis"`
	Coins    []common.Coin `json:"coins"`
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

func (ce *CryptoExchangeSummaryDTO) GetTotal() float64 {
	return ce.Total
}

func (ce *CryptoExchangeSummaryDTO) GetSatoshis() float64 {
	return ce.Satoshis
}

func (ce *CryptoExchangeSummaryDTO) GetCoins() []common.Coin {
	return ce.Coins
}
