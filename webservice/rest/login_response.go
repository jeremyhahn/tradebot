package rest

type LoginResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}
