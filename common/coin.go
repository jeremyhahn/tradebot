package common

type ICoin interface {
	IsBitcoin()
}

type Coin struct {
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	Available float64 `json:"available"`
	Pending   float64 `json:"pending"`
	Price     float64 `json:"price"`
	Address   string  `json:"address"`
	Total     float64 `json:"total"`
	BTC       float64 `json:"btc"`
}

func NewCoin() *Coin {
	return &Coin{}
}

func (c *Coin) IsBitcoin() bool {
	return c.Currency == "BTC"
}
