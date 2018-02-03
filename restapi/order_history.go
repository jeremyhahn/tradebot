package restapi

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type OrderHistoryRestService interface {
	GetOrderHistory(w http.ResponseWriter, r *http.Request)
}

type OrderHistoryRestServiceImpl struct {
	ctx             *common.Context
	exchangeService service.ExchangeService
}

func NewOrderHistoryRestService(ctx *common.Context, exchangeService service.ExchangeService) OrderHistoryRestService {
	return &OrderHistoryRestServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService}
}

func (ohrs *OrderHistoryRestServiceImpl) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	ohrs.ctx.Logger.Debugf("[OrderHistoryRestService.GetOrderHistory]")
	service := service.NewOrderService(ohrs.ctx, ohrs.exchangeService)
	history := service.GetOrderHistory()
	respondWithJSON(w, http.StatusOK, history)
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
