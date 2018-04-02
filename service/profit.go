package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
)

type DefaultProfitService struct {
	ctx       common.Context
	profitDAO dao.ProfitDAO
	ProfitService
}

func NewProfitService(ctx common.Context, profitDAO dao.ProfitDAO) ProfitService {
	return &DefaultProfitService{
		ctx:       ctx,
		profitDAO: profitDAO}
}

func (ps *DefaultProfitService) Save(profit common.Profit) {
	ps.profitDAO.Create(&entity.Profit{
		UserId:   profit.GetUserId(),
		TradeId:  profit.GetTradeId(),
		Quantity: profit.GetQuantity().String(),
		Bought:   profit.GetBought().String(),
		Sold:     profit.GetSold().String(),
		Fee:      profit.GetFee().String(),
		Tax:      profit.GetTax().String(),
		Total:    profit.GetTotal().String()})
}
