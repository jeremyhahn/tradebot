package rest

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type TransactionRestService interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
	GetDepositHistory(w http.ResponseWriter, r *http.Request)
	GetWithdrawalHistory(w http.ResponseWriter, r *http.Request)
	GetImportedTransactions(w http.ResponseWriter, r *http.Request)
	Export(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
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
	ctx.GetLogger().Debugf("[TransactionRestService.GetTransactions]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	transactions, err := txService.GetTransactions()
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	restService.formatTransactions(&transactions)
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: &transactions})
}

func (restService *TransactionRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetOrderHistory]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	transactions := txService.GetOrderHistory()
	restService.formatTransactions(&transactions)
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: transactions})
}

func (restService *TransactionRestServiceImpl) GetDepositHistory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetDepositHistory]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	transactions := txService.GetDepositHistory()
	restService.formatTransactions(&transactions)
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: transactions})
}

func (restService *TransactionRestServiceImpl) GetWithdrawalHistory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetWithdrawalHistory]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	transactions := txService.GetWithdrawalHistory()
	restService.formatTransactions(&transactions)
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: transactions})
}

func (restService *TransactionRestServiceImpl) GetImportedTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetImportedTransactions]")

	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}

	transactions := txService.GetImportedTransactions()
	restService.formatTransactions(&transactions)
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: transactions})
}

func (restService *TransactionRestServiceImpl) Import(w http.ResponseWriter, r *http.Request) {

	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		return
	}
	defer ctx.Close()

	ctx.GetLogger().Debugf("[TransactionRestServiceImpl.Import]")

	orderService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("csv")
	if err != nil {
		ctx.GetLogger().Errorf("[TransactionRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
			Success: false, Payload: "Missing CSV form field"})
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("/tmp/%s", handler.Filename)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	_, err = io.Copy(f, file)
	if err != nil {
		ctx.GetLogger().Errorf("[TransactionRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false, Payload: nil})
		return
	}

	records, err := orderService.ImportCSV(filename, r.FormValue("exchange"))
	if err != nil {
		ctx.GetLogger().Errorf("[TransactionRestServiceImpl.Import] %s", err.Error())
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false, Payload: err.Error()})
		return
	}

	var response []viewmodel.Transaction
	for _, orderDTO := range records {
		response = append(response,
			orderService.GetMapper().MapTransactionDtoToViewModel(orderDTO))
	}

	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: len(records) > 0,
		Payload: response})
}

func (restService *TransactionRestServiceImpl) Export(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		return
	}
	defer ctx.Close()

	filename := fmt.Sprintf("%s-%s", ctx.GetUser().GetUsername(), "-8949.csv")
	file, err := os.Create(filename)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	ctx.GetLogger().Debugf("[TransactionHistoryRestService.Export]")
	orderService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	orderHistory := orderService.GetOrderHistory()
	var orders []viewmodel.Transaction
	for _, order := range orderHistory {
		orders = append(orders, orderService.GetMapper().MapTransactionDtoToViewModel(order))
	}
}

func (restService *TransactionRestServiceImpl) formatTransactions(txs *[]common.Transaction) {
	for i, tx := range *txs {
		id := tx.GetId()
		if id == "" || len(id) <= 0 {
			id = fmt.Sprintf("%d", i)
		}
		txID := fmt.Sprintf("%s-%s", tx.GetNetwork(), id)
		(*txs)[i] = &dto.TransactionDTO{
			Id:                   txID,
			Date:                 tx.GetDate(),
			Type:                 strings.Title(tx.GetType()),
			CurrencyPair:         tx.GetCurrencyPair(),
			Network:              tx.GetNetwork(),
			NetworkDisplayName:   tx.GetNetworkDisplayName(),
			Quantity:             tx.GetQuantity(),
			QuantityCurrency:     tx.GetQuantityCurrency(),
			FiatQuantity:         tx.GetFiatQuantity(),
			FiatQuantityCurrency: tx.GetFiatQuantityCurrency(),
			Price:                tx.GetPrice(),
			PriceCurrency:        tx.GetPriceCurrency(),
			FiatPrice:            tx.GetFiatPrice(),
			FiatPriceCurrency:    tx.GetFiatPriceCurrency(),
			Fee:                  tx.GetFee(),
			FeeCurrency:          tx.GetFeeCurrency(),
			FiatFee:              tx.GetFiatFee(),
			FiatFeeCurrency:      tx.GetFiatFeeCurrency(),
			Total:                tx.GetTotal(),
			TotalCurrency:        tx.GetTotalCurrency(),
			FiatTotal:            tx.GetFiatTotal(),
			FiatTotalCurrency:    tx.GetFiatTotalCurrency()}
	}
}

func (restService *TransactionRestServiceImpl) createTransactionService(ctx common.Context) (service.TransactionService, error) {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	transactionDAO := dao.NewTransactionDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	transactionMapper := mapper.NewTransactionMapper(ctx)
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := service.NewFiatPriceService(ctx, exchangeService)
	if err != nil {
		return nil, err
	}
	ethereumService, err := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	if err != nil {
		return nil, err
	}
	return service.NewTransactionService(ctx, transactionDAO, transactionMapper, exchangeService, ethereumService, fiatPriceService), nil
}
