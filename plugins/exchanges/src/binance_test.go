// +build integration

package main

import (
	"encoding/json"
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestBinance_GetDepositHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "binance")
	assert.Nil(t, err)

	binance := NewBinance(ctx, userExchangeService).(*Binance)
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

	userExchangeService, err := userDAO.GetExchange(userEntity, "binance")
	assert.Nil(t, err)

	binance := NewBinance(ctx, userExchangeService).(*Binance)
	Withdraws, err := binance.GetWithdrawalHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(Withdraws) > 0)

	test.CleanupIntegrationTest()
}
