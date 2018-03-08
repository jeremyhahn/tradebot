package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type JsonResponse struct {
	Error   string      `json:"error"`
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

type JsonWriter struct {
	HttpWriter
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{}
}

func (writer *JsonWriter) Write(w http.ResponseWriter, status int, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		errResponse := JsonResponse{Error: fmt.Sprintf("JsonWriter failed to marshal response entity %s %+v", reflect.TypeOf(response), response)}
		errBytes, err := json.Marshal(errResponse)
		if err != nil {
			errResponse := JsonResponse{Error: "JsonWriter internal server error"}
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
