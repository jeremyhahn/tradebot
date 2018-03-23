// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
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

	totalWithdrawl := decimal.NewFromFloat(0.0)
	totalDeposit := decimal.NewFromFloat(0.0)

	for _, tx := range transactions {
		qty, _ := decimal.NewFromString(tx.GetQuantity())
		if tx.GetType() == "deposit" {
			totalDeposit = totalDeposit.Add(qty)
		} else if tx.GetType() == "withdrawl" {
			totalWithdrawl = totalWithdrawl.Add(qty)
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit.GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, totalWithdrawl.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}
