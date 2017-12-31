package dao

import (
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type Configuration interface {
	Get(key string) string
}

type ConfigurationDAO struct {
	Items []Config
	Configuration
}

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

func NewConfigurationDAO(db *gorm.DB, logger *logging.Logger) *ConfigurationDAO {
	var configs []Config
	db.AutoMigrate(&Config{})
	if err := db.Find(&configs).Error; err != nil {
		logger.Error(err)
	}
	return &ConfigurationDAO{Items: configs}
}

func (dao *ConfigurationDAO) Get(key string) string {
	value := ""
	for _, config := range dao.Items {
		if config.Key == key {
			return config.Value
		}
	}
	return value
}
