package dto

import "github.com/jeremyhahn/tradebot/common"

type EthereumTokenDTO struct {
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Decimals        uint8   `json:"decimals"`
	Balance         float64 `json:"balance"`
	EthBalance      string  `json:"eth_balance"`
	ContractAddress string  `json:"contract_address"`
	WalletAddress   string  `json:"wallet_address"`
	Value           float64 `json:"value"`
	common.EthereumToken
}

func NewEthereumToken() common.EthereumToken {
	return &EthereumTokenDTO{}
}

func (et *EthereumTokenDTO) GetName() string {
	return et.Name
}

func (et *EthereumTokenDTO) GetSymbol() string {
	return et.Symbol
}

func (et *EthereumTokenDTO) GetDecimals() uint8 {
	return et.Decimals
}

func (et *EthereumTokenDTO) GetBalance() float64 {
	return et.Balance
}

func (et *EthereumTokenDTO) GetEthBalance() string {
	return et.EthBalance
}

func (et *EthereumTokenDTO) GetContractAddress() string {
	return et.ContractAddress
}

func (et *EthereumTokenDTO) GetWalletAddress() string {
	return et.WalletAddress
}

func (et *EthereumTokenDTO) GetValue() float64 {
	return et.Value
}
