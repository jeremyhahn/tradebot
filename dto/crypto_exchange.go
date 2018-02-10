package dto

import "github.com/jeremyhahn/tradebot/common"

type CryptoExchangeDTO struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Total    float64       `json:"total"`
	Satoshis float64       `json:"satoshis"`
	Coins    []common.Coin `json:"coins"`
}

func NewCryptExchange() common.CryptoExchange {
	return &CryptoExchangeDTO{}
}

func (ce *CryptoExchangeDTO) GetName() string {
	return ce.Name
}

func (ce *CryptoExchangeDTO) GetURL() string {
	return ce.URL
}

func (ce *CryptoExchangeDTO) GetTotal() float64 {
	return ce.Total
}

func (ce *CryptoExchangeDTO) GetSatoshis() float64 {
	return ce.Satoshis
}

func (ce *CryptoExchangeDTO) GetCoins() []common.Coin {
	return ce.Coins
}
