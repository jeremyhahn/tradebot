package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type IndicatorDAO interface {
	Create(indicator IndicatorEntity) error
	Save(indicator IndicatorEntity) error
	Update(indicator IndicatorEntity) error
	Find() ([]Indicator, error)
	Get(name string) (IndicatorEntity, error)
}

type IndicatorDAOImpl struct {
	ctx        *common.Context
	Indicators []Indicator
	ChartDAO
}

func NewIndicatorDAO(ctx *common.Context) IndicatorDAO {
	ctx.DB.AutoMigrate(&Indicator{})
	return &IndicatorDAOImpl{ctx: ctx}
}

func (dao *IndicatorDAOImpl) Create(indicator IndicatorEntity) error {
	return dao.ctx.DB.Create(indicator).Error
}

func (dao *IndicatorDAOImpl) Save(indicator IndicatorEntity) error {
	return dao.ctx.DB.Save(indicator).Error
}

func (dao *IndicatorDAOImpl) Update(indicator IndicatorEntity) error {
	return dao.ctx.DB.Update(indicator).Error
}

func (dao *IndicatorDAOImpl) Get(name string) (IndicatorEntity, error) {
	var indicators []Indicator
	if err := dao.ctx.DB.Where("name = ?", name).Find(&indicators).Error; err != nil {
		return nil, err
	}
	if indicators == nil || len(indicators) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to get platform indicator: %s", name))
	}
	return &indicators[0], nil
}

func (dao *IndicatorDAOImpl) Find() ([]Indicator, error) {
	var indicators []Indicator
	if err := dao.ctx.DB.Order("name asc").Find(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

type IndicatorEntity interface {
	GetName() string
	GetFilename() string
	GetVersion() string
}

type Indicator struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
}

func (entity *Indicator) GetName() string {
	return entity.Name
}

func (entity *Indicator) GetFilename() string {
	return entity.Filename
}

func (entity *Indicator) GetVersion() string {
	return entity.Version
}
