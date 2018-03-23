// +build broken_integration

package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
)

/*
func TestBittrex_GetCurrencies(t *testing.T) {
	bittrex := createBittrexService(t)
	currencies, err := bittrex.GetCurrencies()
	assert.Nil(t, err)
	assert.Equal(t, true, len(currencies) > 0)
	test.CleanupIntegrationTest()
}*/

func TestBittrex_GetOrderHistory(t *testing.T) {
	bittrex := createBittrexService(t)
	orders := bittrex.GetOrderHistory(&common.CurrencyPair{
		Base:          "BTC",
		Quote:         "ADA",
		LocalCurrency: "USD"})
	assert.Equal(t, true, len(orders) > 0)

	util.DUMP(orders)

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
