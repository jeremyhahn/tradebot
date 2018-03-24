package dto

import "github.com/jeremyhahn/tradebot/common"

type UserCryptoExchangeDTO struct {
	UserID                    uint   `json:"-"`
	Name                      string `json:"name"`
	Key                       string `json:"key"`
	Secret                    string `json:"secret"`
	Extra                     string `json:"extra"`
	common.UserCryptoExchange `json:"-"`
}

func (dto *UserCryptoExchangeDTO) GetUserID() uint {
	return dto.UserID
}

func (dto *UserCryptoExchangeDTO) GetName() string {
	return dto.Name
}

func (dto *UserCryptoExchangeDTO) GetKey() string {
	return dto.Key
}

func (dto *UserCryptoExchangeDTO) GetSecret() string {
	return dto.Secret
}

func (dto *UserCryptoExchangeDTO) GetExtra() string {
	return dto.Extra
}
