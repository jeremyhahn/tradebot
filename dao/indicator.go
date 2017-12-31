package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type Indicator struct {
	Name       string `gorm:"primary_key"`
	Parameters string `gorm:"not null"`
}

type IndicatorDAO struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewIndicatorDAO(db *gorm.DB, logger *logging.Logger) *IndicatorDAO {
	db.AutoMigrate(&Indicator{})
	return &IndicatorDAO{
		db:     db,
		logger: logger}
}

func (dao *IndicatorDAO) Find(indicator string) []common.Indicator {
	var indicators []common.Indicator
	if err := dao.db.Find(&indicators).Error; err != nil {
		dao.logger.Error(err)
	}
	return indicators
}

func (dao *IndicatorDAO) Get(id string) common.Indicator {
	var indicator common.Indicator
	if err := dao.db.First(&id).Error; err != nil {
		dao.logger.Error(err)
	}
	return indicator
}

func (dao *IndicatorDAO) Delete(id string) common.Indicator {
	var indicator common.Indicator
	if err := dao.db.Delete(&Indicator{Name: id}).Error; err != nil {
		dao.logger.Error(err)
	}
	return indicator
}
