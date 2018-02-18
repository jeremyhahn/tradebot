package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
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
			history := ex.GetOrderHistory(&common.CurrencyPair{
				Base:          currencyPair.Base,
				Quote:         currencyPair.Quote,
				LocalCurrency: os.ctx.GetUser().GetLocalCurrency()})
			orders = append(orders, history...)
		}
	}
	sort.Slice(orders, func(i, j int) bool {
		orders[i] = &dto.OrderDTO{
			Id:           orders[i].GetId(),
			Exchange:     orders[i].GetExchange(),
			Date:         orders[i].GetDate(),
			Type:         orders[i].GetType(),
			CurrencyPair: orders[i].GetCurrencyPair(),
			Quantity:     orders[i].GetQuantity(),
			Price:        orders[i].GetPrice(),
			Fee:          orders[i].GetFee(),
			Total:        orders[i].GetTotal()}
		return orders[i].GetDate().After(orders[j].GetDate())
	})
	return orders
}

func (dos *DefaultOrderService) format(f float64, currency string) float64 {
	if currency == "USD" {
		return util.TruncateFloat(f, 2)
	}
	return util.TruncateFloat(f, 8)
}
