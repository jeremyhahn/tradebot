// +build integration

package main

import (
	"os"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
)

func TestBchWallet_GetTransactions(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := service.NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	wallet := CreateBchWallet(&common.WalletParams{
		Context:          ctx,
		Address:          os.Getenv("BCH_ADDRESS"),
		MarketCapService: marketcapService,
		FiatPriceService: fiatPriceService})

	transactions, err := wallet.GetTransactions()
	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	for _, tx := range transactions {
		util.DUMP(tx)
	}

	test.CleanupIntegrationTest()
}
