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
	Addr    string `json:"addr"`
	Value   int64  `json:"value"`
	TxIndex int64  `json:"tx_index"`
}

type BitcoinRawAddrTxInput struct {
	PrevOut BitcoinRawAddrTxItem `json:"prev_out"`
}

type BitcoinRawAddrTx struct {
	Inputs []BitcoinRawAddrTxInput `json:"inputs"`
	Out    []BitcoinRawAddrTxItem  `json:"out"`
	Time   int64                   `json:"time"`
	Hash   string                  `json:"hash"`
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
	logger   *logging.Logger
	client   http.Client
	params   *common.WalletParams
	endpoint string
	common.Wallet
}

var BTCWALLET_RATELIMITER = common.NewRateLimiter(1, 1)

func main() {}

func CreateBtcWallet(params *common.WalletParams) common.Wallet {
	client := http.Client{Timeout: common.HTTP_CLIENT_TIMEOUT}
	return &BtcWallet{
		logger:   params.Context.GetLogger(),
		client:   client,
		params:   params,
		endpoint: "https://blockchain.info"}
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

	var transactions []common.Transaction

	if len(bitcoinRawAddr.Txs) <= 0 {
		return transactions, nil
	}

	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "BTC",
		LocalCurrency: b.params.Context.GetUser().GetLocalCurrency()}

	for _, tx := range bitcoinRawAddr.Txs {

		var txType string
		var quantity, fee decimal.Decimal
		var myOutSum, outSum, myInSum, inSum int64
		for _, out := range tx.Out {
			if out.Addr == address {
				myOutSum += out.Value
			} else {
				outSum += out.Value
			}
		}
		for _, input := range tx.Inputs {
			if input.PrevOut.Addr == address {
				myInSum += input.PrevOut.Value
			} else {
				inSum += input.PrevOut.Value
			}
		}
		zero := decimal.NewFromFloat(0)
		myInputSum := decimal.NewFromFloat(float64(myInSum)).Div(decimal.NewFromFloat(100000000))
		//inputSum := decimal.NewFromFloat(float64(inSum)).Div(decimal.NewFromFloat(100000000))
		myOutputSum := decimal.NewFromFloat(float64(myOutSum)).Div(decimal.NewFromFloat(100000000))
		outputSum := decimal.NewFromFloat(float64(outSum)).Div(decimal.NewFromFloat(100000000))

		//fmt.Printf("myInputSum=%s, inputSum=%s, myOutputSum=%s, outputSum=%s\n", myInputSum.StringFixed(8),
		//	inputSum.StringFixed(8), myOutputSum.StringFixed(8), outputSum.StringFixed(8))

		if myOutputSum.GreaterThan(zero) {
			txType = common.TX_CATEGORY_DEPOSIT
			quantity = myOutputSum
			fee = zero
		} else if myInputSum.GreaterThan(zero) {
			txType = common.TX_CATEGORY_WITHDRAWAL
			quantity = myInputSum
			fee = myInputSum.Sub(outputSum)
		}
		timestamp := time.Unix(tx.Time, 0)
		candlestick, err := b.params.FiatPriceService.GetPriceAt("BTC", timestamp)
		if err != nil {
			return nil, err
		}
		fiatFee := fee.Mul(candlestick.Close)
		fiatTotal := quantity.Mul(candlestick.Close).Add(fiatFee)
		transactions = append(transactions, &dto.TransactionDTO{
			Id:                     tx.Hash,
			Date:                   timestamp,
			MarketPair:             currencyPair,
			CurrencyPair:           currencyPair,
			Type:                   txType,
			Category:               txType,
			Network:                "Bitcoin",
			NetworkDisplayName:     "Bitcoin",
			Quantity:               quantity.StringFixed(8),
			QuantityCurrency:       "BTC",
			FiatQuantity:           fiatTotal.Sub(fiatFee).StringFixed(2),
			FiatQuantityCurrency:   "USD",
			Price:                  candlestick.Close.StringFixed(2),
			PriceCurrency:          "USD",
			FiatPrice:              candlestick.Close.StringFixed(2),
			FiatPriceCurrency:      "USD",
			QuoteFiatPrice:         candlestick.Close.StringFixed(2),
			QuoteFiatPriceCurrency: "USD",
			Fee:               fee.StringFixed(8),
			FeeCurrency:       "BTC",
			FiatFee:           fiatFee.StringFixed(2),
			FiatFeeCurrency:   "USD",
			Total:             quantity.StringFixed(8),
			TotalCurrency:     "BTC",
			FiatTotal:         fiatTotal.StringFixed(2),
			FiatTotalCurrency: "USD"})
	}

	newOffset := offset
	numPages := bitcoinRawAddr.NumberOfTxs / len(bitcoinRawAddr.Txs)

	for i := 1; i < numPages; i++ {
		BTCWALLET_RATELIMITER.RespectRateLimit()
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
