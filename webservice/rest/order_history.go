package rest

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type OrderHistoryRestService interface {
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
}

type OrderHistoryRestServiceImpl struct {
	middlewareService service.Middleware
	jsonWriter        common.HttpWriter
	OrderHistoryRestService
}

func NewOrderHistoryRestService(middlewareService service.Middleware, jsonWriter common.HttpWriter) OrderHistoryRestService {
	return &OrderHistoryRestServiceImpl{
		middlewareService: middlewareService,
		jsonWriter:        jsonWriter}
}

func (restService *OrderHistoryRestServiceImpl) createOrderService(ctx common.Context) service.OrderService {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	orderDAO := dao.NewOrderDAO(ctx)

	userMapper := mapper.NewUserMapper()
	orderMapper := mapper.NewOrderMapper(ctx)
	userExchangeMapper := mapper.NewUserExchangeMapper()

	marketcapService := service.NewMarketCapService(ctx)
	exchangeService := service.NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper)
	ethereumService, _ := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService)
	userService := service.NewUserService(ctx, userDAO, pluginDAO, marketcapService, ethereumService, userMapper, userExchangeMapper)

	return service.NewOrderService(ctx, orderDAO, orderMapper, exchangeService, userService)
}

func (restService *OrderHistoryRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[OrderHistoryRestService.GetOrderHistory]")
	orderService := restService.createOrderService(ctx)
	orderHistory := orderService.GetOrderHistory()
	var orders []viewmodel.Order
	for _, order := range orderHistory {
		orders = append(orders, orderService.GetMapper().MapOrderDtoToViewModel(order))
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: orders})
}

func (restService *OrderHistoryRestServiceImpl) Export(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		return
	}
	defer ctx.Close()

	filename := fmt.Sprintf("%s-%s", ctx.GetUser().GetUsername(), "-8949.csv")
	file, err := os.Create(filename)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	ctx.GetLogger().Debugf("[OrderHistoryRestService.Export]")
	orderService := restService.createOrderService(ctx)
	orderHistory := orderService.GetOrderHistory()
	var orders []viewmodel.Order
	for _, order := range orderHistory {
		orders = append(orders, orderService.GetMapper().MapOrderDtoToViewModel(order))
	}
}

func (restService *OrderHistoryRestServiceImpl) Import(w http.ResponseWriter, r *http.Request) {

	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		return
	}
	defer ctx.Close()

	ctx.GetLogger().Debugf("[OrderHistoryRestServiceImpl.Import]")

	orderService := restService.createOrderService(ctx)

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("csv")
	if err != nil {
		ctx.GetLogger().Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
			Success: false, Payload: "Missing CSV form field"})
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("/tmp/%s", handler.Filename)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	_, err = io.Copy(f, file)
	if err != nil {
		ctx.GetLogger().Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false, Payload: nil})
		return
	}

	records, err := orderService.ImportCSV(filename, r.FormValue("exchange"))
	if err != nil {
		ctx.GetLogger().Errorf("[OrderHistoryRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false, Payload: err.Error()})
		return
	}

	var response []viewmodel.Order
	for _, orderDTO := range records {
		response = append(response,
			orderService.GetMapper().MapOrderDtoToViewModel(orderDTO))
	}

	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: len(records) > 0,
		Payload: response})
}
