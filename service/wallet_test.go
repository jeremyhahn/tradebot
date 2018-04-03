// +build integration

package service

import (
	"os"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWalletService_BTC(t *testing.T) {
	ctx := NewIntegrationTestContext()
	btc := createWalletService(t, ctx, "BTC", os.Getenv("BTC_ADDRESS"))
	zero := decimal.NewFromFloat(0)
	assert.Equal(t, true, btc.GetPrice().GreaterThan(zero))
	wallet, err := btc.GetWallet()
	assert.Nil(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, "BTC", wallet.GetCurrency())
	assert.Equal(t, true, wallet.GetValue().GreaterThan(zero))
	assert.Equal(t, true, wallet.GetBalance().GreaterThan(zero))
	CleanupIntegrationTest()
}

func TestWalletService_ETH(t *testing.T) {
	ctx := NewIntegrationTestContext()
	eth := createWalletService(t, ctx, "ETH", os.Getenv("ETH_ADDRESS"))
	zero := decimal.NewFromFloat(0)
	assert.Equal(t, true, eth.GetPrice().GreaterThan(zero))
	wallet, err := eth.GetWallet()
	assert.Nil(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, "ETH", wallet.GetCurrency())
	assert.Equal(t, true, wallet.GetValue().GreaterThan(zero))
	assert.Equal(t, true, wallet.GetBalance().GreaterThan(zero))
	CleanupIntegrationTest()
}

func TestWalletService_XRP(t *testing.T) {
	ctx := NewIntegrationTestContext()
	xrp := createWalletService(t, ctx, "XRP", os.Getenv("XRP_ADDRESS"))
	zero := decimal.NewFromFloat(0)
	assert.Equal(t, true, xrp.GetPrice().GreaterThan(zero))
	wallet, err := xrp.GetWallet()
	assert.Nil(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, "XRP", wallet.GetCurrency())
	assert.Equal(t, true, wallet.GetValue().GreaterThan(zero))
	assert.Equal(t, true, wallet.GetBalance().GreaterThan(zero))
	CleanupIntegrationTest()
}

func createWalletService(t *testing.T, ctx common.Context, currency, address string) common.Wallet {
	pluginDAO := dao.NewPluginDAO(ctx)
	pluginMapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, pluginMapper)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, _ := NewFiatPriceService(ctx, exchangeService)
	service := NewWalletService(ctx, pluginService, fiatPriceService)
	wallet, err := service.CreateWallet(currency, address)
	assert.Nil(t, err)
	assert.NotNil(t, wallet)
	return wallet
}
