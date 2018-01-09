package dao

import (
	"github.com/jeremyhahn/tradebot/common"
)

type IExchangeDAO interface {
	Get(key string) string
}

type ExchangeDAO struct {
	ctx       *common.Context
	Exchanges []CoinExchange
	IExchangeDAO
}

type CoinExchange struct {
	UserId     User
	Name       string `gorm:"primary_key" sql:"type:varchar(255)"`
	URL        string `gorm:"not null" sql:"type:varchar(255)"`
	Key        string `gorm:"not null" sql:"type:varchar(255)"`
	Secret     string `gorm:"not null" sql:"type:text"`
	Passphrase string `gorm:"not null" sql:"type:varchar(255)"`
}

func NewExchangeDAO(ctx *common.Context) *ExchangeDAO {
	var exchanges []CoinExchange
	ctx.DB.AutoMigrate(&CoinExchange{})
	if err := ctx.DB.Find(&exchanges).Error; err != nil {
		ctx.Logger.Error(err)
	}
	return &ExchangeDAO{
		ctx:       ctx,
		Exchanges: exchanges}
}

func (dao *ExchangeDAO) Create(exchange *CoinExchange) {
	if err := dao.ctx.DB.Create(exchange).Error; err != nil {
		dao.ctx.Logger.Errorf("[ExchangeDAO.Create] Error:%s", err.Error())
	}
}

func (dao *ExchangeDAO) Get(name string) *CoinExchange {
	var exchange CoinExchange
	for _, ex := range dao.Exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
