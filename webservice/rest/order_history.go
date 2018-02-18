package rest

import (
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type OrderHistoryResponse struct {
	Payload interface{}
	RestResponse
}

type OrderHistoryRestService interface {
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
}

type OrderHistoryRestServiceImpl struct {
	ctx              *common.Context
	marketcapService *service.MarketCapService
	exchangeService  service.ExchangeService
	userService      service.UserService
	jsonWriter       common.HttpWriter
}

func NewOrderHistoryRestService(ctx *common.Context, marketcapService *service.MarketCapService,
	exchangeService service.ExchangeService, userService service.UserService, jsonWriter common.HttpWriter) OrderHistoryRestService {
	return &OrderHistoryRestServiceImpl{
		ctx:              ctx,
		marketcapService: marketcapService,
		exchangeService:  exchangeService,
		jsonWriter:       jsonWriter}
}

func (ohrs *OrderHistoryRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ohrs.ctx.Logger.Debugf("[OrderHistoryRestService.GetOrderHistory]")
	service := service.NewOrderService(ohrs.ctx, ohrs.exchangeService, ohrs.userService)
	history := service.GetOrderHistory()
	var orders []viewmodel.Order
	for _, order := range history {
		orders = append(orders, viewmodel.Order{
			Id:           order.GetId(),
			Exchange:     order.GetExchange(),
			Date:         order.GetDate().Format(common.TIME_DISPLAY_FORMAT),
			Type:         order.GetType(),
			CurrencyPair: order.GetCurrencyPair(),
			Quantity:     order.GetQuantity(),
			Price:        order.GetPrice(),
			Fee:          order.GetFee(),
			Total:        order.GetTotal()})
	}
	ohrs.jsonWriter.Write(w, http.StatusOK, RestResponse{
		Success: true,
		Payload: orders})
}
