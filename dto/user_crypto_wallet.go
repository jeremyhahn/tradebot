package dto

import "github.com/jeremyhahn/tradebot/common"

type UserCryptoWalletDTO struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
	common.UserCryptoWallet
}

func NewCryptoWallet() common.UserCryptoWallet {
	return &UserCryptoWalletDTO{}
}

func (dto *UserCryptoWalletDTO) GetAddress() string {
	return dto.Address
}

func (dto *UserCryptoWalletDTO) GetBalance() float64 {
	return dto.Balance
}

func (dto *UserCryptoWalletDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *UserCryptoWalletDTO) GetValue() float64 {
	return dto.Value
}
