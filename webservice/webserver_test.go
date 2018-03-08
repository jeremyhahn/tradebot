// +build broken_integration

package webservice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMarketCapService struct {
	*service.MarketCapService
	mock.Mock
}

type MockExchangeService struct {
	mock.Mock
}

type MockEthereumService struct {
	mock.Mock
}

type MockUserService struct {
	mock.Mock
}

type MockPortfolioService struct {
	mock.Mock
}

type MockOrderService struct {
	mock.Mock
}

func TestWebServer(t *testing.T) {
	ctx := test.NewUnitTestContext()

	mockEthereumService := new(MockEthereumService)
	marketcapService := service.NewMarketCapService(ctx.GetLogger())

	ws := NewWebServer(ctx, 8081, marketcapService, new(MockExchangeService),
		mockEthereumService, new(MockUserService), new(MockPortfolioService), new(MockOrderService))

	go ws.Start()
	go ws.Run()

	creds := &UserCredentials{
		Username: "unittest",
		Password: "unittest"}
	jsonCreds, err := json.Marshal(creds)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8081/api/v1/login", bytes.NewBuffer(jsonCreds))
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "\"token\":")
	assert.Equal(t, 200, res.StatusCode)

	ws.Stop()
}

func (ethereum *MockEthereumService) Login(username, password string) (common.UserContext, error) {
	return &dto.UserDTO{
		Id:            1,
		Username:      "testing",
		LocalCurrency: "USD",
		Etherbase:     "0xabc123"}, nil
}

func (ethereum *MockEthereumService) Register(username, password string) error {
	return nil
}

func (exchange *MockExchangeService) CreateExchange(user common.UserContext, name string) (common.Exchange, error) {
	return nil, nil
}

func (exchange *MockExchangeService) GetCurrencyPairs(user common.UserContext, name string) ([]common.CurrencyPair, error) {
	return []common.CurrencyPair{}, nil
}

func (exchange *MockExchangeService) GetExchange(user common.UserContext, name string) common.Exchange {
	return nil
}

func (exchange *MockExchangeService) GetExchanges(user common.UserContext) []common.Exchange {
	return nil
}

func (exchange *MockExchangeService) GetDisplayNames(user common.UserContext) []string {
	return []string{}
}

func (mock *MockUserService) CreateUser(user common.UserContext) {
}

func (mock *MockUserService) GetExchange(user common.UserContext, name string, currencyPair *common.CurrencyPair) (common.Exchange, error) {
	return nil, nil
}

func (mock *MockUserService) GetExchanges(user common.UserContext, currencyPair *common.CurrencyPair) []common.UserCryptoExchange {
	var exchanges []common.UserCryptoExchange
	return exchanges
}

func (user *MockUserService) GetCurrentUser() (common.UserContext, error) {
	return nil, nil
}

func (user *MockUserService) GetUserById(uint) (common.UserContext, error) {
	return nil, nil
}

func (user *MockUserService) GetUserByName(string) (common.UserContext, error) {
	return nil, nil
}

func (user *MockUserService) GetWallet(common.UserContext, string) common.UserCryptoWallet {
	return nil
}

func (user *MockUserService) GetWallets(common.UserContext) []common.UserCryptoWallet {
	return nil
}

func (user *MockUserService) GetTokens(userContext common.UserContext, wallet string) ([]common.EthereumToken, error) {
	return []common.EthereumToken{}, nil
}

func (user *MockUserService) GetAllTokens(common.UserContext) ([]common.EthereumToken, error) {
	return []common.EthereumToken{}, nil
}

func (user *MockUserService) GetConfiguredExchanges(common.UserContext) []common.UserCryptoExchange {
	return []common.UserCryptoExchange{}
}

func (user *MockUserService) GetExchangeSummary(common.UserContext, *common.CurrencyPair) ([]common.CryptoExchangeSummary, error) {
	return []common.CryptoExchangeSummary{}, nil
}

func (user *MockPortfolioService) Build(common.UserContext, *common.CurrencyPair) (common.Portfolio, error) {
	return nil, nil
}

func (user *MockPortfolioService) Queue(common.UserContext) (<-chan common.Portfolio, error) {
	return make(chan common.Portfolio), nil
}

func (user *MockPortfolioService) IsStreaming(common.UserContext) bool {
	return false
}

func (user *MockPortfolioService) Stop(common.UserContext) {
}

func (user *MockPortfolioService) Stream(common.UserContext, *common.CurrencyPair) (<-chan common.Portfolio, error) {
	return make(chan common.Portfolio), nil
}

func (user *MockOrderService) GetMapper() mapper.OrderMapper {
	return nil
}

func (user *MockOrderService) GetOrderHistory() []common.Order {
	return []common.Order{}
}

func (user *MockOrderService) ImportCSV(file, exchangeName string) ([]common.Order, error) {
	return []common.Order{}, nil
}
