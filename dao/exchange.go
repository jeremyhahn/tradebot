package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type ExchangeDAO interface {
	Create(exchange *CryptoExchange) error
	Get(name string) (*CryptoExchange, error)
	Find() ([]CryptoExchange, error)
}

type ExchangeDAOImpl struct {
	ctx *common.Context
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
	ctx.DB.AutoMigrate(&CryptoExchange{})
	return &ExchangeDAOImpl{
		ctx: ctx}
}

func (dao *ExchangeDAOImpl) Create(exchange *CryptoExchange) error {
	if err := dao.ctx.DB.Create(exchange).Error; err != nil {
		return err
	}
	return nil
}

func (dao *ExchangeDAOImpl) Find() ([]CryptoExchange, error) {
	var exchanges []CryptoExchange
	if err := dao.ctx.DB.Find(&exchanges).Error; err != nil {
		return nil, err
	}
	return exchanges, nil
}

func (dao *ExchangeDAOImpl) Get(name string) (*CryptoExchange, error) {
	exchanges, err := dao.Find()
	if err != nil {
		return nil, err
	}
	for _, ex := range exchanges {
		if ex.Name == name {
			return &ex, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Exchange not found: %s", name))
}
