// Not used - leaving just in case

package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
)

type BlockCypherService interface {
	BlockExplorerService
}

type EthereumTxRef struct {
	Date          time.Time `json:"confirmed"`
	TxInput       int       `json:"tx_input_n"`
	TxOutput      int       `json:"tx_output_n"`
	Value         float64   `json:"value"`
	Confirmations int       `json:"confirmations"`
	IsDoubleSpend bool      `json:"double_spend"`
	BlockHeight   int       `json:"block_height"`
}

type EthereumAccount struct {
	Balance     float64         `json:"balance"`
	NumberOfTxs int             `json:"n_tx"`
	TxRefs      []EthereumTxRef `json:"txrefs"`
}

type BlockCypherServiceImpl struct {
	ctx              common.Context
	endpoint         string
	marketcapService MarketCapService
	BlockCypherService
}

func NewBlockCypherService(ctx common.Context, marketcapService MarketCapService) BlockCypherService {
	return &BlockCypherServiceImpl{
		ctx:              ctx,
		marketcapService: marketcapService,
		endpoint:         "https://api.blockcypher.com/v1/eth/main"}
}

func (service *BlockCypherServiceImpl) GetBalance(address string) common.UserCryptoWallet {
	url := fmt.Sprintf("%s/addrs/%s", service.endpoint, address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] Error: %s", err.Error())
	}
	ethereumAccount := EthereumAccount{}
	jsonErr := json.Unmarshal(body, &ethereumAccount)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] %s", jsonErr.Error())
	}
	priceUSD, err := strconv.ParseFloat(service.marketcapService.GetMarket("ETH").PriceUSD, 64)
	if err != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] Error: %s", err.Error())
	}
	balance := ethereumAccount.Balance / 1000000000000000000
	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  balance,
		Currency: "ETH",
		Value:    priceUSD * balance}
}

func (service *BlockCypherServiceImpl) GetTransactions(address string) ([]common.Transaction, error) {
	/*
		url := fmt.Sprintf("%s/addrs/%s?limit=2000", service.endpoint, address)
		body, err := util.HttpRequest(url)
		if err != nil {
			service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] Error: %s", err.Error())
			return nil, err
		}
		ethereumAccount := EthereumAccount{}
		jsonErr := json.Unmarshal(body, &ethereumAccount)
		if jsonErr != nil {
			service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] %s", jsonErr.Error())
			return nil, err
		}
		transactions := make([]common.Transaction, 0, ethereumAccount.NumberOfTxs)
		for _, txref := range ethereumAccount.TxRefs {
			if txref.Confirmations <= 0 || txref.IsDoubleSpend {
				continue
			}
			var txType string
			if txref.TxInput == -1 {
				txType = "deposit"
			} else if txref.TxOutput == -1 {
				txType = "withdrawl"
			}
			transactions = append(transactions, &dto.TransactionDTO{
				Date:   txref.Date,
				Amount: txref.Value / 1000000000000000000,
				Type:   txType})
		}
		return transactions, nil
	*/
	return service.getTransactionsAt(address, 0)
}

func (service *BlockCypherServiceImpl) GetTokenTransactions(walletAddress, contractAddress string) ([]common.Transaction, error) {
	txs, _, err := service.getTokenTransactionsAt(walletAddress, contractAddress, 0)
	return txs, err
}

/*
func (service *BlockCypherServiceImpl) GetContract(address string) (common.EthereumContract, error) {
	url := fmt.Sprintf("%s/contracts/%s", service.endpoint, address)
	body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] Error: %s", err.Error())
		return nil, err
	}

	fmt.Println(string(body))

	ethereumContract := &dto.EthereumContractDTO{}
	jsonErr := json.Unmarshal(body, ethereumContract)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.GetBalance] %s", jsonErr.Error())
		return nil, err
	}
	return ethereumContract, nil
}*/

func (service *BlockCypherServiceImpl) getTokenTransactionsAt(walletAddress, contractAddress string, blockHeight int) ([]common.Transaction, int, error) {

	service.ctx.GetLogger().Debugf("[BlockCypherService.getTokenTransactionsAt] walletAddress: %s, contractAddress: %s, blockHeight: %d",
		walletAddress, contractAddress, blockHeight)

	url := fmt.Sprintf("%s/addrs/%s?after=%d&confirmations=1&limit=2000", service.endpoint, contractAddress, blockHeight)

	service.ctx.GetLogger().Debugf("[BlockCypherService.getTokenTransactionsAt] url: %s", url)

	statusCode, body, err := util.HttpRequest(url)

	service.ctx.GetLogger().Debugf("%d", statusCode)
	service.ctx.GetLogger().Debugf("%s", string(body))

	os.Exit(1)

	if err != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTokenTransactionsAt] Error: %s", err.Error())
		return nil, blockHeight, err
	}
	if statusCode == 429 {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTokenTransactionsAt] %s", string(body))
		return nil, blockHeight, err
	}
	ethereumAccount := EthereumAccount{}
	jsonErr := json.Unmarshal(body, &ethereumAccount)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTokenTransactionsAt] %s", jsonErr.Error())
		return nil, blockHeight, err
	}

	lastBlockHeight := blockHeight
	transactions := make([]common.Transaction, 0, ethereumAccount.NumberOfTxs)
	txLen := len(ethereumAccount.TxRefs)

	if txLen <= 0 {
		return transactions, lastBlockHeight, nil
	}

	for _, txref := range ethereumAccount.TxRefs {
		if txref.IsDoubleSpend {
			continue
		}

		//service.ctx.GetLogger().Debugf("%+v", txref)

		var txType string
		if txref.TxInput == -1 {
			txType = "deposit"
		} else if txref.TxOutput == -1 {
			txType = "withdrawl"
		}
		lastBlockHeight = txref.BlockHeight
		transactions = append(transactions, &dto.TransactionDTO{
			Date:   txref.Date,
			Amount: txref.Value / 1000000000000000000,
			Type:   txType})
	}

	if lastBlockHeight == blockHeight && lastBlockHeight > 0 {
		return transactions, lastBlockHeight, nil
	}

	service.ctx.GetLogger().Debugf("[BlockCypher.getTokenTransactionsAt] lastBlockHeight: %d,  blockHeight: %d", lastBlockHeight, blockHeight)
	fmt.Printf("txLen: %d\n", txLen)

	service.ctx.GetLogger().Debugf("%s\n", transactions)

	numPages := ethereumAccount.NumberOfTxs / txLen

	fmt.Printf("total # of transactions: %d\n", ethereumAccount.NumberOfTxs)
	fmt.Printf("lastBlockHeight: %d\n", lastBlockHeight)
	fmt.Printf("numPages: %d\n", numPages)

	for i := 0; i < numPages; i++ {

		fmt.Println(i)

		txs, bHeight, err := service.getTokenTransactionsAt(walletAddress, contractAddress, lastBlockHeight)
		lastBlockHeight = bHeight

		fmt.Printf("lastBlockHeight: %d", lastBlockHeight)

		//service.ctx.GetLogger().Debugf("%s\n", txs)

		if err != nil {
			return nil, lastBlockHeight, err
		}
		transactions = append(transactions, txs...)
	}

	fmt.Printf("parsed transactions: %d\n", len(transactions))

	return transactions, lastBlockHeight, nil
}

func (service *BlockCypherServiceImpl) getTransactionsAt(address string, blockHeight int) ([]common.Transaction, error) {

	service.ctx.GetLogger().Debugf("[BlockCypherService.getTransactionsAt] address: %s, blockHeight: %d", address, blockHeight)

	url := fmt.Sprintf("%s/addrs/%s?limit=2000", service.endpoint, address)

	service.ctx.GetLogger().Debugf("[BlockCypherService.getTransactionsAt] url: %s", url)

	statusCode, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTransactionsAt] Error: %s", err.Error())
		return nil, err
	}
	ethereumAccount := EthereumAccount{}
	jsonErr := json.Unmarshal(body, &ethereumAccount)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTransactionsAt] %s", jsonErr.Error())
		return nil, err
	}
	if statusCode == 429 {
		service.ctx.GetLogger().Errorf("[BlockCypherService.getTokenTransactionsAt] %s", string(body))
		return nil, err
	}
	transactions := make([]common.Transaction, 0, ethereumAccount.NumberOfTxs)
	if len(ethereumAccount.TxRefs) <= 0 {
		return transactions, nil
	}
	for _, txref := range ethereumAccount.TxRefs {
		if txref.Confirmations <= 0 || txref.IsDoubleSpend {
			continue
		}
		var txType string
		if txref.TxInput == -1 {
			txType = "deposit"
		} else if txref.TxOutput == -1 {
			txType = "withdrawl"
		}
		transactions = append(transactions, &dto.TransactionDTO{
			Date:   txref.Date,
			Amount: txref.Value / 1000000000000000000,
			Type:   txType})
	}

	newBlockHeight := blockHeight
	numPages := ethereumAccount.NumberOfTxs / len(ethereumAccount.TxRefs)

	for i := 1; i < numPages; i++ {
		newBlockHeight += 2000
		txs, err := service.getTransactionsAt(address, newBlockHeight)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txs...)
		return transactions, nil
	}

	return transactions, nil
}
