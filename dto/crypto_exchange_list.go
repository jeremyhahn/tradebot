package dto

import "github.com/jeremyhahn/tradebot/common"

type CryptoExchangeListDTO struct {
	Exchanges []common.CryptoExchange `json:"exchange"`
	NetWorth  float64                 `json:"net_worth"`
	common.CryptoExchangeList
}

func NewCryptoExchangeListDTO() common.CryptoExchangeList {
	return &CryptoExchangeListDTO{}
}

func (dto *CryptoExchangeListDTO) GetExchanges() []common.CryptoExchange {
	return dto.Exchanges
}

func (dto *CryptoExchangeListDTO) GetNetWorth() float64 {
	return dto.NetWorth
}
