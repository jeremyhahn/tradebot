package dao

import "github.com/jeremyhahn/tradebot/common"

type ChartStrategyDAO interface {
	Create(indicator ChartStrategyEntity) error
	Save(indicator ChartStrategyEntity) error
	Update(indicator ChartStrategyEntity) error
	Find(chart ChartEntity) ([]ChartStrategy, error)
	Get(chart ChartEntity, strategyName string) (ChartStrategyEntity, error)
}

type ChartStrategyDAOImpl struct {
	ctx            *common.Context
	ChartStrategys []ChartStrategy
	ChartDAO
}

func NewChartStrategyDAO(ctx *common.Context) ChartStrategyDAO {
	ctx.DB.AutoMigrate(&ChartStrategy{})
	return &ChartStrategyDAOImpl{ctx: ctx}
}

func (dao *ChartStrategyDAOImpl) Create(indicator ChartStrategyEntity) error {
	return dao.ctx.DB.Create(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Save(indicator ChartStrategyEntity) error {
	return dao.ctx.DB.Save(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Update(indicator ChartStrategyEntity) error {
	return dao.ctx.DB.Update(indicator).Error
}

func (dao *ChartStrategyDAOImpl) Get(chart ChartEntity, strategyName string) (ChartStrategyEntity, error) {
	var strategies []ChartStrategy
	if err := dao.ctx.DB.Where("name = ?", strategyName).Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return &strategies[0], nil
}

func (dao *ChartStrategyDAOImpl) Find(chart ChartEntity) ([]ChartStrategy, error) {
	var strategies []ChartStrategy
	if err := dao.ctx.DB.Order("id asc").Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}

type ChartStrategyEntity interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
}

type ChartStrategy struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_chart_strategy"`
	Name       string `gorm:"unique_index:idx_chart_strategy"`
	Parameters string `gorm:"not null"`
}

func (entity *ChartStrategy) GetId() uint {
	return entity.Id
}

func (entity *ChartStrategy) GetChartId() uint {
	return entity.ChartId
}

func (entity *ChartStrategy) GetName() string {
	return entity.Name
}

func (entity *ChartStrategy) GetParameters() string {
	return entity.Parameters
}
