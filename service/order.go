package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultOrderService struct {
	ctx             *common.Context
	orderDAO        dao.OrderDAO
	orderMapper     mapper.OrderMapper
	exchangeService ExchangeService
	userService     UserService
	OrderService
}

func NewOrderService(ctx *common.Context, orderDAO dao.OrderDAO, orderMapper mapper.OrderMapper,
	exchangeService ExchangeService, userService UserService) OrderService {
	return &DefaultOrderService{
		ctx:             ctx,
		orderDAO:        orderDAO,
		orderMapper:     orderMapper,
		exchangeService: exchangeService,
		userService:     userService}
}

func (os *DefaultOrderService) GetMapper() mapper.OrderMapper {
	return os.orderMapper
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
	orderEntities, err := os.orderDAO.Find()
	if err != nil {
		os.ctx.Logger.Errorf("[OrderService.GetOrderHistory] %s", err.Error())
	} else {
		for _, entity := range orderEntities {
			orders = append(orders, os.orderMapper.MapOrderEntityToDto(&entity))
		}
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].GetDate().After(orders[j].GetDate())
	})
	return orders
}

func (dos *DefaultOrderService) ImportCSV(file, exchangeName string) ([]common.Order, error) {
	dos.ctx.Logger.Debugf("[OrderService.ImportCSV] Creating %s exchange service", exchangeName)
	exchange := dos.exchangeService.GetExchange(dos.ctx.GetUser(), exchangeName)
	orderDTOs, err := exchange.ParseImport(file)
	if err != nil {
		return nil, err
	}
	for _, dto := range orderDTOs {
		entity := dos.orderMapper.MapOrderDtoToEntity(dto)
		dos.orderDAO.Create(entity)
	}
	return orderDTOs, nil
}
