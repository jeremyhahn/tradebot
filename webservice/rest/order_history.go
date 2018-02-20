package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type OrderHistoryRestService interface {
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
}

type OrderHistoryRestServiceImpl struct {
	ctx          *common.Context
	orderService service.OrderService
	jsonWriter   common.HttpWriter
}

func NewOrderHistoryRestService(ctx *common.Context, orderService service.OrderService,
	jsonWriter common.HttpWriter) OrderHistoryRestService {
	return &OrderHistoryRestServiceImpl{
		ctx:          ctx,
		orderService: orderService,
		jsonWriter:   jsonWriter}
}

func (ohrs *OrderHistoryRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ohrs.ctx.Logger.Debugf("[OrderHistoryRestService.GetOrderHistory]")
	history := ohrs.orderService.GetOrderHistory()
	var orders []viewmodel.Order
	for _, order := range history {
		orders = append(orders, ohrs.orderService.GetMapper().MapOrderDtoToViewModel(order))
	}
	ohrs.jsonWriter.Write(w, http.StatusOK, RestResponse{
		Success: true,
		Payload: orders})
}

func (ohrs *OrderHistoryRestServiceImpl) Import(w http.ResponseWriter, r *http.Request) {
	ohrs.ctx.Logger.Debugf("[OrderHistoryRestServiceImpl.Import]")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("csv")
	if err != nil {
		ohrs.ctx.Logger.Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		ohrs.jsonWriter.Write(w, http.StatusBadRequest, RestResponse{
			Success: false, Payload: "Missing CSV form field"})
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%s/%s", common.TMP_DIR, handler.Filename)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	_, err = io.Copy(f, file)
	if err != nil {
		ohrs.ctx.Logger.Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		ohrs.jsonWriter.Write(w, http.StatusInternalServerError, RestResponse{
			Success: false, Payload: nil})
		return
	}

	records, err := ohrs.orderService.ImportCSV(filename, r.FormValue("exchange"))
	if err != nil {
		ohrs.ctx.Logger.Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		ohrs.jsonWriter.Write(w, http.StatusInternalServerError, RestResponse{
			Success: false, Payload: err.Error()})
		return
	}

	var response []viewmodel.Order
	for _, orderDTO := range records {
		response = append(response,
			ohrs.orderService.GetMapper().MapOrderDtoToViewModel(orderDTO))
	}

	ohrs.jsonWriter.Write(w, http.StatusOK, RestResponse{
		Success: len(records) > 0,
		Payload: response})
}
