// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestExchangeService_GetExchanges(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	exchangeDAO := dao.NewExchangeDAO(ctx)
	exchangeService := NewExchangeService(ctx, exchangeDAO)
	exchanges := exchangeService.GetExchanges(ctx.User, &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: ctx.User.LocalCurrency})
	assert.Equal(t, 4, len(exchanges))
	test.CleanupMockContext()
}
