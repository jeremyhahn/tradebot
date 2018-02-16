package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/webservice/rest"
)

type JsonWriter struct {
	common.HttpWriter
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{}
}

func (writer *JsonWriter) Write(w http.ResponseWriter, status int, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		errResponse := rest.RestResponse{Error: fmt.Sprintf("JsonWriter failed to marshal response entity %+v", reflect.TypeOf(response), response)}
		errBytes, err := json.Marshal(errResponse)
		if err != nil {
			errResponse := rest.RestResponse{Error: "JsonWriter internal server error"}
			errBytes, _ := json.Marshal(errResponse)
			http.Error(w, string(errBytes), http.StatusInternalServerError)
		}
		http.Error(w, string(errBytes), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
