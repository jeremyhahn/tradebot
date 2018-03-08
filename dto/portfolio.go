package dto

import "github.com/jeremyhahn/tradebot/common"

type PortfolioDTO struct {
	User      common.UserContext             `json:"user"`
	NetWorth  float64                        `json:"net_worth"`
	Exchanges []common.CryptoExchangeSummary `json:"exchanges"`
	Wallets   []common.UserCryptoWallet      `json:"wallets"`
	Tokens    []common.EthereumToken         `json:"tokens"`
	common.Portfolio
}

func NewPortfolioDTO() *PortfolioDTO {
	return &PortfolioDTO{}
}

func (dto *PortfolioDTO) GetUser() common.UserContext {
	return dto.User
}

func (dto *PortfolioDTO) GetNetWorth() float64 {
	return dto.NetWorth
}

func (dto *PortfolioDTO) GetExchanges() []common.CryptoExchangeSummary {
	return dto.Exchanges
}

func (dto *PortfolioDTO) GetWallets() []common.UserCryptoWallet {
	return dto.Wallets
}

func (dto *PortfolioDTO) GetTokens() []common.EthereumToken {
	return dto.Tokens
}
