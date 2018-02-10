package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
)

type DefaultOrderService struct {
	ctx             *common.Context
	exchangeService ExchangeService
	OrderService
}

func NewOrderService(ctx *common.Context, exchangeService ExchangeService) OrderService {
	return &DefaultOrderService{
		ctx:             ctx,
		exchangeService: exchangeService}
}

func (os DefaultOrderService) GetOrderHistory() []common.Order {
	var orders []common.Order

	// TODO: Look up currency pairs from DB
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         os.ctx.User.GetLocalCurrency(),
		LocalCurrency: os.ctx.User.GetLocalCurrency()}

	exchanges := os.exchangeService.GetExchanges(os.ctx.User)
	for _, ex := range exchanges {
		orders = append(orders, ex.GetOrderHistory(currencyPair)...)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].GetDate().Before(orders[j].GetDate())
	})
	return orders
}
