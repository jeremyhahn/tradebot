package rest

type RestResponse struct {
	Error   string      `json:"error"`
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}
