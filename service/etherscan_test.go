// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
)

func TestEtherScan_GetWallet(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(),
		NewMarketCapService(ctx), NewPriceHistoryService(ctx))
	assert.Nil(t, err)

	wallet, err := service.GetWallet(os.Getenv("ETH_ADDRESS"))
	assert.Nil(t, err)

	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, true, wallet.GetBalance() > 0)
	assert.Equal(t, true, wallet.GetValue() > 0)

	test.CleanupIntegrationTest()
}

func TestEtherScanService_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	walletAddress := os.Getenv("ETH_ADDRESS")

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := NewMarketCapService(ctx)
	priceHistoryService := NewPriceHistoryService(ctx)
	ethereumService, _ := NewEthereumService(ctx, userDAO, userMapper, marketcapService)
	userService := NewUserService(ctx, userDAO, pluginDAO, marketcapService,
		ethereumService, userMapper, userExchangeMapper, priceHistoryService)

	userService.CreateWallet(&dto.UserCryptoWalletDTO{
		Address:  os.Getenv("ETH_ADDRESS"),
		Currency: "ETH"})

	userService.CreateToken(&dto.EthereumTokenDTO{
		Symbol:          os.Getenv("TOKEN_SYMBOL"),
		ContractAddress: os.Getenv("TOKEN_ADDRESS"),
		WalletAddress:   os.Getenv("ETH_ADDRESS")})

	etherscanService, err := NewEtherscanService(ctx, userDAO, userMapper, marketcapService,
		NewPriceHistoryService(ctx))
	assert.Nil(t, err)

	transactions, err := etherscanService.GetTransactions()

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	totalDeposit := 0.0
	totalWithdrawl := 0.0
	totalFees := 0.0

	for _, tx := range transactions {
		if tx.GetType() == "Deposit" {
			totalDeposit += tx.GetAmount()
		} else if tx.GetType() == "Withdrawl" {
			totalWithdrawl += tx.GetAmount()
			totalFees += tx.GetFee()
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit > 0)
	assert.Equal(t, true, totalWithdrawl > 0)

	txSum := totalDeposit - totalWithdrawl - totalFees

	wallet, err := etherscanService.GetWallet(walletAddress)
	assert.Nil(t, err)

	actual := util.TruncateFloat(wallet.GetBalance(), 8)
	expected := util.TruncateFloat(txSum, 8)
	assert.Equal(t, actual, expected)

	test.CleanupIntegrationTest()
}

func TestEtherScan_GetToken(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(),
		NewMarketCapService(ctx), NewPriceHistoryService(ctx))
	assert.Nil(t, err)

	token, err := service.GetToken(os.Getenv("ETH_ADDRESS"), os.Getenv("TOKEN_ADDRESS"))
	assert.Nil(t, err)

	assert.Equal(t, true, token.GetBalance() > 0)
	//fmt.Printf("balance: %.8f\n", balance)

	test.CleanupIntegrationTest()
}

func TestEtherScanService_GetTokenTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(),
		NewMarketCapService(ctx), NewPriceHistoryService(ctx))
	assert.Nil(t, err)

	transactions, err := service.GetTokenTransactions(os.Getenv("TOKEN_INTERNAL_ADDRESS"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	test.CleanupIntegrationTest()
}
