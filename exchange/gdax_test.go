// +build integration

package exchange

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestGDAX_GetBalances(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	userDAO := dao.NewUserDAO(ctx)
	cryptoExchange := userDAO.GetExchange(ctx.User, "gdax")
	gdax := NewGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)
	balance, netWorth := gdax.GetBalances()
	assert.Equal(t, true, len(balance) > 0)
	assert.Equal(t, true, netWorth > 0)
	test.CleanupIntegrationTest()
}

func TestGDAX_GetOrderHistory(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	userDAO := dao.NewUserDAO(ctx)
	cryptoExchange := userDAO.GetExchange(ctx.User, "gdax")
	gdax := NewGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)
	orders := gdax.GetOrderHistory(&common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"})
	assert.Equal(t, true, len(orders) > 0)
	test.CleanupIntegrationTest()
}

func TestGDAX_GetPriceHistory(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	userDAO := dao.NewUserDAO(ctx)
	cryptoExchange := userDAO.GetExchange(ctx.User, "gdax")
	gdax := NewGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	history := gdax.GetPriceHistory(currencyPair, start, end, 900)
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

func TestGDAX_GetExchange(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	userDAO := dao.NewUserDAO(ctx)
	cryptoExchange := userDAO.GetExchange(ctx.User, "gdax")
	gdax := NewGDAX(ctx, cryptoExchange)
	assert.NotNil(t, gdax)
	exchange := gdax.GetExchange()
	assert.NotNil(t, exchange)
	assert.Equal(t, "gdax", exchange.Name)
	test.CleanupIntegrationTest()
}
