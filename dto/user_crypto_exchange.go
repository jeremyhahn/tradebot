package dto

import "github.com/jeremyhahn/tradebot/common"

type UserCryptoExchangeDTO struct {
	UserId   uint   `json:"user_id"`
	Currency string `json:"currency"`
	Address  string `json:"address"`
	common.CryptoExchange
}
