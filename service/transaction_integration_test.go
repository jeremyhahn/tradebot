// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

/*
func TestTransactionService_Synchronize(t *testing.T) {
	_, transactionService := createTransactionService()
	txs, err := transactionService.Synchronize()
	assert.Nil(t, err)
	assert.Equal(t, true, len(txs) > 0)
	CleanupIntegrationTest()
}*/

/*
func TestTransactionService_GetOrderHistory(t *testing.T) {
	_, transactionService := createTransactionService()
	history := transactionService.GetOrderHistory()
	assert.NotNil(t, history)
	assert.Equal(t, true, len(history) > 0)
	CleanupIntegrationTest()
}*/

func TestTransactionService_ImportCSV(t *testing.T) {

	transactionDAO, transactionService := createTransactionService()

	transactions, err := transactionService.ImportCSV("../test/data/bittrex.csv", "Bittrex")
	assert.Nil(t, err)
	assert.Equal(t, true, len(transactions) > 0)

	persistedTransactions, err := transactionDAO.Find()
	assert.Nil(t, err)
	assert.Equal(t, true, len(persistedTransactions) == 2)

	CleanupIntegrationTest()
}

func createTransactionService() (dao.TransactionDAO, TransactionService) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	transactionDAO := dao.NewTransactionDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	transactionMapper := mapper.NewTransactionMapper(ctx)
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := NewMarketCapService(ctx)
	pluginService := CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	ethereumService, _ := NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	fiatPriceService, _ := NewFiatPriceService(ctx, exchangeService)
	walletService := NewWalletService(ctx, pluginService, fiatPriceService)
	userService := NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService, ethereumService, exchangeService, walletService)
	return transactionDAO, NewTransactionService(ctx, transactionDAO, transactionMapper,
		exchangeService, userService, ethereumService, fiatPriceService)
}
