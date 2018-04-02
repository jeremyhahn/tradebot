package dto

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type UserCryptoWalletDTO struct {
	Address                 string          `json:"address"`
	Balance                 decimal.Decimal `json:"balance"`
	Currency                string          `json:"currency"`
	Value                   decimal.Decimal `json:"value"`
	common.UserCryptoWallet `json:"-"`
}

func NewCryptoWallet() common.UserCryptoWallet {
	return &UserCryptoWalletDTO{}
}

func (dto *UserCryptoWalletDTO) GetAddress() string {
	return dto.Address
}

func (dto *UserCryptoWalletDTO) GetBalance() decimal.Decimal {
	return dto.Balance
}

func (dto *UserCryptoWalletDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *UserCryptoWalletDTO) GetValue() decimal.Decimal {
	return dto.Value
}

func (dto *UserCryptoWalletDTO) String() string {
	return fmt.Sprintf("[UserCryptoWalletDTO] Currency: %s, Address: %s, Balance: %s, Value: %s",
		dto.GetCurrency(), dto.GetAddress(), dto.GetBalance(), dto.GetValue())
}
