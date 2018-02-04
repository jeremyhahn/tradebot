package dao

import "github.com/jeremyhahn/tradebot/common"

type IndicatorDAO interface {
	Create(indicator IndicatorEntity) error
	Save(indicator IndicatorEntity) error
	Update(indicator IndicatorEntity) error
	Find(chart ChartEntity) ([]Indicator, error)
}

type IndicatorDAOImpl struct {
	ctx        *common.Context
	Indicators []Indicator
	ChartDAO
}

func NewIndicatorDAO(ctx *common.Context) IndicatorDAO {
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

func (dao *IndicatorDAOImpl) Find(chart ChartEntity) ([]Indicator, error) {
	var indicators []Indicator
	if err := dao.ctx.DB.Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

type IndicatorEntity interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
	GetFilename() string
}

type Indicator struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_indicator"`
	Name       string `gorm:"unique_index:idx_indicator"`
	Parameters string `gorm:"not null"`
	Filename   string `gorm:"not null"`
}

func (entity *Indicator) GetId() uint {
	return entity.Id
}

func (entity *Indicator) GetChartId() uint {
	return entity.ChartId
}

func (entity *Indicator) GetName() string {
	return entity.Name
}

func (entity *Indicator) GetParameters() string {
	return entity.Parameters
}

func (entity *Indicator) GetFilename() string {
	return entity.Filename
}
