// +build integration

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBtcWallet_GetPrice(t *testing.T) {
	ctx := test.NewUnitTestContext()
	wallet := CreateBtcWallet(&common.WalletParams{Context: ctx})
	price := wallet.GetPrice()
	assert.Equal(t, true, price.GreaterThan(decimal.NewFromFloat(0)))
}

func TestBtcWallet_GetBalance(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	wallet := CreateBtcWallet(&common.WalletParams{
		Context: ctx,
		Address: os.Getenv("BTC_ADDRESS")})
	userWallet, err := wallet.GetWallet()
	assert.Nil(t, err)
	assert.Equal(t, true, len(userWallet.GetAddress()) > 0)
	assert.Equal(t, true, len(userWallet.GetCurrency()) > 0)
	assert.Equal(t, true, userWallet.GetBalance().GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, userWallet.GetValue().GreaterThan(decimal.NewFromFloat(0)))
	test.CleanupIntegrationTest()
}

func TestBtcWallet_GetTransactions(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := service.NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	wallet := CreateBtcWallet(&common.WalletParams{
		Context:          ctx,
		Address:          os.Getenv("BTC_ADDRESS"),
		FiatPriceService: fiatPriceService})

	transactions, err := wallet.GetTransactions()
	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	totalWithdrawl := decimal.NewFromFloat(0.0)
	totalDeposit := decimal.NewFromFloat(0.0)

	for _, tx := range transactions {
		qty, _ := decimal.NewFromString(tx.GetQuantity())
		if tx.GetType() == common.DEPOSIT_ORDER_TYPE {
			totalDeposit = totalDeposit.Add(qty)
		} else if tx.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			totalWithdrawl = totalWithdrawl.Add(qty)
		}
		fmt.Printf("%+v\n", tx)
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit.GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, totalWithdrawl.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}
