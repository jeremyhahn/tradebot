package dao

import "github.com/jeremyhahn/tradebot/common"

type IndicatorDAO interface {
	Create(chart ChartEntity)
	Save(chart ChartEntity)
	Update(chart ChartEntity)
	Find(chart ChartEntity) []Indicator
}

type IndicatorDAOImpl struct {
	ctx        *common.Context
	Indicators []Indicator
	ChartDAO
}

type Indicator struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_indicator"`
	Name       string `gorm:"unique_index:idx_indicator"`
	Parameters string `gorm:"not null"`
}

func NewIndicatorDAO(ctx *common.Context) IndicatorDAO {
	return &IndicatorDAOImpl{ctx: ctx}
}

func (dao *IndicatorDAOImpl) Create(chart ChartEntity) {
	if err := dao.ctx.DB.Create(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[IndicatorDAOImpl.Create] Error:%s", err.Error())
	}
}

func (dao *IndicatorDAOImpl) Save(chart ChartEntity) {
	if err := dao.ctx.DB.Save(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[IndicatorDAOImpl.Save] Error:%s", err.Error())
	}
}

func (dao *IndicatorDAOImpl) Update(chart ChartEntity) {
	if err := dao.ctx.DB.Update(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[IndicatorDAOImpl.Update] Error:%s", err.Error())
	}
}

func (dao *IndicatorDAOImpl) Find(chart ChartEntity) []Indicator {
	var indicators []Indicator
	if err := dao.ctx.DB.Model(chart).Related(&indicators).Error; err != nil {
		dao.ctx.Logger.Errorf("[IndicatorDAOImpl.GetIndicators] Error: %s", err.Error())
	}
	return indicators
}
