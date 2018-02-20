package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type StrategyDAO interface {
	Create(indicator entity.StrategyEntity) error
	Save(indicator entity.StrategyEntity) error
	Update(indicator entity.StrategyEntity) error
	Find() ([]entity.Strategy, error)
	Get(name string) (entity.StrategyEntity, error)
}

type StrategyDAOImpl struct {
	ctx       *common.Context
	Strategys []entity.Strategy
	ChartDAO
}

func NewStrategyDAO(ctx *common.Context) StrategyDAO {
	ctx.CoreDB.AutoMigrate(&entity.Strategy{})
	return &StrategyDAOImpl{ctx: ctx}
}

func (dao *StrategyDAOImpl) Create(indicator entity.StrategyEntity) error {
	return dao.ctx.CoreDB.Create(indicator).Error
}

func (dao *StrategyDAOImpl) Save(indicator entity.StrategyEntity) error {
	return dao.ctx.CoreDB.Save(indicator).Error
}

func (dao *StrategyDAOImpl) Update(indicator entity.StrategyEntity) error {
	return dao.ctx.CoreDB.Update(indicator).Error
}

func (dao *StrategyDAOImpl) Get(name string) (entity.StrategyEntity, error) {
	var strategies []entity.Strategy
	if err := dao.ctx.CoreDB.Where("name = ?", name).Find(&strategies).Error; err != nil {
		return nil, err
	}
	if strategies == nil || len(strategies) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to get platform strategy: %s", name))
	}
	return &strategies[0], nil
}

func (dao *StrategyDAOImpl) Find() ([]entity.Strategy, error) {
	var strategies []entity.Strategy
	if err := dao.ctx.CoreDB.Order("name asc").Find(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}
