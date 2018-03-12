// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestRippleService_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	marketcapService := NewMarketCapService(ctx)
	rippleService := NewRipple(ctx, marketcapService)

	transactions, err := rippleService.GetTransactions(os.Getenv("XRP_ADDRESS"))
	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, true, len(transactions) > 1)

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
