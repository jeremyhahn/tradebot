package main

import (
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type IConfiguration interface {
	Get(key string) string
}

type Configuration struct {
	Items []Config
	IConfiguration
}

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

func NewConfiguration(db *gorm.DB, logger *logging.Logger) *Configuration {
	var configs []Config
	db.AutoMigrate(&Config{})
	if err := db.Find(&configs).Error; err != nil {
		logger.Error(err)
	}
	return &Configuration{Items: configs}
}

func (config *Configuration) Get(key string) string {
	value := ""
	for _, config := range config.Items {
		if config.Key == key {
			return config.Value
		}
	}
	return value
}
