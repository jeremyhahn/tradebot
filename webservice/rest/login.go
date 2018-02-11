package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type LoginRestService interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type LoginRestServiceImpl struct {
	ctx         *common.Context
	authService service.AuthService
}

func NewLoginRestService(ctx *common.Context, authService service.AuthService) LoginRestService {
	return &LoginRestServiceImpl{
		ctx:         ctx,
		authService: authService}
}

func (service *LoginRestServiceImpl) Login(w http.ResponseWriter, r *http.Request) {
	service.ctx.Logger.Debugf("[LoginRestService.Login]")
	var loginRequest LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		respondWithJSON(w, http.StatusBadRequest, LoginResponse{
			Error: err.Error()})
		return
	}
	err := service.authService.Login(loginRequest.Password)
	var response LoginResponse
	if err != nil {
		response.Error = err.Error()
		response.Success = false
	} else {
		response.Success = true
	}
	respondWithJSON(w, http.StatusOK, response)
}
