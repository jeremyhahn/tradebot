package dto

import "github.com/jeremyhahn/tradebot/common"

type UserCryptoExchangeDTO struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Extra  string `json:"extra"`
	common.UserCryptoExchange
}

func (dto *UserCryptoExchangeDTO) GetName() string {
	return dto.Name
}

func (dto *UserCryptoExchangeDTO) GetURL() string {
	return dto.URL
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
