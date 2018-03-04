package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ChartIndicatorDAO interface {
	Create(indicator entity.ChartIndicatorEntity) error
	Save(indicator entity.ChartIndicatorEntity) error
	Update(indicator entity.ChartIndicatorEntity) error
	Find(chart entity.ChartEntity) ([]entity.ChartIndicator, error)
	Get(chart entity.ChartEntity, indicatorName string) (entity.ChartIndicatorEntity, error)
}

type ChartIndicatorDAOImpl struct {
	ctx common.Context
	ChartIndicatorDAO
}

func NewChartIndicatorDAO(ctx common.Context) ChartIndicatorDAO {
	return &ChartIndicatorDAOImpl{ctx: ctx}
}

func (dao *ChartIndicatorDAOImpl) Create(indicator entity.ChartIndicatorEntity) error {
	return dao.ctx.GetCoreDB().Create(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Save(indicator entity.ChartIndicatorEntity) error {
	return dao.ctx.GetCoreDB().Save(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Update(indicator entity.ChartIndicatorEntity) error {
	return dao.ctx.GetCoreDB().Update(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Get(chart entity.ChartEntity, indicatorName string) (entity.ChartIndicatorEntity, error) {
	var indicators []entity.ChartIndicator
	if err := dao.ctx.GetCoreDB().Where("name = ?", indicatorName).Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return &indicators[0], nil
}

func (dao *ChartIndicatorDAOImpl) Find(chart entity.ChartEntity) ([]entity.ChartIndicator, error) {
	var indicators []entity.ChartIndicator
	if err := dao.ctx.GetCoreDB().Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}
