package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ExchangeDAO interface {
	Create(exchange *entity.CryptoExchange) error
	Get(name string) (*entity.CryptoExchange, error)
	Find() ([]entity.CryptoExchange, error)
}

type ExchangeDAOImpl struct {
	ctx *common.Context
	ExchangeDAO
}

func NewExchangeDAO(ctx *common.Context) ExchangeDAO {
	ctx.CoreDB.AutoMigrate(&entity.CryptoExchange{})
	return &ExchangeDAOImpl{
		ctx: ctx}
}

func (dao *ExchangeDAOImpl) Create(exchange *entity.CryptoExchange) error {
	return dao.ctx.CoreDB.Create(exchange).Error
}

func (dao *ExchangeDAOImpl) Find() ([]entity.CryptoExchange, error) {
	var exchanges []entity.CryptoExchange
	if err := dao.ctx.CoreDB.Find(&exchanges).Error; err != nil {
		return nil, err
	}
	return exchanges, nil
}

func (dao *ExchangeDAOImpl) Get(name string) (*entity.CryptoExchange, error) {
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
