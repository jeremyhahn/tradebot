package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ChartStrategyDAO interface {
	Create(indicator entity.ChartStrategyEntity) error
	Save(indicator entity.ChartStrategyEntity) error
	Update(indicator entity.ChartStrategyEntity) error
	Find(chart entity.ChartEntity) ([]entity.ChartStrategy, error)
	Get(chart entity.ChartEntity, strategyName string) (entity.ChartStrategyEntity, error)
}

type ChartStrategyDAOImpl struct {
	ctx common.Context
	ChartStrategyDAO
}

func NewChartStrategyDAO(ctx common.Context) ChartStrategyDAO {
	ctx.GetCoreDB().AutoMigrate(&entity.ChartStrategy{})
	return &ChartStrategyDAOImpl{ctx: ctx}
}

func (dao *ChartStrategyDAOImpl) Create(indicator entity.ChartStrategyEntity) error {
	return dao.ctx.GetCoreDB().Create(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Save(indicator entity.ChartStrategyEntity) error {
	return dao.ctx.GetCoreDB().Save(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Update(indicator entity.ChartStrategyEntity) error {
	return dao.ctx.GetCoreDB().Update(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Get(chart entity.ChartEntity, strategyName string) (entity.ChartStrategyEntity, error) {
	var strategies []entity.ChartStrategy
	if err := dao.ctx.GetCoreDB().Where("name = ?", strategyName).Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return &strategies[0], nil
}

func (dao *ChartStrategyDAOImpl) Find(chart entity.ChartEntity) ([]entity.ChartStrategy, error) {
	var strategies []entity.ChartStrategy
	if err := dao.ctx.GetCoreDB().Order("id asc").Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}
