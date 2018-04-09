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

type CashExplorerResponse struct {
	Txs []CashExplorerTx `json:"txs"`
}

type CashExplorerTx struct {
	TxID     string               `json:"txid"`
	Time     int64                `json:"time"`
	Hash     string               `json:"blockhash"`
	ValueIn  float64              `json:"valueIn"`
	ValueOut float64              `json:"valueOut"`
	Fees     float64              `json:"fees"`
	Outputs  []CashExplorerOutput `json:"vout"`
	Inputs   []CashExplorerInput  `json:"vin"`
}

type CashExplorerOutput struct {
	Value     string   `json:"value"`
	Addresses []string `json:"addresses"`
	SpentTxID string   `json:"spentTxId"`
}

type CashExplorerInput struct {
	TxId    string  `json:"txid"`
	Address string  `json:"addr"`
	Value   float64 `json:"value"`
}

type BchWallet struct {
	logger   *logging.Logger
	client   http.Client
	params   *common.WalletParams
	endpoint string
	common.Wallet
}

var BCHWALLET_RATELIMITER = common.NewRateLimiter(1, 1)

func CreateBchWallet(params *common.WalletParams) common.Wallet {
	client := http.Client{Timeout: common.HTTP_CLIENT_TIMEOUT}
	return &BchWallet{
		logger:   params.Context.GetLogger(),
		client:   client,
		params:   params,
		endpoint: "https://cashexplorer.bitcoin.com/api"}
}

func (bch *BchWallet) GetPrice() decimal.Decimal {
	BCHWALLET_RATELIMITER.RespectRateLimit()
	marketcap := bch.params.MarketCapService.GetMarket("BCH")
	usd, err := decimal.NewFromString(marketcap.PriceUSD)
	if err != nil {
		bch.params.Context.GetLogger().Errorf("[BchWallet.GetPrice] Error parsing price string into decimal: %s", err.Error())
		return decimal.NewFromFloat(0)
	}
	return usd.Truncate(2)
}

func (bch *BchWallet) GetWallet() (common.UserCryptoWallet, error) {
	BCHWALLET_RATELIMITER.RespectRateLimit()
	if len(bch.params.Address) <= 0 {
		return nil, errors.New("BitcoinCash address is nil")
	}
	url := fmt.Sprintf("%s/tx/%s?format=json", bch.endpoint, bch.params.Address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		bch.logger.Errorf("[BchWallet.GetBalance] %s", err.Error())
		return nil, err
	}
	wallet := BlockchainWallet{}
	jsonErr := json.Unmarshal(body, &wallet)
	if jsonErr != nil {
		bch.logger.Errorf("[BchWallet.GetBalance] %s", jsonErr.Error())
		return nil, err
	}
	balance := decimal.NewFromFloat(wallet.Balance).Div(decimal.NewFromFloat(100000000))
	return &dto.UserCryptoWalletDTO{
		Address:  bch.params.Address,
		Balance:  balance.Truncate(8),
		Currency: "BCH",
		Value:    balance.Mul(bch.GetPrice()).Truncate(2)}, nil
}

func (bch *BchWallet) GetTransactions() ([]common.Transaction, error) {
	BCHWALLET_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s/txs?address=%s", bch.endpoint, bch.params.Address)
	bch.logger.Debugf("[BchWallet.GetTransactions] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		return nil, err
	}

	var response CashExplorerResponse
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		bch.logger.Errorf("[BchWallet.GetTransactions] Error: %s", jsonErr.Error())
	}

	var transactions []common.Transaction

	if len(response.Txs) <= 0 {
		return transactions, nil
	}

	currencyPair := &common.CurrencyPair{
		Base:          "BCH",
		Quote:         "BCH",
		LocalCurrency: bch.params.Context.GetUser().GetLocalCurrency()}

	for _, tx := range response.Txs {
		timestamp := time.Unix(tx.Time, 0)
		candlestick, err := bch.params.FiatPriceService.GetPriceAt("BCH", timestamp)
		if err != nil {
			return nil, err
		}
		for _, input := range tx.Inputs {
			if input.Address == bch.params.Address {
				quantity := decimal.NewFromFloat(input.Value)
				fee := quantity.Div(decimal.NewFromFloat(tx.Fees))
				fiatTotal := quantity.Mul(candlestick.Close)
				transactions = append(transactions, &dto.TransactionDTO{
					Id:                   input.TxId,
					Date:                 timestamp,
					CurrencyPair:         currencyPair,
					Type:                 common.WITHDRAWAL_ORDER_TYPE,
					Network:              "BitcoinCash",
					NetworkDisplayName:   "Bitcoin Cash",
					Quantity:             quantity.StringFixed(8),
					QuantityCurrency:     "BCH",
					FiatQuantity:         fiatTotal.StringFixed(2),
					FiatQuantityCurrency: "USD",
					Price:                candlestick.Close.StringFixed(2),
					PriceCurrency:        "USD",
					FiatPrice:            candlestick.Close.StringFixed(2),
					FiatPriceCurrency:    "USD",
					Fee:                  fee.StringFixed(8),
					FeeCurrency:          "BCH",
					FiatFee:              fee.Mul(candlestick.Close).StringFixed(2),
					FiatFeeCurrency:      "USD",
					Total:                quantity.StringFixed(8),
					TotalCurrency:        "BCH",
					FiatTotal:            fiatTotal.StringFixed(2),
					FiatTotalCurrency:    "USD"})
			}
		}
		for _, output := range tx.Outputs {
			for _, address := range output.Addresses {
				if address == bch.params.Address {
					quantity, err := decimal.NewFromString(output.Value)
					if err != nil {
						bch.logger.Errorf("[BchWallet.GetTransactions] Error parsing value from string to decimal: %s", err.Error())
						return transactions, err
					}
					fee := quantity.Div(decimal.NewFromFloat(tx.Fees))
					fiatTotal := quantity.Mul(candlestick.Close)
					transactions = append(transactions, &dto.TransactionDTO{
						Id:                   output.SpentTxID,
						Date:                 timestamp,
						CurrencyPair:         currencyPair,
						Type:                 common.DEPOSIT_ORDER_TYPE,
						Network:              "BitcoinCash",
						NetworkDisplayName:   "Bitcoin Cash",
						Quantity:             quantity.StringFixed(8),
						QuantityCurrency:     "BCH",
						FiatQuantity:         fiatTotal.StringFixed(2),
						FiatQuantityCurrency: "USD",
						Price:                candlestick.Close.StringFixed(2),
						PriceCurrency:        "USD",
						FiatPrice:            candlestick.Close.StringFixed(2),
						FiatPriceCurrency:    "USD",
						Fee:                  fee.StringFixed(8),
						FeeCurrency:          "BCH",
						FiatFee:              fee.Mul(candlestick.Close).StringFixed(2),
						FiatFeeCurrency:      "USD",
						Total:                quantity.StringFixed(8),
						TotalCurrency:        "BCH",
						FiatTotal:            fiatTotal.StringFixed(2),
						FiatTotalCurrency:    "USD"})
				}
			}
		}
	}

	return transactions, nil
}
