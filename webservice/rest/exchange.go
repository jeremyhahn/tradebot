package rest

import (
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
)

type ExchangeRestService interface {
	GetDisplayNames(w http.ResponseWriter, r *http.Request)
}

type ExchangeRestServiceImpl struct {
	middlewareService service.Middleware
	jsonWriter        common.HttpWriter
}

func NewExchangeRestService(middlewareService service.Middleware,
	jsonWriter common.HttpWriter) ExchangeRestService {
	return &ExchangeRestServiceImpl{
		middlewareService: middlewareService,
		jsonWriter:        jsonWriter}
}

func (restService *ExchangeRestServiceImpl) createExchangeService(ctx common.Context) service.ExchangeService {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
	return service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
}

func (restService *ExchangeRestServiceImpl) GetDisplayNames(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
		return
	}
	defer ctx.Close()
	exchangeService := restService.createExchangeService(ctx)
	names, err := exchangeService.GetDisplayNames()
	if err != nil {
		restService.jsonWriter.Write(w, http.StatusInternalServerError, common.JsonResponse{
			Success: false,
			Payload: err.Error()})
	}
	ctx.GetLogger().Debugf("[ExchangeRestService.GetDisplayNames]")
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: names})
}
