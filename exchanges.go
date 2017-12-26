package main

import (
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type ICoinExchanges interface {
	Get(key string) string
}

type CoinExchanges struct {
	Exchanges []CoinExchange
	ICoinExchanges
}

type CoinExchange struct {
	Name       string `gorm:"primary_key" sql:"type:varchar(255)"`
	URL        string `gorm:"not null" sql:"type:varchar(255)"`
	Key        string `gorm:"not null" sql:"type:varchar(255)"`
	Secret     string `gorm:"not null" sql:"type:text"`
	Passphrase string `gorm:"not null" sql:"type:varchar(255)"`
}

func NewCoinExchanges(db *gorm.DB, logger *logging.Logger) *CoinExchanges {
	var exchanges []CoinExchange
	db.AutoMigrate(&CoinExchange{})
	if err := db.Find(&exchanges).Error; err != nil {
		logger.Error(err)
	}
	return &CoinExchanges{Exchanges: exchanges}
}

func (ce *CoinExchanges) Get(name string) *CoinExchange {
	var exchange CoinExchange
	for _, ex := range ce.Exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
