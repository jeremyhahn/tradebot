package dto

import "github.com/jeremyhahn/tradebot/common"

type UserContextDTO struct {
	Id            uint   `json:"id"`
	Username      string `json:"username"`
	LocalCurrency string `json:"local_currency"`
	FiatExchange  string `json:"fiat_exchange"`
	Etherbase     string `json:"etherbase"`
	Keystore      string `json:"keystore"`
	common.UserContext
}

func NewUserContextDTO() common.UserContext {
	return &UserContextDTO{}
}

func (dto *UserContextDTO) GetId() uint {
	return dto.Id
}

func (dto *UserContextDTO) GetUsername() string {
	return dto.Username
}

func (dto *UserContextDTO) GetLocalCurrency() string {
	return dto.LocalCurrency
}

func (dto *UserContextDTO) GetFiatExchange() string {
	return dto.FiatExchange
}

func (dto *UserContextDTO) GetEtherbase() string {
	return dto.Etherbase
}

func (dto *UserContextDTO) GetKeystore() string {
	return dto.Keystore
}
