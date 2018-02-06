package dao

import (
	"github.com/jeremyhahn/tradebot/common"
)

type ChartIndicatorDAO interface {
	Create(indicator ChartIndicatorEntity) error
	Save(indicator ChartIndicatorEntity) error
	Update(indicator ChartIndicatorEntity) error
	Find(chart ChartEntity) ([]ChartIndicator, error)
	Get(chart ChartEntity, indicatorName string) (ChartIndicatorEntity, error)
}

type ChartIndicatorDAOImpl struct {
	ctx             *common.Context
	ChartIndicators []ChartIndicator
	ChartDAO
}

func NewChartIndicatorDAO(ctx *common.Context) ChartIndicatorDAO {
	ctx.DB.AutoMigrate(&ChartIndicator{})
	return &ChartIndicatorDAOImpl{ctx: ctx}
}

func (dao *ChartIndicatorDAOImpl) Create(indicator ChartIndicatorEntity) error {
	return dao.ctx.DB.Create(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Save(indicator ChartIndicatorEntity) error {
	return dao.ctx.DB.Save(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Update(indicator ChartIndicatorEntity) error {
	return dao.ctx.DB.Update(indicator).Error
}

func (dao *ChartIndicatorDAOImpl) Get(chart ChartEntity, indicatorName string) (ChartIndicatorEntity, error) {
	var indicators []ChartIndicator
	if err := dao.ctx.DB.Where("name = ?", indicatorName).Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return &indicators[0], nil
}

func (dao *ChartIndicatorDAOImpl) Find(chart ChartEntity) ([]ChartIndicator, error) {
	var indicators []ChartIndicator
	if err := dao.ctx.DB.Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

type ChartIndicatorEntity interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
}

type ChartIndicator struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_chart_indicator"`
	Name       string `gorm:"unique_index:idx_chart_indicator"`
	Parameters string `gorm:"not null"`
}

func (entity *ChartIndicator) GetId() uint {
	return entity.Id
}

func (entity *ChartIndicator) GetChartId() uint {
	return entity.ChartId
}

func (entity *ChartIndicator) GetName() string {
	return entity.Name
}

func (entity *ChartIndicator) GetParameters() string {
	return entity.Parameters
}
