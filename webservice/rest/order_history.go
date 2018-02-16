package rest

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	ctx             *common.Context
	exchangeService service.ExchangeService
	userService     service.UserService
	jsonWriter      common.HttpWriter
}

func NewOrderHistoryRestService(ctx *common.Context, exchangeService service.ExchangeService,
	userService service.UserService, jsonWriter common.HttpWriter) OrderHistoryRestService {
	return &OrderHistoryRestServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService,
		jsonWriter:      jsonWriter}
}

func (ohrs *OrderHistoryRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ohrs.ctx.Logger.Debugf("[OrderHistoryRestService.GetOrderHistory]")
	service := service.NewOrderService(ohrs.ctx, ohrs.exchangeService, ohrs.userService)
	history := service.GetOrderHistory()
	var orders []viewmodel.Order
	for _, order := range history {
		orders = append(orders, viewmodel.Order{
			Id:       order.GetId(),
			Date:     order.GetDate().Format(common.TIME_DISPLAY_FORMAT),
			Type:     order.GetType(),
			Price:    order.GetPrice(),
			Currency: order.GetCurrency(),
			Quantity: order.GetQuantity(),
			Exchange: order.GetExchange()})
	}
	ohrs.jsonWriter.Write(w, http.StatusOK, RestResponse{
		Success: true,
		Payload: orders})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	var restResponse RestResponse
	response, err := json.Marshal(payload)
	if err != nil {
		var buf bytes.Buffer
		restResponse = RestResponse{Error: err.Error()}
		data, _ := json.Marshal(restResponse)
		binary.Write(&buf, binary.BigEndian, data)
		w.Write(buf.Bytes())
	} else {
		w.Write(response)
	}
}
