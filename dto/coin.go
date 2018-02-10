package dto

import "github.com/jeremyhahn/tradebot/common"

type CoinDTO struct {
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	Available float64 `json:"available"`
	Pending   float64 `json:"pending"`
	Price     float64 `json:"price"`
	Address   string  `json:"address"`
	Total     float64 `json:"total"`
	BTC       float64 `json:"btc"`
	common.Coin
}

func NewCoinDTO() common.Coin {
	return &CoinDTO{}
}

func (dto *CoinDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *CoinDTO) GetBalance() float64 {
	return dto.Balance
}

func (dto *CoinDTO) GetAvailable() float64 {
	return dto.Available
}

func (dto *CoinDTO) GetPending() float64 {
	return dto.Pending
}

func (dto *CoinDTO) GetPrice() float64 {
	return dto.Price
}

func (dto *CoinDTO) GetAddress() string {
	return dto.Address
}

func (dto *CoinDTO) GetTotal() float64 {
	return dto.Total
}

func (dto *CoinDTO) GetBTC() float64 {
	return dto.BTC
}

func (dto *CoinDTO) IsBitcoin() bool {
	return dto.Currency == "BTC"
}
