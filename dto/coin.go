package dto

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type CoinDTO struct {
	Currency    string          `json:"currency"`
	Exchange    string          `json:"exchange"`
	Balance     decimal.Decimal `json:"balance"`
	Available   decimal.Decimal `json:"available"`
	Pending     decimal.Decimal `json:"pending"`
	Address     string          `json:"address"`
	Total       decimal.Decimal `json:"total"`
	Price       decimal.Decimal `json:"price"`
	BTC         decimal.Decimal `json:"btc"`
	USD         decimal.Decimal `json:"usd"`
	common.Coin `json:"-"`
}

func NewCoinDTO() common.Coin {
	return &CoinDTO{}
}

func (dto *CoinDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *CoinDTO) GetPrice() decimal.Decimal {
	return dto.Price
}

func (dto *CoinDTO) GetExchange() string {
	return dto.Exchange
}

func (dto *CoinDTO) GetBalance() decimal.Decimal {
	return dto.Balance
}

func (dto *CoinDTO) GetAvailable() decimal.Decimal {
	return dto.Available
}

func (dto *CoinDTO) GetPending() decimal.Decimal {
	return dto.Pending
}

func (dto *CoinDTO) GetAddress() string {
	return dto.Address
}

func (dto *CoinDTO) GetTotal() decimal.Decimal {
	return dto.Total
}

func (dto *CoinDTO) GetBTC() decimal.Decimal {
	return dto.BTC
}

func (dto *CoinDTO) GetUSD() decimal.Decimal {
	return dto.USD
}

func (dto *CoinDTO) IsBitcoin() bool {
	return dto.Currency == "BTC"
}

func (dto *CoinDTO) String() string {
	return fmt.Sprintf("[CoinDTO] Currency: %s, Balance: %s, Available: %s, Pending: %s, Price: %s, Address: %s, Total: %s, BTC: %s, USD: %s",
		dto.Currency, dto.Balance, dto.Available, dto.Pending, dto.Price, dto.Address, dto.Total, dto.BTC, dto.USD)
}
