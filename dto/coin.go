package dto

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type CoinDTO struct {
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Exchange    string  `json:"exchange"`
	Balance     float64 `json:"balance"`
	Available   float64 `json:"available"`
	Pending     float64 `json:"pending"`
	Address     string  `json:"address"`
	Total       float64 `json:"total"`
	BTC         float64 `json:"btc"`
	USD         float64 `json:"usd"`
	common.Coin `json:"-"`
}

func NewCoinDTO() common.Coin {
	return &CoinDTO{}
}

func (dto *CoinDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *CoinDTO) GetPrice() float64 {
	return dto.Price
}

func (dto *CoinDTO) GetExchange() string {
	return dto.Exchange
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

func (dto *CoinDTO) GetAddress() string {
	return dto.Address
}

func (dto *CoinDTO) GetTotal() float64 {
	return dto.Total
}

func (dto *CoinDTO) GetBTC() float64 {
	return dto.BTC
}

func (dto *CoinDTO) GetUSD() float64 {
	return dto.USD
}

func (dto *CoinDTO) IsBitcoin() bool {
	return dto.Currency == "BTC"
}

func (dto *CoinDTO) String() string {
	return fmt.Sprintf("[CoinDTO] Currency: %s, Balance: %f, Available: %f, Pending: %f, Price: %f, Address: %s, Total: %f, BTC: %f, USD: %f",
		dto.Currency, dto.Balance, dto.Available, dto.Pending, dto.Price, dto.Address, dto.Total, dto.BTC, dto.USD)
}
