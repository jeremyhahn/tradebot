package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
)

type EtherScanTx struct {
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

type EthereumWebClient struct {
	ctx              common.Context
	userDAO          dao.UserDAO
	authService      AuthService
	marketcapService common.MarketCapService
	fiatPriceService common.FiatPriceService
	apiKeyToken      string
	endpoint         string
	EthereumService
}

var ETHERSCAN_RATELIMITER = common.NewRateLimiter(5, 1)

func NewEthereumWebClient(ctx common.Context, userDAO dao.UserDAO, authService AuthService,
	marketcapService common.MarketCapService, fiatPriceService common.FiatPriceService) (EthereumService, error) {
	return &EthereumWebClient{
		ctx:              ctx,
		userDAO:          userDAO,
		authService:      authService,
		marketcapService: marketcapService,
		fiatPriceService: fiatPriceService,
		apiKeyToken:      "YourApiKeyToken",
		endpoint:         "https://api.etherscan.io/api"}, nil
}

func (service *EthereumWebClient) Login(username, password string) (common.UserContext, error) {
	return service.authService.Login(username, password)
}

func (service *EthereumWebClient) Register(username, password string) error {
	return service.authService.Register(username, password)
}

func (service *EthereumWebClient) GetAccounts() ([]common.UserContext, error) {
	var accounts []common.UserContext
	walletEntities := service.userDAO.GetWallets(
		&entity.User{
			Id: service.ctx.GetUser().GetId()})
	for _, wallet := range walletEntities {
		if wallet.GetCurrency() != "ETH" {
			continue
		}
		accounts = append(accounts, &dto.UserContextDTO{
			Etherbase: wallet.GetAddress(),
			Keystore:  service.apiKeyToken})
	}
	return accounts, nil
}

func (service *EthereumWebClient) GetPrice() decimal.Decimal {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=stats&action=ethprice&apikey=%s", service.endpoint, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EthereumWebClient.GetPrice] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetPrice] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetLastPriceResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetPrice] JSON Unmarshal Error: %s", jsonErr.Error())
	}

	price, err := decimal.NewFromString(response.Result.USD)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetPrice] String to decimal conversation error: %s", err.Error())
	}

	return price
}

func (service *EthereumWebClient) GetWallet(address string) (common.UserCryptoWallet, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=balance&address=0x%s&tag=latest&apikey=%s",
		service.endpoint, address, service.apiKeyToken)

	service.ctx.GetLogger().Debugf("[EthereumWebClient.GetWallet] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetWallet] Error: %s", err.Error())
		return nil, err
	}

	response := EtherScanResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetWallet] %s", jsonErr.Error())
		return nil, jsonErr
	}

	result, err := decimal.NewFromString(response.Result)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetWallet] Float conversation error: %s", err.Error())
		return nil, err
	}

	balance := result.Div(decimal.NewFromFloat(1000000000000000000))

	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  balance,
		Currency: "ETH",
		Value:    service.GetPrice().Mul(balance)}, nil
}

func (service *EthereumWebClient) GetTransactions() ([]common.Transaction, error) {
	var transactions []common.Transaction
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	wallets := service.userDAO.GetWallets(userEntity)
	for _, wallet := range wallets {
		if wallet.GetCurrency() != "ETH" {
			continue
		}
		txs, err := service.GetTransactionsFor(wallet.GetAddress())
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txs...)
	}
	return transactions, nil
}

func (service *EthereumWebClient) GetTransactionsFor(address string) ([]common.Transaction, error) {
	ETHERSCAN_RATELIMITER.RespectRateLimit()
	url := fmt.Sprintf("%s?module=account&action=txlist&address=0x%s&startblock=0&endblock=99999999999999999999999&sort=asc&apikey=%s",
		service.endpoint, address, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EthereumWebClient.GetTransactionsFor] url: %s", url)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetTransactionsFor] HTTP Request Error: %s", err.Error())
	}
	response := EtherScanGetTransactionsResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetTransactionsFor] JSON Unmarshal error: %s", jsonErr.Error())
	}
	var transactions []common.Transaction
	for _, tx := range response.Result {
		var txType string
		hexAddress := fmt.Sprintf("0x%s", address)
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
		fiatTotal := finalAmount.Mul(candlestick.Close).Add(fee)
		transactions = append(transactions, &dto.TransactionDTO{
			Id:                   fmt.Sprintf("eth-%s", tx.Timestamp),
			Date:                 time.Unix(timestamp, 0),
			CurrencyPair:         &common.CurrencyPair{Base: "ETH", Quote: "ETH", LocalCurrency: localCurrency},
			Type:                 txType,
			Category:             common.TX_CATEGORY_TRANSFER,
			Network:              "Ethereum",
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
			FiatFee:              finalFee.Mul(candlestick.Close).StringFixed(2),
			FiatFeeCurrency:      "USD",
			Total:                finalAmount.Mul(candlestick.Close).StringFixed(2),
			TotalCurrency:        "USD",
			FiatTotal:            fiatTotal.StringFixed(2),
			FiatTotalCurrency:    "USD"})
	}
	return transactions, nil
}

func (service *EthereumWebClient) GetToken(walletAddress, contractAddress string) (common.EthereumToken, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	var symbol string
	url := fmt.Sprintf("%s?module=account&action=tokenbalance&contractaddress=0x%s&&address=0x%s&tag=latest&apikey=%s",
		service.endpoint, contractAddress, walletAddress, service.apiKeyToken)

	service.ctx.GetLogger().Debugf("[EthereumWebClient.GetToken] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetToken] Error: %s", err.Error())
		return nil, err
	}

	response := EtherScanResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetToken] %s", jsonErr.Error())
	}

	result, err := decimal.NewFromString(response.Result)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetToken] Result decimal conversation error: %s", err.Error())
	}

	tokens := service.userDAO.GetTokens(&entity.User{
		Id: service.ctx.GetUser().GetId()})

	for _, token := range tokens {
		if token.GetWalletAddress() == walletAddress && token.GetContractAddress() == contractAddress {
			symbol = token.GetSymbol()
			break
		}
	}

	balance := result.Div(decimal.NewFromFloat(100000000))
	decPriceUSD, err := decimal.NewFromString(service.marketcapService.GetMarket(symbol).PriceUSD)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetToken] Balance decimal conversation error: %s", err.Error())
	}
	value := balance.Mul(decPriceUSD)
	return &dto.EthereumTokenDTO{
		Symbol:          symbol,
		WalletAddress:   walletAddress,
		ContractAddress: contractAddress,
		Balance:         balance.Truncate(8),
		Value:           value.Truncate(2)}, nil
}

func (service *EthereumWebClient) GetTokenTransactions(contractAddress string) ([]common.Transaction, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=txlistinternal&address=0x%s&startblock=0&endblock=99999999999999999999999&sort=asc&apikey=%s",
		service.endpoint, contractAddress, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EthereumWebClient.GetTokenTransactions] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetTokenTransactions] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetTransactionsResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetTokenTransactions] JSON Unmarshall error: %s", jsonErr.Error())
	}

	priceUSD, err := decimal.NewFromString(service.marketcapService.GetMarket("ETH").PriceUSD)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EthereumWebClient.GetTokenTransactions] Error converting price to decimal: %s", err.Error())
	}

	var transactions []common.Transaction
	for _, tx := range response.Result {
		var txType string
		hexAddress := fmt.Sprintf("0x%s", contractAddress)
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
		transactions = append(transactions, &dto.TransactionDTO{
			Id:               tx.Timestamp,
			Date:             time.Unix(timestamp, 0),
			CurrencyPair:     &common.CurrencyPair{Base: "ETH", Quote: localCurrency, LocalCurrency: localCurrency},
			Type:             strings.Title(txType),
			Network:          "Ethereum",
			Quantity:         finalAmount.StringFixed(8),
			QuantityCurrency: "ETH",
			Price:            priceUSD.StringFixed(2),
			PriceCurrency:    "USD",
			Fee:              fee.Div(decimal.NewFromFloat(100000000)).StringFixed(8),
			FeeCurrency:      "ETH",
			Total:            finalAmount.Mul(priceUSD).StringFixed(2),
			TotalCurrency:    "USD"})
	}
	return transactions, nil
}
