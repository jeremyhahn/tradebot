// +build integration

package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCoinbase_GetBalance(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	coins, sum := cb.GetBalances()

	assert.Equal(t, true, len(coins) > 0)
	assert.Equal(t, true, sum.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetOrderHistory(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	orders := cb.GetOrderHistory(&common.CurrencyPair{
		Base:          "ETH",
		Quote:         "USD",
		LocalCurrency: "USD"})

	assert.Nil(t, err)
	assert.Equal(t, true, len(orders) > 0)

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetDeposits(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	deposits, err := cb.GetDepositHistory()
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, true, len(deposits) > 0)

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetWithdrawls(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	withdrawls, err := cb.GetWithdrawalHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(withdrawls) > 0)

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetCurrencies(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	currencies, err := cb.GetCurrencies()
	assert.Nil(t, err)
	assert.Equal(t, true, len(currencies) > 0)

	test.CleanupIntegrationTest()
}
