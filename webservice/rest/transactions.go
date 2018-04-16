package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jeremyhahn/tradebot/accounting"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type TransactionRestService interface {
	GetHistory(w http.ResponseWriter, r *http.Request)
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
	GetDepositHistory(w http.ResponseWriter, r *http.Request)
	GetWithdrawalHistory(w http.ResponseWriter, r *http.Request)
	GetImportedTransactions(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	Synchronize(w http.ResponseWriter, r *http.Request)
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

func (restService *TransactionRestServiceImpl) Synchronize(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.Synchronize]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	txs, err := txService.Synchronize()
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
}

func (restService *TransactionRestServiceImpl) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.GetHistory]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	txs, err := txService.GetHistory("desc")
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
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
	txs := txService.GetOrderHistory()
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
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
	txs := txService.GetDepositHistory()
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
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
	txs := txService.GetWithdrawalHistory()
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
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

	txs := txService.GetImportedTransactions()
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: restService.formatTransactions(ctx, txs)})
}

func (restService *TransactionRestServiceImpl) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[TransactionRestService.UpdateCategory]")
	txService, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	params := mux.Vars(r)
	ctx.GetLogger().Debugf("[TransactionRestService.UpdateCategory]")
	if params["id"] == "" {
		restService.jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
			Success: false,
			Payload: "Transaction id required"})
		return
	}
	r.ParseMultipartForm(32 << 20)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	category := r.FormValue("category")
	if category == "" {
		restService.jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
			Success: false,
			Payload: "Category required"})
		return
	}
	err = txService.UpdateCategory(params["id"], category)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: true})
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
	ctx.GetLogger().Debugf("[TransactionHistoryRestService.Export]")
	service, err := restService.createTransactionService(ctx)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	transactions, err := service.GetHistory("asc")
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	fifo := accounting.NewFifoReport(ctx, transactions)
	location := time.Now().Location()
	start := time.Date(2017, 01, 01, 0, 0, 0, 0, location)
	end := time.Date(2017, 12, 31, 0, 0, 0, 0, location)
	form8948 := fifo.Run(start, end)

	filename := fmt.Sprintf("/tmp/%s-%s", ctx.GetUser().GetUsername(), "-8949.csv")
	form8948.WriteCSV(filename)

	csvBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}

	restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
		Success: true,
		Payload: string(csvBytes)})

	os.Remove(filename)
}

func (restService *TransactionRestServiceImpl) formatTransactions(ctx common.Context, txs []common.Transaction) []viewmodel.Transaction {
	mapper := mapper.NewTransactionMapper(ctx)
	var viewModels []viewmodel.Transaction
	for _, tx := range txs {
		viewModels = append(viewModels, mapper.MapTransactionDtoToViewModel(tx))
	}
	return viewModels
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
	walletService := service.NewWalletService(ctx, pluginService, fiatPriceService)
	userService := service.NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService, ethereumService, exchangeService, walletService)
	return service.NewTransactionService(ctx, transactionDAO, transactionMapper, exchangeService, userService, ethereumService, fiatPriceService), nil
}
