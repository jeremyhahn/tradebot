package dto

import "github.com/jeremyhahn/tradebot/common"

type CryptoWalletDTO struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	NetWorth float64 `json:"net_worth"`
	common.CryptoWallet
}

func NewCryptoWallet() common.CryptoWallet {
	return &CryptoWalletDTO{}
}

func (dto *CryptoWalletDTO) GetAddress() string {
	return dto.Address
}

func (dto *CryptoWalletDTO) GetBalance() float64 {
	return dto.Balance
}

func (dto *CryptoWalletDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *CryptoWalletDTO) GetNetWorth() float64 {
	return dto.NetWorth
}
