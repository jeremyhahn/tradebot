package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type OrderDAO interface {
	Create(Order entity.OrderEntity) error
	Save(Order entity.OrderEntity) error
	Find() ([]entity.Order, error)
	GetByTrade(trade entity.TradeEntity) (entity.OrderEntity, error)
}

type OrderDAOImpl struct {
	ctx   *common.Context
	Items []entity.Order
	OrderDAO
}

func NewOrderDAO(ctx *common.Context) OrderDAO {
	ctx.CoreDB.AutoMigrate(&entity.Order{})
	return &OrderDAOImpl{ctx: ctx}
}

func (dao *OrderDAOImpl) Create(Order entity.OrderEntity) error {
	return dao.ctx.CoreDB.Create(Order).Error
}

func (dao *OrderDAOImpl) Save(Order entity.OrderEntity) error {
	return dao.ctx.CoreDB.Save(Order).Error
}

func (dao *OrderDAOImpl) Find() ([]entity.Order, error) {
	var Orders []entity.Order
	daoUser := &entity.User{Id: dao.ctx.User.GetId()}
	if err := dao.ctx.CoreDB.Model(daoUser).Related(&Orders).Error; err != nil {
		return nil, err
	}
	return Orders, nil
}
