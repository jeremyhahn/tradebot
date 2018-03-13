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
	userExchangeMapper := mapper.NewUserExchangeMapper()
	priceHistoryService := service.NewPriceHistoryService(ctx)
	return service.NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper, priceHistoryService)
}

func (restService *ExchangeRestServiceImpl) GetDisplayNames(w http.ResponseWriter, r *http.Request) {
	ctx, err := restService.middlewareService.CreateContext(w, r)
	if err != nil {
		RestError(w, r, err, restService.jsonWriter)
		return
	}
	defer ctx.Close()
	exchangeService := restService.createExchangeService(ctx)
	ctx.GetLogger().Debugf("[ExchangeRestService.GetDisplayNames]")
	restService.jsonWriter.Write(w, http.StatusOK, common.JsonResponse{
		Success: true,
		Payload: exchangeService.GetDisplayNames()})
}
