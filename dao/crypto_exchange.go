package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type CryptoExchangeDAO interface {
	Create(indicator entity.PluginEntity) error
	Save(indicator entity.PluginEntity) error
	Update(indicator entity.PluginEntity) error
	Find(chart entity.ChartEntity) ([]entity.CryptoExchange, error)
	Get(chart entity.ChartEntity, indicatorName string) (entity.PluginEntity, error)
}

type CryptoExchangeDAOImpl struct {
	ctx *common.Context
}

func NewCryptoExchangeDAO(ctx *common.Context) CryptoExchangeDAO {
	return &CryptoExchangeDAOImpl{ctx: ctx}
}

func (dao *CryptoExchangeDAOImpl) Create(indicator entity.PluginEntity) error {
	return dao.ctx.CoreDB.Create(indicator).Error
}

func (dao *CryptoExchangeDAOImpl) Save(indicator entity.PluginEntity) error {
	return dao.ctx.CoreDB.Save(indicator).Error
}

func (dao *CryptoExchangeDAOImpl) Update(indicator entity.PluginEntity) error {
	return dao.ctx.CoreDB.Update(indicator).Error
}

func (dao *CryptoExchangeDAOImpl) Get(chart entity.ChartEntity, exchangeName string) (entity.PluginEntity, error) {
	var exchanges []entity.CryptoExchange
	if err := dao.ctx.CoreDB.Where("name = ?", exchangeName).Model(&entity.CryptoExchange{}).Error; err != nil {
		return nil, err
	}
	return &exchanges[0], nil
}

func (dao *CryptoExchangeDAOImpl) Find(chart entity.ChartEntity) ([]entity.CryptoExchange, error) {
	var exchanges []entity.CryptoExchange
	if err := dao.ctx.CoreDB.Order("id asc").Model(&exchanges).Error; err != nil {
		return nil, err
	}
	return exchanges, nil
}
