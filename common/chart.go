package common

type Chart struct {
	ID         uint        `json:"id"`
	Base       string      `json:"base"`
	Quote      string      `json:"quote"`
	Exchange   string      `json:"exchange"`
	Period     int         `json:"period"`
	Price      float64     `json:"price"`
	Indicators []Indicator `json:"indicators"`
	Trades     []Trade     `json:"trades"`
}

func NewChart() *Chart {
	return &Chart{}
}

/*
func (chart *Chart) GetExchange() Exchange {
	return
}

func (chart *Chart) GetCurrencyPair() common.CurrencyPair {
	return service.Exchange.GetCurrencyPair()
}*/

func (chart *Chart) GetPrice() float64 {
	return chart.Price
}
