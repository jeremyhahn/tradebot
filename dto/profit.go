package dto

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type ProfitDTO struct {
	UserId   uint            `json:"id"`
	TradeId  uint            `json:"trade_id"`
	Quantity decimal.Decimal `json:"quantity"`
	Bought   decimal.Decimal `json:"bought"`
	Sold     decimal.Decimal `json:"sold"`
	Fee      decimal.Decimal `json:"fee"`
	Tax      decimal.Decimal `json:"tax"`
	Total    decimal.Decimal `json:"total"`
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

func (dto *ProfitDTO) GetQuantity() decimal.Decimal {
	return dto.Quantity
}

func (dto *ProfitDTO) GetBought() decimal.Decimal {
	return dto.Bought
}

func (dto *ProfitDTO) GetSold() decimal.Decimal {
	return dto.Sold
}

func (dto *ProfitDTO) GetFee() decimal.Decimal {
	return dto.Fee
}

func (dto *ProfitDTO) GetTax() decimal.Decimal {
	return dto.Tax
}

func (dto *ProfitDTO) GetTotal() decimal.Decimal {
	return dto.Total
}
