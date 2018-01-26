package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
)

type OrderService interface {
	GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order
}

type OrderServiceImpl struct {
	ctx             *common.Context
	exchangeService ExchangeService
	OrderService
}

func NewOrderService(ctx *common.Context, exchangeService ExchangeService) OrderService {
	return &OrderServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService}
}

func (os OrderServiceImpl) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order {
	var orders []common.Order
	exchanges := os.exchangeService.GetExchanges(os.ctx.User, currencyPair)
	for _, ex := range exchanges {
		orders = append(orders, ex.GetOrderHistory()...)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Date.Before(orders[j].Date)
	})
	return orders
}
