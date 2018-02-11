package rest

type RegisterResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}
