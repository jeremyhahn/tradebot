// +build integration

package service

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRippleService_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	marketcapService := NewMarketCapService(ctx)
	rippleService := NewRippleService(ctx, dao.NewUserDAO(ctx), marketcapService)

	transactions, err := rippleService.GetTransactions()
	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, true, len(transactions) > 1)

	totalWithdrawl := decimal.NewFromFloat(0)
	totalDeposit := decimal.NewFromFloat(0)

	for _, tx := range transactions {
		qty, err := decimal.NewFromString(tx.GetQuantity())
		if err != nil {
			panic(err)
		}
		if tx.GetType() == common.DEPOSIT_ORDER_TYPE {
			totalDeposit = totalDeposit.Add(qty)
		} else if tx.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			totalWithdrawl = totalWithdrawl.Add(qty)
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}

	assert.Equal(t, true, totalDeposit.GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, totalWithdrawl.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}
