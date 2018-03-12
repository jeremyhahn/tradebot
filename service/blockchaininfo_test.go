// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainInfo_GetPrice(t *testing.T) {
	ctx := test.NewUnitTestContext()

	blockchainInfo := NewBlockchainInfo(ctx)
	price := blockchainInfo.GetPrice()

	assert.Equal(t, true, price > 0)
}

func TestBlockchainInfo_GetBalance(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	blockchainInfo := NewBlockchainInfo(ctx)
	balance := blockchainInfo.GetBalance(os.Getenv("BTC_ADDRESS"))

	assert.Equal(t, true, len(balance.GetAddress()) > 0)
	assert.Equal(t, true, balance.GetBalance() > 0)
	assert.Equal(t, true, balance.GetValue() > 0)

	test.CleanupIntegrationTest()
}

func TestBlockchainInfo_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	blockchainInfo := NewBlockchainInfo(ctx)
	transactions, err := blockchainInfo.GetTransactions(os.Getenv("BTC_ADDRESS"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	totalWithdrawl := 0.0
	totalDeposit := 0.0

	for _, tx := range transactions {
		if tx.GetType() == "deposit" {
			totalDeposit += tx.GetAmount()
		} else if tx.GetType() == "withdrawl" {
			totalWithdrawl += tx.GetAmount()
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit > 0)
	assert.Equal(t, true, totalWithdrawl > 0)

	test.CleanupIntegrationTest()
}
