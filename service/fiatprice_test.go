// +build integration

package service

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestFiatPrice(t *testing.T) {
	ctx := NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	atDate := time.Now().Add(-5 * time.Hour)
	candlestick, err := fiatPriceService.GetPriceAt("BTC", atDate)
	assert.Nil(t, err)
	assert.Equal(t, true, candlestick.Close > 0)

	CleanupIntegrationTest()
}
