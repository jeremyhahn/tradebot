package rest

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type UserRestService interface {
	GetExchanges(w http.ResponseWriter, r *http.Request)
	CreateExchange(w http.ResponseWriter, r *http.Request)
	DeleteExchanges(w http.ResponseWriter, r *http.Request)
}

type UserRestServiceImpl struct {
	middlewareService service.Middleware
	jsonWriter        common.HttpWriter
}

func NewUserRestService(middlewareService service.Middleware, jsonWriter common.HttpWriter) UserRestService {
	return &UserRestServiceImpl{
		middlewareService: middlewareService,
		jsonWriter:        jsonWriter}
}

func (restService *UserRestServiceImpl) createUserService(ctx common.Context) service.UserService {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	ethereumService, _ := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	walletService := service.NewWalletService(ctx, pluginService)
	return service.NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService,
		ethereumService, exchangeService, walletService)
}

func (restService *UserRestServiceImpl) GetExchanges(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
		return
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[UserRestService.GetExchanges]")
	userCryptoExchanges := restService.createUserService(ctx).GetConfiguredExchanges()
	viewModels := make([]*viewmodel.UserCryptoExchange, len(userCryptoExchanges))
	for i, userCryptoExchange := range userCryptoExchanges {
		viewModels[i] = mapper.NewUserExchangeMapper().MapDtoToViewModel(userCryptoExchange)
	}
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: viewModels})
}

func (restService *UserRestServiceImpl) CreateExchange(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
		return
	}
	defer ctx.Close()
	ctx.GetLogger().Debugf("[UserRestServiceImpl.CreateExchange]")
	userService := restService.createUserService(ctx)
	userCryptoExchange, err := userService.CreateExchange(&dto.UserCryptoExchangeDTO{
		UserID: ctx.GetUser().GetId(),
		Name:   r.FormValue("name"),
		Key:    r.FormValue("key"),
		Secret: r.FormValue("secret"),
		Extra:  r.FormValue("extra")})
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
		return
	}
	ctx.GetLogger().Debugf("[UserRestServiceImpl.CreateExchange]")
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: userCryptoExchange})
}

func (restService *UserRestServiceImpl) DeleteExchanges(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
		return
	}
	defer ctx.Close()
	params := mux.Vars(r)
	ctx.GetLogger().Debugf("[UserRestService.DeleteExchange]")
	exchanges := strings.Split(params["id"], ",")
	for _, exchange := range exchanges {
		err = restService.createUserService(ctx).DeleteExchange(exchange)
		if err != nil {
			restService.jsonWriter.Write(w, http.StatusNotFound, common.JsonResponse{
				Success: false,
				Payload: err.Error()})
			return
		}
	}
	if err == nil {
		restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
			Success: true,
			Payload: nil})
		return
	}
	restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
		Success: false,
		Payload: "Internal server error"})
}
