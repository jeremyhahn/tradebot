package dto

import "github.com/jeremyhahn/tradebot/common"

type UserDTO struct {
	Id            uint   `json:"id"`
	Username      string `json:"username"`
	LocalCurrency string `json:"local_currency"`
	Etherbase     string `json:"etherbase"`
	Keystore      string `json:"keystore"`
	common.User
}

func NewUserDTO() common.User {
	return &UserDTO{}
}

func (dto *UserDTO) GetId() uint {
	return dto.Id
}

func (dto *UserDTO) GetUsername() string {
	return dto.Username
}

func (dto *UserDTO) GetLocalCurrency() string {
	return dto.LocalCurrency
}

func (dto *UserDTO) GetEtherbase() string {
	return dto.Etherbase
}

func (dto *UserDTO) GetKeystore() string {
	return dto.Keystore
}

func (dto *UserDTO) SetEtherbase(etherbase string) {
	dto.Etherbase = etherbase
}

func (dto *UserDTO) SetKeystore(keystore string) {
	dto.Keystore = keystore
}
