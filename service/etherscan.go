package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/util"
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

type EtherscanServiceImpl struct {
	ctx                 common.Context
	userDAO             dao.UserDAO
	userMapper          mapper.UserMapper
	marketcapService    MarketCapService
	pricehistoryService common.PriceHistoryService
	apiKeyToken         string
	endpoint            string
	rateLimiter         *common.RateLimiter
	EthereumService
}

var ETHERSCAN_RATELIMITER = common.NewRateLimiter(5, 1)

func NewEtherscanService(ctx common.Context, userDAO dao.UserDAO, userMapper mapper.UserMapper,
	marketcapService MarketCapService, priceHistoryService common.PriceHistoryService) (EthereumService, error) {
	return &EtherscanServiceImpl{
		ctx:                 ctx,
		userDAO:             userDAO,
		userMapper:          userMapper,
		marketcapService:    marketcapService,
		pricehistoryService: priceHistoryService,
		apiKeyToken:         ctx.GetKeystore(),
		endpoint:            "https://api.etherscan.io/api"}, nil
}

func (service *EtherscanServiceImpl) Login(username, password string) (common.UserContext, error) {
	return &dto.UserContextDTO{
		Id:            1,
		Username:      username,
		LocalCurrency: "USD",
		Etherbase:     username,
		Keystore:      password}, nil
}

func (service *EtherscanServiceImpl) Register(username, password string) error {
	return nil
}

func (service *EtherscanServiceImpl) GetAccounts() ([]common.UserContext, error) {
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

func (service *EtherscanServiceImpl) GetWallet(address string) (common.UserCryptoWallet, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=balance&address=0x%s&tag=latest&apikey=%s",
		service.endpoint, address, service.apiKeyToken)

	service.ctx.GetLogger().Debugf("[EtherscanService.GetBalance] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetBalance] Error: %s", err.Error())
		return nil, err
	}

	response := EtherScanResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetBalance] %s", jsonErr.Error())
		return nil, jsonErr
	}

	result, err := strconv.ParseFloat(response.Result, 64)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetBalance] Float conversation error: %s", err.Error())
		return nil, err
	}

	balance := result / 1000000000000000000

	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  balance,
		Currency: "ETH",
		Value:    service.getLastPrice() * balance}, nil
}

func (service *EtherscanServiceImpl) GetTransactions() ([]common.Transaction, error) {
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

func (service *EtherscanServiceImpl) GetTransactionsFor(address string) ([]common.Transaction, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=txlist&address=0x%s&startblock=0&endblock=99999999999999999999999&sort=asc&apikey=%s",
		service.endpoint, address, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EtherscanService.GetTransactionsFor] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetTransactionsFor] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetTransactionsResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetTransactionsFor] JSON Unmarshall error: %s", jsonErr.Error())
	}

	var transactions []common.Transaction
	for _, tx := range response.Result {
		var txType string
		hexAddress := fmt.Sprintf("0x%s", address)
		if tx.From == hexAddress {
			txType = "Withdrawl"
		} else if tx.To == hexAddress {
			txType = "Deposit"
		}
		timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(tx.Value, 64)
		if err != nil {
			return nil, err
		}
		fee, err := strconv.ParseFloat(tx.GasUsed, 64)
		if err != nil {
			return nil, err
		}
		localCurrency := service.ctx.GetUser().GetLocalCurrency()
		finalAmount := amount / 1000000000000000000

		closePrice := service.pricehistoryService.GetClosePriceOn("ETH", time.Unix(timestamp, 0))

		transactions = append(transactions, &dto.TransactionDTO{
			Date:               time.Unix(timestamp, 0),
			CurrencyPair:       &common.CurrencyPair{Base: "ETH", Quote: localCurrency, LocalCurrency: localCurrency},
			Type:               txType,
			Source:             "Ethereum",
			Amount:             finalAmount,
			AmountCurrency:     "ETH",
			Fee:                fee / 100000000,
			FeeCurrency:        "ETH",
			Total:              finalAmount * closePrice,
			TotalCurrency:      "USD",
			HistoricalPrice:    closePrice,
			HistoricalCurrency: "USD"})
	}
	return transactions, nil
}

func (service *EtherscanServiceImpl) GetToken(walletAddress, contractAddress string) (common.EthereumToken, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	var symbol string
	url := fmt.Sprintf("%s?module=account&action=tokenbalance&contractaddress=0x%s&&address=0x%s&tag=latest&apikey=%s",
		service.endpoint, contractAddress, walletAddress, service.apiKeyToken)

	service.ctx.GetLogger().Debugf("[EtherscanService.GetToken] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetToken] Error: %s", err.Error())
		return nil, err
	}

	response := EtherScanResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetToken] %s", jsonErr.Error())
	}

	result, err := strconv.ParseFloat(response.Result, 64)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetToken] Float conversation error: %s", err.Error())
	}

	tokens := service.userDAO.GetTokens(&entity.User{
		Id: service.ctx.GetUser().GetId()})

	for _, token := range tokens {
		if token.GetWalletAddress() == walletAddress && token.GetContractAddress() == contractAddress {
			symbol = token.GetSymbol()
			break
		}
	}

	balance := result / 100000000
	fPriceUSD, _ := strconv.ParseFloat(service.marketcapService.GetMarket(symbol).PriceUSD, 64)
	value := balance * fPriceUSD

	return &dto.EthereumTokenDTO{
		Symbol:          symbol,
		WalletAddress:   walletAddress,
		ContractAddress: contractAddress,
		Balance:         balance,
		Value:           value}, nil
}

func (service *EtherscanServiceImpl) GetTokenTransactions(contractAddress string) ([]common.Transaction, error) {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=account&action=txlistinternal&address=0x%s&startblock=0&endblock=99999999999999999999999&sort=asc&apikey=%s",
		service.endpoint, contractAddress, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EtherscanService.GetTokenTransactions] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetTokenTransactions] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetTransactionsResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetTokenTransactions] JSON Unmarshall error: %s", jsonErr.Error())
	}

	priceUSD, err := strconv.ParseFloat(service.marketcapService.GetMarket("ETH").PriceUSD, 64)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.GetTokenTransactions] Error converting price to float: %s", err.Error())
	}

	var transactions []common.Transaction
	for _, tx := range response.Result {
		var txType string
		hexAddress := fmt.Sprintf("0x%s", contractAddress)
		if tx.From == hexAddress {
			txType = "Withdrawl"
		} else if tx.To == hexAddress {
			txType = "Deposit"
		}
		timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(tx.Value, 64)
		if err != nil {
			return nil, err
		}
		fee, err := strconv.ParseFloat(tx.GasUsed, 64)
		if err != nil {
			return nil, err
		}
		localCurrency := service.ctx.GetUser().GetLocalCurrency()
		finalAmount := amount / 1000000000000000000
		transactions = append(transactions, &dto.TransactionDTO{
			Date:           time.Unix(timestamp, 0),
			CurrencyPair:   &common.CurrencyPair{Base: "ETH", Quote: localCurrency, LocalCurrency: localCurrency},
			Type:           txType,
			Source:         "Ethereum",
			Amount:         finalAmount,
			AmountCurrency: "ETH",
			Fee:            fee / 100000000,
			FeeCurrency:    "ETH",
			Total:          finalAmount * priceUSD,
			TotalCurrency:  "USD"})
	}

	return transactions, nil
}

func (service *EtherscanServiceImpl) getLastPrice() float64 {

	ETHERSCAN_RATELIMITER.RespectRateLimit()

	url := fmt.Sprintf("%s?module=stats&action=ethprice&apikey=%s", service.endpoint, service.apiKeyToken)
	service.ctx.GetLogger().Debugf("[EtherscanService.getLastPrice] url: %s", url)

	_, body, err := util.HttpRequest(url)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.getLastPrice] HTTP Request Error: %s", err.Error())
	}

	response := EtherScanGetLastPriceResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.getLastPrice] JSON Unmarshall Error: %s", jsonErr.Error())
	}

	price, err := strconv.ParseFloat(response.Result.USD, 64)
	if err != nil {
		service.ctx.GetLogger().Errorf("[EtherscanService.getLastPrice] Float Conversation Error: %s", err.Error())
	}

	return price
}
