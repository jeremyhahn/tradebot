package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type DefaultProfitService struct {
	ctx       *common.Context
	profitDAO dao.ProfitDAO
	ProfitService
}

func NewProfitService(ctx *common.Context, profitDAO dao.ProfitDAO) ProfitService {
	return &DefaultProfitService{
		ctx:       ctx,
		profitDAO: profitDAO}
}

func (ps *DefaultProfitService) Save(profit common.Profit) {
	ps.profitDAO.Create(&dao.Profit{
		UserId:   profit.GetUserId(),
		TradeId:  profit.GetTradeId(),
		Quantity: profit.GetQuantity(),
		Bought:   profit.GetBought(),
		Sold:     profit.GetSold(),
		Fee:      profit.GetFee(),
		Tax:      profit.GetTax(),
		Total:    profit.GetTotal()})
}
