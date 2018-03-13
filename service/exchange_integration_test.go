// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestExchangeService_GetExchanges(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	priceHistoryService := NewPriceHistoryService(ctx)
	exchangeService := NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper, priceHistoryService)
	exchanges := exchangeService.GetExchanges()
	assert.Equal(t, 3, len(exchanges))
	CleanupIntegrationTest()
}

func TestExchangeService_GetExchange(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	priceHistoryService := NewPriceHistoryService(ctx)
	exchangeService := NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper, priceHistoryService)
	gdax := exchangeService.GetExchange("gdax")
	assert.Equal(t, "gdax", gdax.GetName())
	CleanupIntegrationTest()
}

func TestExchangeService_CreateExchange(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	priceHistoryService := NewPriceHistoryService(ctx)
	exchangeService := NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper, priceHistoryService)
	gdax, err := exchangeService.CreateExchange("gdax")
	assert.Nil(t, err)
	assert.Equal(t, "gdax", gdax.GetName())
	CleanupIntegrationTest()
}
