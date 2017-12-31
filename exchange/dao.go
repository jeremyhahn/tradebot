package exchange

import (
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type IExchangeDAO interface {
	Get(key string) string
}

type ExchangeDAO struct {
	Exchanges []CoinExchange
	IExchangeDAO
}

type CoinExchange struct {
	Name       string `gorm:"primary_key" sql:"type:varchar(255)"`
	URL        string `gorm:"not null" sql:"type:varchar(255)"`
	Key        string `gorm:"not null" sql:"type:varchar(255)"`
	Secret     string `gorm:"not null" sql:"type:text"`
	Passphrase string `gorm:"not null" sql:"type:varchar(255)"`
}

func NewExchangeDAO(db *gorm.DB, logger *logging.Logger) *ExchangeDAO {
	var exchanges []CoinExchange
	db.AutoMigrate(&CoinExchange{})
	if err := db.Find(&exchanges).Error; err != nil {
		logger.Error(err)
	}
	return &ExchangeDAO{Exchanges: exchanges}
}

func (el *ExchangeDAO) Get(name string) *CoinExchange {
	var exchange CoinExchange
	for _, ex := range el.Exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
