package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ProfitDAO interface {
	Create(profit entity.ProfitEntity) error
	Save(profit entity.ProfitEntity) error
	Find() ([]entity.Profit, error)
	GetByTrade(trade entity.TradeEntity) (entity.ProfitEntity, error)
}

type ProfitDAOImpl struct {
	ctx common.Context
	ProfitDAO
}

func NewProfitDAO(ctx common.Context) ProfitDAO {
	return &ProfitDAOImpl{ctx: ctx}
}

func (dao *ProfitDAOImpl) Create(profit entity.ProfitEntity) error {
	return dao.ctx.GetCoreDB().Create(profit).Error
}

func (dao *ProfitDAOImpl) Save(profit entity.ProfitEntity) error {
	return dao.ctx.GetCoreDB().Save(profit).Error
}

func (dao *ProfitDAOImpl) Find() ([]entity.Profit, error) {
	var profits []entity.Profit
	daoUser := &entity.User{Id: dao.ctx.GetUser().GetId()}
	if err := dao.ctx.GetCoreDB().Model(daoUser).Related(&profits).Error; err != nil {
		return nil, err
	}
	return profits, nil
}

func (dao *ProfitDAOImpl) GetByTrade(trade entity.TradeEntity) (entity.ProfitEntity, error) {
	var profit entity.Profit
	if err := dao.ctx.GetCoreDB().Model(trade).Related(&profit).Error; err != nil {
		return nil, err
	}
	return &profit, nil
}
