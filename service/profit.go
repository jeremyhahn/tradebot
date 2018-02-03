package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type ProfitServiceImpl struct {
	ctx       *common.Context
	profitDAO dao.ProfitDAO
	ProfitService
}

func NewProfitService(ctx *common.Context, profitDAO dao.ProfitDAO) ProfitService {
	return &ProfitServiceImpl{
		ctx:       ctx,
		profitDAO: profitDAO}
}

func (ps *ProfitServiceImpl) Save(profit *common.Profit) {
	ps.profitDAO.Create(&dao.Profit{
		UserID:   profit.UserID,
		TradeID:  profit.TradeID,
		Quantity: profit.Quantity,
		Bought:   profit.Bought,
		Sold:     profit.Sold,
		Fee:      profit.Fee,
		Tax:      profit.Tax,
		Total:    profit.Total})
}
