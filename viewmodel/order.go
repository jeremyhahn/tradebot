package viewmodel

type Order struct {
	Id       string  `json:"id"`
	Exchange string  `json:"exchange"`
	Date     string  `json:"date"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}
