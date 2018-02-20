// +build broken_integration

package webservic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMarketCapService_OrderHistory struct {
	*service.MarketCapService
	mock.Mock
}

type MockExchange_OrderHistory struct {
	mock.Mock
}

type MockEthereum_OrderHistory struct {
	mock.Mock
}

type MockUser_OrderHistory struct {
	mock.Mock
}

type MockPortfolio_OrderHistory struct {
	mock.Mock
}

type MockUserService_OrderHistory struct {
	mock.Mock
	service.UserService
}

func TestOrderHistory(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	orderDAO := dao.NewOrderDAO(ctx)

	exchangeService := new(MockExchange_OrderHistory)
	userService := new(MockUserService_OrderHistory)
	orderMapper := mapper.NewOrderMapper(ctx)

	mockEthereumService := new(MockEthereum_OrderHistory)
	marketcapService := service.NewMarketCapService(ctx.Logger)
	orderService := service.NewOrderService(ctx, orderDAO, orderMapper, exchangeService, userService)
	//priceHistoryService := service.NewPriceHistoryService(ctx)

	rsaKeyPair, err := common.CreateRsaKeyPair(ctx, "../test/keys")
	jwt := CreateJsonWebToken(ctx, mockEthereumService, NewJsonWriter(), 10, rsaKeyPair)
	ws := NewWebServer(ctx, 8081, marketcapService, exchangeService,
		mockEthereumService, new(MockUser_OrderHistory), new(MockPortfolio_OrderHistory),
		orderService, jwt)

	go ws.Start()
	go ws.Run()

	creds := &UserCredentials{
		Username: "unittest",
		Password: "unittest"}
	jsonCreds, err := json.Marshal(creds)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://localhost:8081/api/v1/login", bytes.NewBuffer(jsonCreds))
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	util.DUMP(res)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	jwtResponse := string(bodyBytes)
	assert.Nil(t, err)
	assert.Contains(t, jwtResponse, "\"token\":")
	assert.Equal(t, 200, res.StatusCode)
	token := JsonWebTokenDTO{}
	err = json.Unmarshal(bodyBytes, &token)

	req, _ = http.NewRequest("GET", "https://localhost:8081/api/v1/orderhistory", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Value)
	res, _ = client.Do(req)
	bodyBytes, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(bodyBytes), "\"success\":true")

	ws.Stop()
	test.CleanupIntegrationTest()
}

func (ethereum *MockEthereum_OrderHistory) Login(username, password string) (common.User, error) {
	return &dto.UserDTO{
		Id:            1,
		Username:      "test",
		LocalCurrency: "USD",
		Etherbase:     "0xabc123",
		Keystore:      "/tmp"}, nil
}

func (ethereum *MockEthereum_OrderHistory) Register(username, password string) error {
	return nil
}

func (exchange *MockExchange_OrderHistory) CreateExchange(user common.User, name string) common.Exchange {
	return nil
}

func (exchange *MockExchange_OrderHistory) GetExchange(user common.User, name string) common.Exchange {
	return nil
}

func (exchange *MockExchange_OrderHistory) GetExchanges(user common.User) []common.Exchange {
	return nil
}

func (exchange *MockExchange_OrderHistory) GetCurrencyPairs(user common.User, pairs string) ([]common.CurrencyPair, error) {
	return []common.CurrencyPair{
		common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"}}, nil
}

func (exchange *MockExchange_OrderHistory) GetDisplayNames(user common.User) []string {
	return []string{"Exchange 1", "Exchange 2", "Exchange 3"}
}

func (mock *MockUser_OrderHistory) CreateUser(user common.User) {
}

func (mock *MockUser_OrderHistory) GetExchange(user common.User, name string, currencyPair *common.CurrencyPair) common.Exchange {
	return nil
}

func (mock *MockUser_OrderHistory) GetExchanges(user common.User, currencyPair *common.CurrencyPair) []common.CryptoExchange {
	var exchanges []common.CryptoExchange
	return exchanges
}

func (user *MockUser_OrderHistory) GetCurrentUser() (common.User, error) {
	return nil, nil
}

func (user *MockUser_OrderHistory) GetUserById(uint) (common.User, error) {
	return nil, nil
}

func (user *MockUser_OrderHistory) GetUserByName(string) (common.User, error) {
	return nil, nil
}

func (user *MockUser_OrderHistory) GetWallet(common.User, string) common.CryptoWallet {
	return nil
}

func (user *MockUser_OrderHistory) GetWallets(common.User) []common.CryptoWallet {
	return nil
}

func (user *MockPortfolio_OrderHistory) Build(common.User, *common.CurrencyPair) common.Portfolio {
	return nil
}

func (user *MockPortfolio_OrderHistory) Queue(common.User) <-chan common.Portfolio {
	return make(chan common.Portfolio)
}

func (user *MockPortfolio_OrderHistory) IsStreaming(common.User) bool {
	return false
}

func (user *MockPortfolio_OrderHistory) Stop(common.User) {
}

func (user *MockPortfolio_OrderHistory) Stream(common.User, *common.CurrencyPair) <-chan common.Portfolio {
	return make(chan common.Portfolio)
}
