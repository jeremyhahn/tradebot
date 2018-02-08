package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
)

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

func (os OrderServiceImpl) GetOrderHistory() []common.Order {
	var orders []common.Order

	// TODO: Look up currency pairs from DB
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         os.ctx.User.LocalCurrency,
		LocalCurrency: os.ctx.User.LocalCurrency}

	exchanges := os.exchangeService.GetExchanges(os.ctx.User)
	for _, ex := range exchanges {
		orders = append(orders, ex.GetOrderHistory(currencyPair)...)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].GetDate().Before(orders[j].GetDate())
	})
	return orders
}
