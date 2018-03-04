package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type PriceHistoryDAO interface {
}

type PriceHistoryDAOImpl struct {
	ctx common.Context
}

func NewPriceHistoryDAO(ctx common.Context) PriceHistoryDAO {
	return &PriceHistoryDAOImpl{
		ctx: ctx}
}

func (phDAO *PriceHistoryDAOImpl) Create(priceHistory entity.PriceHistoryEntity) error {
	return phDAO.ctx.GetPriceDB().Create(priceHistory).Error
}

func (phDAO *PriceHistoryDAOImpl) Save(priceHistory entity.PriceHistoryEntity) error {
	return phDAO.ctx.GetPriceDB().Save(priceHistory).Error
}

func (phDAO *PriceHistoryDAOImpl) Update(priceHistory entity.PriceHistoryEntity) error {
	return phDAO.ctx.GetPriceDB().Update(priceHistory).Error
}
