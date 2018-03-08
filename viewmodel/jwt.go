package viewmodel

type JsonWebToken struct {
	Value string `json:"token"`
	Error string `json:"error"`
}
