package dto

import "github.com/jeremyhahn/tradebot/common"

type ProfitDTO struct {
	UserId   uint    `json:"id"`
	TradeId  uint    `json:"trade_id"`
	Quantity float64 `json:"quantity"`
	Bought   float64 `json:"bought"`
	Sold     float64 `json:"sold"`
	Fee      float64 `json:"fee"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
	common.Profit
}

func NewProfitDTO() common.Profit {
	return &ProfitDTO{}
}

func (dto *ProfitDTO) GetUserId() uint {
	return dto.UserId
}

func (dto *ProfitDTO) GetTradeId() uint {
	return dto.TradeId
}

func (dto *ProfitDTO) GetQuantity() float64 {
	return dto.Quantity
}

func (dto *ProfitDTO) GetBought() float64 {
	return dto.Bought
}

func (dto *ProfitDTO) GetSold() float64 {
	return dto.Sold
}

func (dto *ProfitDTO) GetFee() float64 {
	return dto.Fee
}

func (dto *ProfitDTO) GetTax() float64 {
	return dto.Tax
}

func (dto *ProfitDTO) GetTotal() float64 {
	return dto.Total
}
