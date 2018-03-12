package dto

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type EthereumContractDTO struct {
	Address      string
	Source       string
	Bin          string
	ABI          string
	CreationDate time.Time
	common.EthereumContract
}

func NewEthereumContract() common.EthereumContract {
	return &EthereumContractDTO{}
}

func (contract *EthereumContractDTO) GetAddress() string {
	return contract.Address
}

func (contract *EthereumContractDTO) GetSource() string {
	return contract.Source
}

func (contract *EthereumContractDTO) GetBin() string {
	return contract.Bin
}

func (contract *EthereumContractDTO) GetABI() string {
	return contract.ABI
}

func (contract *EthereumContractDTO) GetCreationDate() time.Time {
	return contract.CreationDate
}
