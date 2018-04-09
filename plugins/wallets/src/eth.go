package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
)

type EtherScanTx struct {
	Hash            string `json:"hash"`
	Timestamp       string `json:"timestamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Gas             string `json:"gas"`
	GasPrice        string `json:"gasPrice"`
	GasUsed         string `json:"gasUsed"`
	ContractAddress string `json:"contractAddress"`
}

type EtherScanGetTransactionsResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []EtherScanTx `json:"result"`
}

type EtherScanGetLastPrice struct {
	//BTC          string `json:"ethbtc"`
	//BTCTimestamp string `json:"ethbtc_timestamp"`
	USD string `json:"ethusd"`
	//USDTimestamp string `json:"ethusd_timestamp"`
}

type EtherScanGetLastPriceResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Result  EtherScanGetLastPrice `json:"result"`
}

type EtherScanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type EthWallet struct {
	ctx              common.Context
	marketcapService common.MarketCapService
	fiatPriceService common.FiatPriceService
	apiKeyToken      string
	endpoint         string
	params           *common.WalletParams
	common.Wallet
}

var ETHWALLET_RATELIMITER = common.NewRateLimiter(5, 1)

func CreateEthWallet(params *common.WalletParams) common.Wallet {
	return &EthWallet{
		ctx:              params.Context,
		marketcapService: params.MarketCapService,
		fiatPriceService: params.FiatPriceService,
		apiKeyToken:      params.WalletSecret,
		endpoint:         "https://api.etherscan.io/api",
		params:           params}
}

func (service *EthWallet) GetPrice() decimal.Decimal {

	ETHWALLET_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=stats&action=ethprice&apikey=%s", service.endpoint, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EthWallet.GetPrice] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetPrice] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetLastPriceResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetPrice] JSON Unmarshal Error: %s", jsonErr.Error())
	}

	price, err := decimal.NewFromString(response.Result.USD)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetPrice] String to decimal conversation error: %s", err.Error())
	}

	return price.Truncate(2)
}

func (service *EthWallet) GetWallet() (common.UserCryptoWallet, error) {

	ETHWALLET_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=balance&address=0x%s&tag=latest&apikey=%s",
		service.endpoint, service.params.Address, service.params.WalletSecret)

	service.ctx.GetLogger().Debugf("[EthWallet.GetWallet] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetWallet] Error: %s", err.Error())
		return nil, err
	}

	response := EtherScanResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetWallet] %s", jsonErr.Error())
		return nil, jsonErr
	}

	result, err := decimal.NewFromString(response.Result)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetWallet] Float conversation error: %s", err.Error())
		return nil, err
	}

	balance := result.Div(decimal.NewFromFloat(1000000000000000000))

	return &dto.UserCryptoWalletDTO{
		Address:  service.params.Address,
		Balance:  balance.Truncate(8),
		Currency: "ETH",
		Value:    service.GetPrice().Mul(balance).Truncate(2)}, nil
}

func (service *EthWallet) GetTransactions() ([]common.Transaction, error) {
	ETHWALLET_RATELIMITER.RespectRateLimit()
	url := fmt.Sprintf("%s?module=account&action=txlist&address=0x%s&startblock=0&endblock=99999999999999999999999&sort=asc&apikey=%s",
		service.endpoint, service.params.Address, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EthWallet.GetTransactionsFor] url: %s", url)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetTransactionsFor] HTTP Request Error: %s", err.Error())
	}
	response := EtherScanGetTransactionsResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthWallet.GetTransactionsFor] JSON Unmarshal error: %s", jsonErr.Error())
	}
	var transactions []common.Transaction
	for _, tx := range response.Result {
		var txType string
		hexAddress := fmt.Sprintf("0x%s", service.params.Address)
		if tx.From == hexAddress {
			txType = common.WITHDRAWAL_ORDER_TYPE
		} else if tx.To == hexAddress {
			txType = common.DEPOSIT_ORDER_TYPE
		}
		timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		amount, err := decimal.NewFromString(tx.Value)
		if err != nil {
			return nil, err
		}
		fee, err := decimal.NewFromString(tx.GasUsed)
		if err != nil {
			return nil, err
		}
		localCurrency := service.ctx.GetUser().GetLocalCurrency()
		finalAmount := amount.Div(decimal.NewFromFloat(1000000000000000000))
		finalFee := fee.Div(decimal.NewFromFloat(100000000))
		candlestick, err := service.fiatPriceService.GetPriceAt("ETH", time.Unix(timestamp, 0))
		if err != nil {
			return nil, err
		}
		fiatFee := finalFee.Mul(candlestick.Close)
		fiatTotal := finalAmount.Mul(candlestick.Close).Add(fiatFee)
		transactions = append(transactions, &dto.TransactionDTO{
			Id:                   tx.Hash,
			Date:                 time.Unix(timestamp, 0),
			CurrencyPair:         &common.CurrencyPair{Base: "ETH", Quote: "ETH", LocalCurrency: localCurrency},
			Type:                 txType,
			Category:             common.TX_CATEGORY_TRANSFER,
			Network:              "etherscan",
			NetworkDisplayName:   "Ethereum",
			Quantity:             finalAmount.StringFixed(8),
			QuantityCurrency:     "ETH",
			FiatQuantity:         finalAmount.Mul(candlestick.Close).StringFixed(2),
			FiatQuantityCurrency: "USD",
			Price:                candlestick.Close.StringFixed(2),
			PriceCurrency:        "USD",
			FiatPrice:            candlestick.Close.StringFixed(2),
			FiatPriceCurrency:    "USD",
			Fee:                  finalFee.StringFixed(8),
			FeeCurrency:          "ETH",
			FiatFee:              fiatFee.StringFixed(2),
			FiatFeeCurrency:      "USD",
			Total:                finalAmount.StringFixed(8),
			TotalCurrency:        "ETH",
			FiatTotal:            fiatTotal.StringFixed(2),
			FiatTotalCurrency:    "USD"})
	}
	return transactions, nil
}
