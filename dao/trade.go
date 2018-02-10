package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type TradeDAO interface {
	Create(trade entity.TradeEntity)
	Save(trade entity.TradeEntity)
	Update(trade entity.TradeEntity)
	Find(user common.User) []entity.Trade
	FindByChart(chart entity.ChartEntity) []entity.Trade
	GetLastTrade(chart entity.ChartEntity) entity.TradeEntity
}

type TradeDAOImpl struct {
	ctx *common.Context
	TradeDAO
}

func NewTradeDAO(ctx *common.Context) TradeDAO {
	ctx.DB.AutoMigrate(&entity.Trade{})
	return &TradeDAOImpl{ctx: ctx}
}

func (dao *TradeDAOImpl) Create(trade entity.TradeEntity) {
	if err := dao.ctx.DB.Create(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Create] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Save(trade entity.TradeEntity) {
	if err := dao.ctx.DB.Save(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Save] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Update(trade entity.TradeEntity) {
	if err := dao.ctx.DB.Update(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Update] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Find(user common.User) []entity.Trade {
	var trades []entity.Trade
	daoUser := &entity.User{Id: user.GetId(), Username: user.GetUsername()}
	if err := dao.ctx.DB.Model(daoUser).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.GetTrades] Error: %s", err.Error())
	}
	return trades
}
