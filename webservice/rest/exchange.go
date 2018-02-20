package rest

import (
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type ExchangeRestService interface {
	GetExchanges(w http.ResponseWriter, r *http.Request)
}

type ExchangeRestServiceImpl struct {
	ctx             *common.Context
	exchangeService service.ExchangeService
	userService     service.UserService
	jsonWriter      common.HttpWriter
}

func NewExchangeRestService(ctx *common.Context, exchangeService service.ExchangeService,
	userService service.UserService, jsonWriter common.HttpWriter) ExchangeRestService {
	return &ExchangeRestServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService,
		jsonWriter:      jsonWriter}
}

func (irs *ExchangeRestServiceImpl) GetExchanges(w http.ResponseWriter, r *http.Request) {
	irs.ctx.Logger.Debugf("[ExchangeRestService.GetExchanges]")
	exchanges := irs.exchangeService.GetDisplayNames(irs.ctx.GetUser())
	irs.jsonWriter.Write(w, http.StatusOK, RestResponse{
		Success: true,
		Payload: exchanges})
}
