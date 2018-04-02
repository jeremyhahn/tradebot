// +build integration

package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestBinance_GetPriceAt(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "Binance")
	assert.Nil(t, err)

	binance := CreateBinance(ctx, userExchangeService).(*Binance)
	atDate := time.Now().Add(-24 * time.Hour)
	candle, err := binance.GetPriceAt("BTC", atDate)
	assert.Nil(t, err)
	assert.Equal(t, atDate.Month(), candle.Date.Month())
	assert.Equal(t, atDate.Day(), candle.Date.Day())
	assert.Equal(t, atDate.Year(), candle.Date.Year())

	test.CleanupIntegrationTest()
}

func TestBinance_GetDepositHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "Binance")
	assert.Nil(t, err)

	binance := CreateBinance(ctx, userExchangeService).(*Binance)
	deposits, err := binance.GetDepositHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(deposits) > 0)

	jsonData, _ := json.Marshal(deposits)
	ctx.GetLogger().Debugf("[Binance.GetDepositHistory] orders: %s", string(jsonData))

	test.CleanupIntegrationTest()
}

func TestBinance_GetWithdrawHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "Binance")
	assert.Nil(t, err)

	binance := CreateBinance(ctx, userExchangeService).(*Binance)
	Withdraws, err := binance.GetWithdrawalHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(Withdraws) > 0)

	test.CleanupIntegrationTest()
}
