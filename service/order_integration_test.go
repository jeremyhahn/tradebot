// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
)

func TestOrderService_GetOrderHistory(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	exchangeDAO := dao.NewExchangeDAO(ctx)
	exchangeService := NewExchangeService(ctx, exchangeDAO)
	orderService := NewOrderService(ctx, exchangeService)
	history := orderService.GetOrderHistory()
	//for _, o := range history {
	//	util.DUMP(o)
	//}
	actual := len(history)
	if actual <= 0 {
		t.Errorf("[TestOrderService_GetOrderHistory] Expected order history, got %d", actual)
	}
	test.CleanupIntegrationTest()
}
