// +build integration

package main

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEthWallet_GetWallet(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := service.NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	ethwallet := CreateEthWallet(&common.WalletParams{
		Context:          ctx,
		Address:          os.Getenv("ETH_ADDRESS"),
		MarketCapService: service.NewMarketCapService(ctx),
		FiatPriceService: fiatPriceService,
		WalletSecret:     "YourApiKeyToken"})

	wallet, err := ethwallet.GetWallet()
	assert.Nil(t, err)

	zero := decimal.NewFromFloat(0)

	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, true, wallet.GetBalance().GreaterThan(zero))
	assert.Equal(t, true, wallet.GetValue().GreaterThan(zero))

	test.CleanupIntegrationTest()
}

func TestEthWalletService_GetTransactions(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := service.NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)
	ethereumService, err := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	assert.Nil(t, err)
	userService := service.NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService,
		ethereumService, exchangeService, service.NewWalletService(ctx, pluginService))

	userService.CreateWallet(&dto.UserCryptoWalletDTO{
		Address:  os.Getenv("ETH_ADDRESS"),
		Currency: "ETH"})

	userService.CreateToken(&dto.EthereumTokenDTO{
		Symbol:          os.Getenv("TOKEN_SYMBOL"),
		ContractAddress: os.Getenv("TOKEN_ADDRESS"),
		WalletAddress:   os.Getenv("ETH_ADDRESS")})

	ethwallet := CreateEthWallet(&common.WalletParams{
		Context:          ctx,
		Address:          os.Getenv("ETH_ADDRESS"),
		MarketCapService: marketcapService,
		FiatPriceService: fiatPriceService,
		WalletSecret:     "YourApiKeyToken"})

	transactions, err := ethwallet.GetTransactions()

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	var totalDeposit, totalWithdrawl, totalFees decimal.Decimal

	for _, tx := range transactions {
		qty, err := decimal.NewFromString(tx.GetQuantity())
		if err != nil {
			panic(err)
		}
		if tx.GetType() == common.DEPOSIT_ORDER_TYPE {
			totalDeposit = totalDeposit.Add(qty)
		} else if tx.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			fee, err := decimal.NewFromString(tx.GetFee())
			if err != nil {
				panic(err)
			}
			totalWithdrawl = totalWithdrawl.Add(qty)
			totalFees = totalFees.Add(fee)
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}

	assert.Equal(t, true, totalDeposit.GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, totalWithdrawl.GreaterThan(decimal.NewFromFloat(0)))

	txSum := totalDeposit.Sub(totalWithdrawl).Sub(totalFees)
	wallet, err := ethwallet.GetWallet()
	assert.Nil(t, err)

	actual := wallet.GetBalance().StringFixed(8)
	expected := txSum.StringFixed(8)
	assert.Equal(t, actual, expected)

	test.CleanupIntegrationTest()
}
