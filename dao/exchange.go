package dao

import "github.com/jeremyhahn/tradebot/common"

type ExchangeDAO interface {
	Get(key string) *CryptoExchange
}

type ExchangeDAOImpl struct {
	ctx       *common.Context
	Exchanges []CryptoExchange
	ExchangeDAO
}

type CryptoExchange struct {
	UserID uint
	Name   string `gorm:"primary_key" sql:"type:varchar(255)"`
	URL    string `gorm:"not null" sql:"type:varchar(255)"`
	Key    string `gorm:"not null" sql:"type:varchar(255)"`
	Secret string `gorm:"not null" sql:"type:text"`
	Extra  string `gorm:"not null" sql:"type:varchar(255)"`
}

func NewExchangeDAO(ctx *common.Context) ExchangeDAO {
	var exchanges []CryptoExchange
	ctx.DB.AutoMigrate(&CryptoExchange{})
	if err := ctx.DB.Find(&exchanges).Error; err != nil {
		ctx.Logger.Error(err)
	}
	return &ExchangeDAOImpl{
		ctx:       ctx,
		Exchanges: exchanges}
}

func (dao *ExchangeDAOImpl) Create(exchange *CryptoExchange) {
	if err := dao.ctx.DB.Create(exchange).Error; err != nil {
		dao.ctx.Logger.Errorf("[ExchangeDAO.Create] Error:%s", err.Error())
	}
}

func (dao *ExchangeDAOImpl) Get(name string) *CryptoExchange {
	var exchange CryptoExchange
	for _, ex := range dao.Exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
