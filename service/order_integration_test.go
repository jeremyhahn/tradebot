// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestOrderService_GetOrderHistory(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	exchangeDAO := dao.NewExchangeDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	orderDAO := dao.NewOrderDAO(ctx)

	userMapper := mapper.NewUserMapper()
	orderMapper := mapper.NewOrderMapper(ctx)
	exchangeMapper := mapper.NewExchangeMapper()

	marketcapService := NewMarketCapService(ctx.Logger)
	ethereumService, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO, userMapper)
	assert.Nil(t, err)

	exchangeService := NewExchangeService(ctx, exchangeDAO, userDAO, userMapper, exchangeMapper)
	userService := NewUserService(ctx, userDAO, marketcapService, ethereumService, userMapper, exchangeMapper)
	orderService := NewOrderService(ctx, orderDAO, orderMapper, exchangeService, userService)

	history := orderService.GetOrderHistory()
	actual := len(history)
	if actual <= 0 {
		t.Errorf("[TestOrderService_GetOrderHistory] Expected order history, got %d", actual)
	}
	test.CleanupIntegrationTest()
}

func TestOrderService_ImportCSV(t *testing.T) {

	ctx := test.NewIntegrationTestContext()
	exchangeDAO := dao.NewExchangeDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	orderDAO := dao.NewOrderDAO(ctx)

	userMapper := mapper.NewUserMapper()
	orderMapper := mapper.NewOrderMapper(ctx)
	exchangeMapper := mapper.NewExchangeMapper()

	marketcapService := NewMarketCapService(ctx.Logger)
	ethereumService, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO, userMapper)
	assert.Nil(t, err)

	exchangeService := NewExchangeService(ctx, exchangeDAO, userDAO, userMapper, exchangeMapper)
	userService := NewUserService(ctx, userDAO, marketcapService, ethereumService, userMapper, exchangeMapper)
	orderService := NewOrderService(ctx, orderDAO, orderMapper, exchangeService, userService)

	orders, err := orderService.ImportCSV("../test/data/bittrex.csv", "bittrex")
	assert.Nil(t, err)
	assert.Equal(t, true, len(orders) > 0)

	persistedOrders, err := orderDAO.Find()
	assert.Nil(t, err)
	assert.Equal(t, true, len(persistedOrders) > 0)

	test.CleanupIntegrationTest()
}
