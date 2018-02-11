package rest

type RestResponse struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
}
