package webservice

import (
	"encoding/json"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
)

type JsonWriter struct {
	common.HttpWriter
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{}
}

func (writer *JsonWriter) Write(w http.ResponseWriter, status int, response interface{}) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "[JsonWriter.Write] Error creating response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
