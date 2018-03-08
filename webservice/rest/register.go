package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type RegisterRestService interface {
	Register(w http.ResponseWriter, r *http.Request)
}

type RegisterRestServiceImpl struct {
	ctx         common.Context
	authService service.AuthService
	jsonWriter  common.HttpWriter
}

func NewRegisterRestService(ctx common.Context, authService service.AuthService, jsonWriter common.HttpWriter) RegisterRestService {
	return &RegisterRestServiceImpl{
		ctx:         ctx,
		authService: authService,
		jsonWriter:  jsonWriter}
}

func (service *RegisterRestServiceImpl) Register(w http.ResponseWriter, r *http.Request) {
	service.ctx.GetLogger().Debugf("[RegisterRestService.Register]")
	var response RegisterResponse
	var request RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		service.jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
			Success: false,
			Error:   err.Error()})
		return
	}
	err := service.authService.Register(request.Username, request.Password)
	if err != nil {
		response.Error = err.Error()
		response.Success = false
	} else {
		response.Success = true
	}
	service.jsonWriter.Write(w, http.StatusOK, response)
}
