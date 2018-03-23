package dto

import "github.com/jeremyhahn/tradebot/common"

type UserDTO struct {
	Id            uint                        `json:"id"`
	Username      string                      `json:"username"`
	LocalCurrency string                      `json:"local_currency"`
	FiatExchange  string                      `json:"fiat_exchange"`
	Etherbase     string                      `json:"etherbase"`
	Keystore      string                      `json:"keystore"`
	Charts        []common.Chart              `json:"charts"`
	Exchanges     []common.UserCryptoExchange `json:"exchanges"`
	Wallets       []common.UserCryptoWallet   `json:"wallets"`
	Tokens        []common.EthereumToken      `json:"tokens"`
	common.UserContext
}

func NewUserDTO() common.UserContext {
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

func (dto *UserDTO) GetFiatExchange() string {
	return dto.FiatExchange
}

func (dto *UserDTO) GetEtherbase() string {
	return dto.Etherbase
}

func (dto *UserDTO) GetKeystore() string {
	return dto.Keystore
}

func (dto *UserDTO) GetCharts() []common.Chart {
	return dto.Charts
}

func (dto *UserDTO) GetExchanges() []common.UserCryptoExchange {
	return dto.Exchanges
}

func (dto *UserDTO) GetWallets() []common.UserCryptoWallet {
	return dto.Wallets
}

func (dto *UserDTO) GetTokens() []common.EthereumToken {
	return dto.Tokens
}
