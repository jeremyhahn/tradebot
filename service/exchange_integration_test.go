// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestExchangeService_GetExchanges(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	exchangeDAO := dao.NewExchangeDAO(ctx)
	exchangeService := NewExchangeService(ctx, exchangeDAO)
	exchanges := exchangeService.GetExchanges(ctx.User)
	assert.Equal(t, 3, len(exchanges))
	test.CleanupIntegrationTest()
}
