package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type IndicatorDAO interface {
	Create(indicator entity.IndicatorEntity) error
	Save(indicator entity.IndicatorEntity) error
	Update(indicator entity.IndicatorEntity) error
	Find() ([]entity.Indicator, error)
	Get(name string) (entity.IndicatorEntity, error)
}

type IndicatorDAOImpl struct {
	ctx        *common.Context
	Indicators []entity.Indicator
	ChartDAO
}

func NewIndicatorDAO(ctx *common.Context) IndicatorDAO {
	ctx.CoreDB.AutoMigrate(&entity.Indicator{})
	return &IndicatorDAOImpl{ctx: ctx}
}

func (dao *IndicatorDAOImpl) Create(indicator entity.IndicatorEntity) error {
	return dao.ctx.CoreDB.Create(indicator).Error
}

func (dao *IndicatorDAOImpl) Save(indicator entity.IndicatorEntity) error {
	return dao.ctx.CoreDB.Save(indicator).Error
}

func (dao *IndicatorDAOImpl) Update(indicator entity.IndicatorEntity) error {
	return dao.ctx.CoreDB.Update(indicator).Error
}

func (dao *IndicatorDAOImpl) Get(name string) (entity.IndicatorEntity, error) {
	var indicators []entity.Indicator
	if err := dao.ctx.CoreDB.Where("name = ?", name).Find(&indicators).Error; err != nil {
		return nil, err
	}
	if indicators == nil || len(indicators) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to get platform indicator: %s", name))
	}
	return &indicators[0], nil
}

func (dao *IndicatorDAOImpl) Find() ([]entity.Indicator, error) {
	var indicators []entity.Indicator
	if err := dao.ctx.CoreDB.Order("name asc").Find(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}
