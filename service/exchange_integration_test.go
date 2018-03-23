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
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	exchanges, err := exchangeService.GetExchanges()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(exchanges))
	CleanupIntegrationTest()
}

func TestExchangeService_GetExchange(t *testing.T) {
	ctx := CreateIntegrationTestContext("../.env", "../")
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	gdax, err := exchangeService.GetExchange("GDAX")
	assert.Nil(t, err)
	assert.Equal(t, "GDAX", gdax.GetName())
	CleanupIntegrationTest()
}

func TestExchangeService_CreateExchange(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	gdax, err := exchangeService.CreateExchange("GDAX")
	assert.Nil(t, err)
	assert.Equal(t, "GDAX", gdax.GetName())
	CleanupIntegrationTest()
}
