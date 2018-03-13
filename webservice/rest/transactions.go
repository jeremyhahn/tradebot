package rest

import (
	"net/http"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
)

type TransactionRestService interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
}

type TransactionRestServiceImpl struct {
	middlewareService service.Middleware
	jsonWriter        common.HttpWriter
	TransactionRestService
}

func NewTransactionRestService(middlewareService service.Middleware, jsonWriter common.HttpWriter) TransactionRestService {
	return &TransactionRestServiceImpl{
		middlewareService: middlewareService,
		jsonWriter:        jsonWriter}
}

func (restService *TransactionRestServiceImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetOrderHistory]")
	orderService, _, _ := restService.createOrderEthereumServices(ctx)
	orderService, ethereumService, _ := restService.createOrderEthereumServices(ctx)
	orderHistory := orderService.GetOrderHistory()
	transactions, err := ethereumService.GetTransactions()
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: true,
			Payload: err.Error()})
		return
	}
	for _, order := range orderHistory {
		transactions = append(transactions, &dto.TransactionDTO{
			Date:               order.GetDate(),
			Type:               strings.Title(order.GetType()),
			CurrencyPair:       order.GetCurrencyPair(),
			Source:             strings.Title(order.GetExchange()),
			Amount:             order.GetQuantity(),
			AmountCurrency:     order.GetQuantityCurrency(),
			Fee:                order.GetFee(),
			FeeCurrency:        order.GetFeeCurrency(),
			Total:              order.GetTotal(),
			TotalCurrency:      order.GetTotalCurrency(),
			HistoricalPrice:    order.GetHistoricalPrice(),
			HistoricalCurrency: order.GetHistoricalCurrency()})
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: transactions})
}

func (restService *TransactionRestServiceImpl) createOrderEthereumServices(ctx common.Context) (service.OrderService,
	service.EthereumService, common.PriceHistoryService) {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	orderDAO := dao.NewOrderDAO(ctx)
	userMapper := mapper.NewUserMapper()
	orderMapper := mapper.NewOrderMapper(ctx)
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	priceHistoryService := service.NewPriceHistoryService(ctx)
	exchangeService := service.NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper, priceHistoryService)
	ethereumService, _ := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService)
	userService := service.NewUserService(ctx, userDAO, pluginDAO, marketcapService, ethereumService, userMapper, userExchangeMapper, priceHistoryService)
	return service.NewOrderService(ctx, orderDAO, orderMapper, exchangeService, userService), ethereumService, priceHistoryService
}
