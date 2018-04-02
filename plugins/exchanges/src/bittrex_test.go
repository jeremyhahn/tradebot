// +build integration

package main

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
)

func TestBittrex_GetPriceAt(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "Bittrex")
	assert.Nil(t, err)

	bittrex := CreateBittrex(ctx, userExchangeService).(*Bittrex)
	atDate, err := time.Parse("2006-01-02 15:04:05", "2017-12-21 04:32:49")

	candle, err := bittrex.GetPriceAt("BCC", atDate)

	assert.Nil(t, err)
	assert.Equal(t, atDate.Month(), candle.Date.Month())
	assert.Equal(t, atDate.Day(), candle.Date.Day())
	assert.Equal(t, atDate.Year(), candle.Date.Year())

	test.CleanupIntegrationTest()
}

func TestBittrex_GetCurrencies(t *testing.T) {
	bittrex := createBittrexService(t)
	currencies, err := bittrex.GetCurrencies()
	assert.Nil(t, err)
	assert.Equal(t, true, len(currencies) > 0)
	test.CleanupIntegrationTest()
}

func TestBittrex_GetPriceHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}
	cryptoExchange, err := userDAO.GetExchange(userEntity, "Bittrex")
	assert.Nil(t, err)
	bittrex := CreateBittrex(ctx, cryptoExchange)
	assert.NotNil(t, bittrex)
	currencyPair := &common.CurrencyPair{
		Base:          "USDT",
		Quote:         "BTC",
		LocalCurrency: "USD"}
	start := time.Now().AddDate(0, 0, -2)
	end := time.Now().AddDate(0, 0, -1)
	history, err := bittrex.GetPriceHistory(currencyPair, start, end, 100) // Hourly
	assert.Nil(t, err)
	lastIdx := len(history) - 1
	assert.Equal(t, true, len(history) > 0)
	assert.Equal(t, start.Month(), history[0].Date.Month())
	assert.Equal(t, start.Day(), history[0].Date.Day())
	assert.Equal(t, start.Year(), history[0].Date.Year())
	assert.Equal(t, end.Month(), history[lastIdx].Date.Month())
	assert.Equal(t, end.Day(), history[lastIdx].Date.Day())
	assert.Equal(t, end.Year(), history[lastIdx].Date.Year())
	test.CleanupIntegrationTest()
}

/* Only orders within last 30 days available via Bittrex API
func TestBittrex_GetOrderHistory(t *testing.T) {
	bittrex := createBittrexService(t)
	orders := bittrex.GetOrderHistory(&common.CurrencyPair{
		Base:          "BTC",
		Quote:         "ADA",
		LocalCurrency: "USD"})
	assert.Equal(t, true, len(orders) > 0)
	util.DUMP(orders)
	test.CleanupIntegrationTest()
}*/

func TestBittrex_GetDepositHistory(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	userExchangeService, err := userDAO.GetExchange(userEntity, "Bittrex")
	assert.Nil(t, err)

	bittrex := CreateBittrex(ctx, userExchangeService).(*Bittrex)
	deposits, err := bittrex.GetDepositHistory()

	for _, deposit := range deposits {
		util.DUMP(deposit)
	}

	assert.Nil(t, err)

	test.CleanupIntegrationTest()
}

func createBittrexService(t *testing.T) common.Exchange {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := service.CreatePluginService(ctx, "../../", pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	bittrex, err := exchangeService.GetExchange("Bittrex")
	assert.Nil(t, err)
	return bittrex
}
