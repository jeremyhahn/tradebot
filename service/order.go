package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
)

type DefaultOrderService struct {
	ctx             *common.Context
	exchangeService ExchangeService
	userService     UserService
	OrderService
}

func NewOrderService(ctx *common.Context, exchangeService ExchangeService, userService UserService) OrderService {
	return &DefaultOrderService{
		ctx:             ctx,
		exchangeService: exchangeService,
		userService:     userService}
}

func (os *DefaultOrderService) GetOrderHistory() []common.Order {
	var orders []common.Order
	exchanges := os.exchangeService.GetExchanges(os.ctx.User)
	for _, ex := range exchanges {
		if ex.GetName() == "gdax" {
			balances, _ := ex.GetBalances()
			for _, coin := range balances {
				currencyPair := &common.CurrencyPair{
					Base:          coin.GetCurrency(),
					Quote:         os.ctx.GetUser().GetLocalCurrency(),
					LocalCurrency: os.ctx.GetUser().GetLocalCurrency()}
				orders = append(orders, ex.GetOrderHistory(currencyPair)...)
			}
			continue
		}
		currencyPairs, err := os.exchangeService.GetCurrencyPairs(os.ctx.GetUser(), ex.GetName())
		if err != nil {
			os.ctx.Logger.Errorf("[OrderService.GetOrderHistory] %s", err.Error())
			return orders
		}
		for _, currencyPair := range currencyPairs {
			orders = append(orders, ex.GetOrderHistory(&currencyPair)...)
		}
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].GetDate().After(orders[j].GetDate())
	})
	return orders
}
