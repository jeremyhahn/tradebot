package rest

import (
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
)

func RestError(w http.ResponseWriter, r *http.Request, err error, jsonWriter common.HttpWriter) {
	jsonWriter.Write(w, http.StatusBadRequest, common.JsonResponse{
		Success: false,
		Payload: err.Error()})
}
