package dao

import "github.com/jeremyhahn/tradebot/common"

type StrategyDAO interface {
	Create(indicator StrategyEntity) error
	Save(indicator StrategyEntity) error
	Update(indicator StrategyEntity) error
	Find() ([]Strategy, error)
}

type StrategyDAOImpl struct {
	ctx       *common.Context
	Strategys []Strategy
	ChartDAO
}

func NewStrategyDAO(ctx *common.Context) StrategyDAO {
	ctx.DB.AutoMigrate(&Strategy{})
	return &StrategyDAOImpl{ctx: ctx}
}

func (dao *StrategyDAOImpl) Create(indicator StrategyEntity) error {
	return dao.ctx.DB.Create(indicator).Error
}

func (dao *StrategyDAOImpl) Save(indicator StrategyEntity) error {
	return dao.ctx.DB.Save(indicator).Error
}

func (dao *StrategyDAOImpl) Update(indicator StrategyEntity) error {
	return dao.ctx.DB.Update(indicator).Error
}

func (dao *StrategyDAOImpl) Find() ([]Strategy, error) {
	var indicators []Strategy
	if err := dao.ctx.DB.Order("id asc").Find(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

type StrategyEntity interface {
	GetId() uint
	GetName() string
	GetFilename() string
	GetVersion() string
}

type Strategy struct {
	Id       uint   `gorm:"primary_key"`
	Name     string `gorm:"unique_index:idx_indicator"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
}

func (entity *Strategy) GetId() uint {
	return entity.Id
}

func (entity *Strategy) GetName() string {
	return entity.Name
}

func (entity *Strategy) GetFilename() string {
	return entity.Filename
}

func (entity *Strategy) GetVersion() string {
	return entity.Version
}
