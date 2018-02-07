package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type StrategyDAO interface {
	Create(indicator StrategyEntity) error
	Save(indicator StrategyEntity) error
	Update(indicator StrategyEntity) error
	Find() ([]Strategy, error)
	Get(name string) (StrategyEntity, error)
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

func (dao *StrategyDAOImpl) Get(name string) (StrategyEntity, error) {
	var strategies []Strategy
	if err := dao.ctx.DB.Where("name = ?", name).Find(&strategies).Error; err != nil {
		return nil, err
	}
	if strategies == nil || len(strategies) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to get platform strategy: %s", name))
	}
	return &strategies[0], nil
}

func (dao *StrategyDAOImpl) Find() ([]Strategy, error) {
	var strategies []Strategy
	if err := dao.ctx.DB.Order("name asc").Find(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}

type StrategyEntity interface {
	GetName() string
	GetFilename() string
	GetVersion() string
}

type Strategy struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
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
