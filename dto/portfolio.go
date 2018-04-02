package dto

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type PortfolioDTO struct {
	User             common.UserContext             `json:"user"`
	NetWorth         decimal.Decimal                `json:"net_worth"`
	Exchanges        []common.CryptoExchangeSummary `json:"exchanges"`
	Wallets          []common.UserCryptoWallet      `json:"wallets"`
	Tokens           []common.EthereumToken         `json:"tokens"`
	common.Portfolio `json:"-"`
}

func NewPortfolioDTO() *PortfolioDTO {
	return &PortfolioDTO{}
}

func (dto *PortfolioDTO) GetUser() common.UserContext {
	return dto.User
}

func (dto *PortfolioDTO) GetNetWorth() decimal.Decimal {
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
