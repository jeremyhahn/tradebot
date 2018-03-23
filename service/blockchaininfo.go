package service

import (
	"encoding/json"
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

type BlockchainInfo struct {
	logger     *logging.Logger
	client     http.Client
	items      BlockchainTickerItem
	lastPrice  float64
	lastLookup time.Time
	BlockExplorerService
}

func NewBlockchainInfo(ctx common.Context) *BlockchainInfo {
	client := http.Client{Timeout: common.HTTP_CLIENT_TIMEOUT}
	return &BlockchainInfo{
		logger:     ctx.GetLogger(),
		client:     client,
		lastPrice:  0.0,
		lastLookup: time.Now().Add(-20 * time.Minute)}
}

func (b *BlockchainInfo) GetPrice() float64 {
	elapsed := float64(time.Since(b.lastLookup))
	if elapsed/float64(time.Second) >= 900 {
		_, body, err := util.HttpRequest("https://blockchain.info/ticker")
		if err != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] Error: %s", err.Error())
		}
		t := BlockchainTickerItem{}
		jsonErr := json.Unmarshal(body, &t)
		if jsonErr != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] %s", jsonErr.Error())
		}
		b.lastLookup = time.Now()
		b.lastPrice = t.USD.Last
	}
	return b.lastPrice
}

func (b *BlockchainInfo) GetBalance(address string) common.UserCryptoWallet {
	url := fmt.Sprintf("https://blockchain.info/address/%s?format=json", address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", err.Error())
	}
	wallet := BlockchainWallet{}
	jsonErr := json.Unmarshal(body, &wallet)
	if jsonErr != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", jsonErr.Error())
	}
	balance := wallet.Balance / 100000000
	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  balance,
		Currency: "BTC",
		Value:    balance * b.GetPrice()}
}

func (b *BlockchainInfo) GetTransactions(address string) ([]common.Transaction, error) {
	return b.getTransactionsAt(address, 0)
}

func (b *BlockchainInfo) getTransactionsAt(address string, offset int) ([]common.Transaction, error) {

	url := fmt.Sprintf("%s/%s?limit=%d&offset=%d", "https://blockchain.info/rawaddr", address, 50, offset)

	b.logger.Debugf("[BlockchainInfo.GetTransactions] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		return nil, err
	}

	var bitcoinRawAddr BitcoinRawAddr
	jsonErr := json.Unmarshal(body, &bitcoinRawAddr)
	if jsonErr != nil {
		b.logger.Errorf("[BlockchainInfo.GetTransactions] Error: %s", jsonErr.Error())
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
					Type:     "withdrawl",
					Quantity: decimal.New(input.PrevOut.Value, 8).StringFixed(8)})
			}
		}
		for _, out := range tx.Out {
			if out.Addr == address {
				transactions = append(transactions, &dto.TransactionDTO{
					Date:     time.Unix(tx.Time, 0),
					Type:     "deposit",
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

func (b *BlockchainInfo) ConvertToUSD(currency string, btc float64) float64 {
	price := b.GetPrice()
	b.logger.Debugf("[BlockchainTicker] currency: %s, btc: %.8f, price: %f", currency, btc, price)
	return btc * price
}
