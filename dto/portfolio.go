package dto

import "github.com/jeremyhahn/tradebot/common"

type PortfolioDTO struct {
	User      common.User             `json:"user"`
	NetWorth  float64                 `json:"net_worth"`
	Exchanges []common.CryptoExchange `json:"exchanges"`
	Wallets   []common.CryptoWallet   `json:"wallets"`
	Tokens    []common.EthereumToken  `json:"tokens"`
	common.Portfolio
}

func NewPortfolioDTO() *PortfolioDTO {
	return &PortfolioDTO{}
}

func (dto *PortfolioDTO) GetUser() common.User {
	return dto.User
}

func (dto *PortfolioDTO) GetNetWorth() float64 {
	return dto.NetWorth
}

func (dto *PortfolioDTO) GetExchanges() []common.CryptoExchange {
	return dto.Exchanges
}

func (dto *PortfolioDTO) GetWallets() []common.CryptoWallet {
	return dto.Wallets
}

func (dto *PortfolioDTO) GetTokens() []common.EthereumToken {
	return dto.Tokens
}
