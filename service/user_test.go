// +build integration

package service

import (
	"os"
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	service := createTestUserService()

	userById, err := service.GetUserById(1)
	assert.Nil(t, err)
	assert.Equal(t, uint(1), userById.GetId())
	assert.Equal(t, "test", userById.GetUsername())

	userByName, err := service.GetUserByName("test")
	assert.Nil(t, err)
	assert.Equal(t, uint(1), userByName.GetId())
	assert.Equal(t, "test", userByName.GetUsername())

	test.CleanupIntegrationTest()
}

func TestUserService_CreateGetWallet(t *testing.T) {
	service := createTestUserService()

	wallet := &dto.UserCryptoWalletDTO{
		Address:  os.Getenv("ETH_ADDRESS"),
		Currency: "ETH"}

	err := service.CreateWallet(wallet)
	assert.Nil(t, err)

	ethWallet := service.GetWallet("ETH")
	assert.NotNil(t, ethWallet)
	assert.Equal(t, wallet.GetAddress(), ethWallet.GetAddress())
	assert.Equal(t, wallet.GetCurrency(), ethWallet.GetCurrency())
	assert.Equal(t, true, ethWallet.GetBalance() > 0)
	assert.Equal(t, true, ethWallet.GetValue() > 0)

	test.CleanupIntegrationTest()
}

func TestUserService_GetWallets(t *testing.T) {
	service := createTestUserService()
	wallets := service.GetWallets()
	assert.Equal(t, os.Getenv("BTC_ADDRESS"), wallets[0].GetAddress())
	assert.Equal(t, os.Getenv("XRP_ADDRESS"), wallets[1].GetAddress())
	test.CleanupIntegrationTest()
}

func TestUserService_GetExchanges(t *testing.T) {
	service := createTestUserService()

	exchanges := service.GetConfiguredExchanges()

	assert.Equal(t, "gdax", exchanges[0].GetName())
	assert.Equal(t, "bittrex", exchanges[1].GetName())
	assert.Equal(t, "binance", exchanges[2].GetName())

	test.CleanupIntegrationTest()
}

func TestUserService_GetTokens(t *testing.T) {
	service := createTestUserService()

	token := &dto.EthereumTokenDTO{
		Symbol:          "GLA",
		ContractAddress: os.Getenv("TOKEN_ADDRESS"),
		WalletAddress:   os.Getenv("ETH_ADDRESS")}

	err := service.CreateToken(token)
	assert.Nil(t, err)

	tokens, err := service.GetTokens(os.Getenv("ETH_ADDRESS"))

	assert.Nil(t, err)
	assert.Equal(t, 1, len(tokens))

	assert.Equal(t, token.GetContractAddress(), tokens[0].GetContractAddress())
	assert.Equal(t, token.GetWalletAddress(), tokens[0].GetWalletAddress())
	assert.Equal(t, true, tokens[0].GetBalance() > 0)

	test.CleanupIntegrationTest()
}

func createTestUserService() UserService {
	ctx := test.NewIntegrationTestContext()
	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := NewMarketCapService(ctx)
	ethereumService, _ := NewEthereumService(ctx, userDAO, userMapper, marketcapService)
	return NewUserService(ctx, userDAO, pluginDAO, marketcapService,
		ethereumService, userMapper, userExchangeMapper)
}
