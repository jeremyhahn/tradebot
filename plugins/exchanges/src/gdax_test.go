// +build integration

package main

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGDAX_GetExchange(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)

	exchange := gdax.GetSummary()
	assert.NotNil(t, exchange)
	assert.Equal(t, "GDAX", exchange.GetName())
	test.CleanupIntegrationTest()
}

func TestGDAX_GetSummary(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)

	exchange := gdax.GetSummary()
	assert.NotNil(t, exchange)
	assert.Equal(t, "GDAX", exchange.GetName())
	test.CleanupIntegrationTest()
}

func TestGDAX_GetBalances(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)

	balance, netWorth := gdax.GetBalances()
	assert.Equal(t, true, len(balance) > 0)
	assert.Equal(t, true, netWorth.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}

func TestGDAX_GetOrderHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)

	orders := gdax.GetOrderHistory(&common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"})
	assert.Equal(t, true, len(orders) > 0)

	for _, o := range orders {
		util.DUMP(o)
	}

	test.CleanupIntegrationTest()
}

func TestGDAX_GetPriceHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}
	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)
	gdax := CreateGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	history, err := gdax.GetPriceHistory(currencyPair, start, end, 900)
	assert.Nil(t, err)
	lastIdx := len(history) - 1
	assert.Equal(t, true, len(history) > 0)
	assert.Equal(t, start.Month(), history[0].Date.Month())
	assert.Equal(t, start.Day(), history[0].Date.Day())
	assert.Equal(t, start.Year(), history[0].Date.Year())
	assert.Equal(t, end.Month(), history[lastIdx].Date.Month())
	assert.Equal(t, end.Day(), history[lastIdx].Date.Day())
	assert.Equal(t, end.Year(), history[lastIdx].Date.Year())
	assert.Equal(t, false, end.Before(time.Now().Add(-900*time.Second)))
	test.CleanupIntegrationTest()
}

func TestGDAX_GetDespositHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange).(*GDAX)
	assert.NotNil(t, gdax)

	deposits, err := gdax.GetDepositHistory()
	assert.Nil(t, err)

	assert.Equal(t, true, len(deposits) > 0)
	for _, deposit := range deposits {
		assert.Equal(t, common.DEPOSIT_ORDER_TYPE, deposit.GetType())
	}

	test.CleanupIntegrationTest()
}

func TestGDAX_GetWithdrawHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange).(*GDAX)
	assert.NotNil(t, gdax)

	withdrawals, err := gdax.GetWithdrawalHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(withdrawals) > 0)
	for _, withdrawal := range withdrawals {
		assert.Equal(t, common.WITHDRAWAL_ORDER_TYPE, withdrawal.GetType())
		//ctx.GetLogger().Debug(withdrawal)
	}
	test.CleanupIntegrationTest()
}

func TestGDAX_GetCurrencies(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "GDAX")
	assert.Nil(t, err)

	gdax := CreateGDAX(ctx, cryptoExchange).(*GDAX)
	assert.NotNil(t, gdax)

	currencies, err := gdax.GetCurrencies()
	assert.Nil(t, err)

	assert.Equal(t, true, len(currencies) > 0)

	test.CleanupIntegrationTest()
}
