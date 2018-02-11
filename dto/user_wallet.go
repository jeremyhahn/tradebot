package dto

import "github.com/jeremyhahn/tradebot/common"

type UserWalletDTO struct {
	UserId   uint   `json:"user_id"`
	Currency string `json:"currency"`
	Address  string `json:"address"`
	common.UserWallet
}

func (dto *UserWalletDTO) GetUserId() uint {
	return dto.UserId
}

func (dto *UserWalletDTO) GetCurrency() string {
	return dto.Currency
}

func (dto *UserWalletDTO) GetAddress() string {
	return dto.Address
}
