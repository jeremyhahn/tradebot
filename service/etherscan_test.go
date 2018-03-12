// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/stretchr/testify/assert"
)

func TestEtherScan_GetWallet(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(), NewMarketCapService(ctx))
	assert.Nil(t, err)

	wallet, err := service.GetWallet(os.Getenv("ETH_ADDRESS"))
	assert.Nil(t, err)

	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, true, wallet.GetBalance() > 0)
	assert.Equal(t, true, wallet.GetValue() > 0)

	test.CleanupIntegrationTest()
}

func TestEtherScanService_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	walletAddress := os.Getenv("ETH_ADDRESS")

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(), NewMarketCapService(ctx))
	assert.Nil(t, err)

	transactions, err := service.GetTransactions(walletAddress)

	service.GetTransactions(walletAddress)

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	totalDeposit := 0.0
	totalWithdrawl := 0.0
	totalFees := 0.0

	for _, tx := range transactions {
		if tx.GetType() == "deposit" {
			totalDeposit += tx.GetAmount()
		} else if tx.GetType() == "withdrawl" {
			totalWithdrawl += tx.GetAmount()
			totalFees += tx.GetFee()
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit > 0)
	assert.Equal(t, true, totalWithdrawl > 0)

	txSum := totalDeposit - totalWithdrawl - totalFees

	wallet, err := service.GetWallet(walletAddress)
	assert.Nil(t, err)

	actual := util.TruncateFloat(wallet.GetBalance(), 8)
	expected := util.TruncateFloat(txSum, 8)
	assert.Equal(t, actual, expected)

	test.CleanupIntegrationTest()
}

func TestEtherScan_GetToken(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(), NewMarketCapService(ctx))
	assert.Nil(t, err)

	token, err := service.GetToken(os.Getenv("ETH_ADDRESS"), os.Getenv("TOKEN_ADDRESS"))
	assert.Nil(t, err)

	assert.Equal(t, true, token.GetBalance() > 0)
	//fmt.Printf("balance: %.8f\n", balance)

	test.CleanupIntegrationTest()
}

func TestEtherScanService_GetTokenTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service, err := NewEtherscanService(ctx, dao.NewUserDAO(ctx), mapper.NewUserMapper(), NewMarketCapService(ctx))
	assert.Nil(t, err)

	transactions, err := service.GetTokenTransactions(os.Getenv("TOKEN_INTERMEDIARY"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	test.CleanupIntegrationTest()
}
