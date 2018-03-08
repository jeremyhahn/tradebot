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
	ctx common.Context
	OrderDAO
}

func NewOrderDAO(ctx common.Context) OrderDAO {
	return &OrderDAOImpl{ctx: ctx}
}

func (dao *OrderDAOImpl) Create(Order entity.OrderEntity) error {
	return dao.ctx.GetCoreDB().Create(Order).Error
}

func (dao *OrderDAOImpl) Save(Order entity.OrderEntity) error {
	return dao.ctx.GetCoreDB().Save(Order).Error
}

func (dao *OrderDAOImpl) Find() ([]entity.Order, error) {
	var Orders []entity.Order
	daoUser := &entity.User{Id: dao.ctx.GetUser().GetId()}
	if err := dao.ctx.GetCoreDB().Model(daoUser).Related(&Orders).Error; err != nil {
		return nil, err
	}
	return Orders, nil
}
