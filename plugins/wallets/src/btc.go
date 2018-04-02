package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
	logging "github.com/op/go-logging"
	"github.com/shopspring/decimal"
)

type BitcoinRawAddrTxItem struct {
	Addr  string `json:"addr"`
	Value int64  `json:"value"`
}

type BitcoinRawAddrTxInput struct {
	PrevOut BitcoinRawAddrTxItem `json:"prev_out"`
}

type BitcoinRawAddrTx struct {
	Inputs []BitcoinRawAddrTxInput `json:"inputs"`
	Out    []BitcoinRawAddrTxItem  `json:"out"`
	Time   int64                   `json:"time"`
}

type BitcoinRawAddr struct {
	Txs         []BitcoinRawAddrTx `json:"txs"`
	NumberOfTxs int                `json:"n_tx"`
}

type BlockchainWallet struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"final_balance"`
	Currency string  `json:"currency"`
}

type BlockchainTickerSubItem struct {
	Last float64 `json:"last"`
}

type BlockchainTickerItem struct {
	USD BlockchainTickerSubItem
}

type BtcWallet struct {
	logger    *logging.Logger
	client    http.Client
	items     BlockchainTickerItem
	lastPrice decimal.Decimal
	params    *common.WalletParams
	endpoint  string
	common.Wallet
}

var BTCWALLET_RATELIMITER = common.NewRateLimiter(1, 1)

func main() {}

func CreateBtcWallet(params *common.WalletParams) common.Wallet {
	client := http.Client{Timeout: common.HTTP_CLIENT_TIMEOUT}
	return &BtcWallet{
		logger:    params.Context.GetLogger(),
		client:    client,
		lastPrice: decimal.NewFromFloat(0.0),
		params:    params,
		endpoint:  "https://blockchain.info"}
}

func (b *BtcWallet) GetPrice() decimal.Decimal {
	BTCWALLET_RATELIMITER.RespectRateLimit()
	_, body, err := util.HttpRequest(fmt.Sprintf("%s/ticker", b.endpoint))
	if err != nil {
		b.logger.Errorf("[BtcWallet.GetPrice] Error: %s", err.Error())
	}
	t := BlockchainTickerItem{}
	jsonErr := json.Unmarshal(body, &t)
	if jsonErr != nil {
		b.logger.Errorf("[BtcWallet.GetPrice] %s", jsonErr.Error())
	}
	return decimal.NewFromFloat(t.USD.Last).Truncate(2)
}

func (b *BtcWallet) GetWallet() (common.UserCryptoWallet, error) {
	BTCWALLET_RATELIMITER.RespectRateLimit()
	if len(b.params.Address) <= 0 {
		return nil, errors.New("Bitcoin address is nil")
	}
	url := fmt.Sprintf("%s/address/%s?format=json", b.endpoint, b.params.Address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		b.logger.Errorf("[BtcWallet.GetBalance] %s", err.Error())
		return nil, err
	}
	wallet := BlockchainWallet{}
	jsonErr := json.Unmarshal(body, &wallet)
	if jsonErr != nil {
		b.logger.Errorf("[BtcWallet.GetBalance] %s", jsonErr.Error())
		return nil, err
	}
	balance := decimal.NewFromFloat(wallet.Balance).Div(decimal.NewFromFloat(100000000))
	return &dto.UserCryptoWalletDTO{
		Address:  b.params.Address,
		Balance:  balance.Truncate(8),
		Currency: "BTC",
		Value:    balance.Mul(b.GetPrice()).Truncate(2)}, nil
}

func (b *BtcWallet) GetTransactions() ([]common.Transaction, error) {
	BTCWALLET_RATELIMITER.RespectRateLimit()
	return b.getTransactionsAt(b.params.Address, 0)
}

func (b *BtcWallet) getTransactionsAt(address string, offset int) ([]common.Transaction, error) {

	url := fmt.Sprintf("%s/rawaddr/%s?limit=%d&offset=%d", b.endpoint, address, 50, offset)
	b.logger.Debugf("[BtcWallet.GetTransactions] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		return nil, err
	}

	var bitcoinRawAddr BitcoinRawAddr
	jsonErr := json.Unmarshal(body, &bitcoinRawAddr)
	if jsonErr != nil {
		b.logger.Errorf("[BtcWallet.GetTransactions] Error: %s", jsonErr.Error())
	}

	transactions := make([]common.Transaction, 0, bitcoinRawAddr.NumberOfTxs)

	if len(bitcoinRawAddr.Txs) <= 0 {
		return transactions, nil
	}

	for _, tx := range bitcoinRawAddr.Txs {
		for _, input := range tx.Inputs {
			if input.PrevOut.Addr == address {
				transactions = append(transactions, &dto.TransactionDTO{
					Date:     time.Unix(tx.Time, 0),
					Type:     common.WITHDRAWAL_ORDER_TYPE,
					Quantity: decimal.New(input.PrevOut.Value, 8).StringFixed(8)})
			}
		}
		for _, out := range tx.Out {
			if out.Addr == address {
				transactions = append(transactions, &dto.TransactionDTO{
					Date:     time.Unix(tx.Time, 0),
					Type:     common.DEPOSIT_ORDER_TYPE,
					Quantity: decimal.New(out.Value, 8).StringFixed(8)}) // 100000000
			}
		}
	}

	newOffset := offset
	numPages := bitcoinRawAddr.NumberOfTxs / len(bitcoinRawAddr.Txs)

	for i := 1; i < numPages; i++ {
		newOffset += 50
		txs, err := b.getTransactionsAt(address, newOffset)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txs...)
		return transactions, nil
	}

	return transactions, nil
}
