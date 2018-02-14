// +build integration_webservice

package webservice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
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

func TestWebServer(t *testing.T) {
	ctx := test.NewUnitTestContext()

	mockEthereumService := new(MockEthereumService)
	marketcapService := service.NewMarketCapService(ctx.Logger)

	rsaKeyPair, err := common.CreateRsaKeyPair(ctx, "../test/keys")
	assert.Nil(t, err)

	jwt := CreateJsonWebToken(ctx, mockEthereumService, NewJsonWriter(), 10, rsaKeyPair)
	assert.Nil(t, err)

	assert.Equal(t, "../test/keys", jwt.rsaKeyPair.Directory)

	ws := NewWebServer(ctx, 8081, marketcapService, new(MockExchangeService),
		mockEthereumService, new(MockUserService), new(MockPortfolioService), jwt)

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

	claims := jwt.GetClaims()
	assert.NotNil(t, claims)
	assert.NotNil(t, "1", claims["user_id"])
	assert.NotNil(t, "testing", claims["username"])
	assert.NotNil(t, "USD", claims["local_currency"])
	assert.NotNil(t, "0xabc123", claims["etherbase"])

	ws.Stop()
}

func (ethereum *MockEthereumService) Login(username, password string) (common.User, error) {
	return &dto.UserDTO{
		Id:            1,
		Username:      "testing",
		LocalCurrency: "USD",
		Etherbase:     "0xabc123"}, nil
}

func (ethereum *MockEthereumService) Register(username, password string) error {
	return nil
}

func (exchange *MockExchangeService) CreateExchange(user common.User, name string) common.Exchange {
	return nil
}

func (exchange *MockExchangeService) GetExchange(user common.User, name string) common.Exchange {
	return nil
}

func (exchange *MockExchangeService) GetExchanges(user common.User) []common.Exchange {
	return nil
}

func (mock *MockUserService) CreateUser(user common.User) {
}

func (mock *MockUserService) GetExchange(user common.User, name string, currencyPair *common.CurrencyPair) common.Exchange {
	return nil
}

func (mock *MockUserService) GetExchanges(user common.User, currencyPair *common.CurrencyPair) []common.CryptoExchange {
	var exchanges []common.CryptoExchange
	return exchanges
}

func (user *MockUserService) GetCurrentUser() (common.User, error) {
	return nil, nil
}

func (user *MockUserService) GetUserById(uint) (common.User, error) {
	return nil, nil
}

func (user *MockUserService) GetUserByName(string) (common.User, error) {
	return nil, nil
}

func (user *MockUserService) GetWallet(common.User, string) common.CryptoWallet {
	return nil
}

func (user *MockUserService) GetWallets(common.User) []common.CryptoWallet {
	return nil
}

func (user *MockPortfolioService) Build(common.User, *common.CurrencyPair) common.Portfolio {
	return nil
}

func (user *MockPortfolioService) Queue(common.User) <-chan common.Portfolio {
	return make(chan common.Portfolio)
}

func (user *MockPortfolioService) IsStreaming(common.User) bool {
	return false
}

func (user *MockPortfolioService) Stop(common.User) {
}

func (user *MockPortfolioService) Stream(common.User, *common.CurrencyPair) <-chan common.Portfolio {
	return make(chan common.Portfolio)
}
