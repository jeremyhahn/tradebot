package dao

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type TradeDAO interface {
	Create(trade TradeEntity)
	Save(trade TradeEntity)
	Update(trade TradeEntity)
	Find(user *common.User) []Trade
	FindByChart(chart ChartEntity) []Trade
	GetLastTrade(chart ChartEntity) *Trade
	TradeEntity
}

type TradeDAOImpl struct {
	ctx *common.Context
	TradeDAO
}

type TradeEntity interface {
	GetId() uint
	GetChartId() uint
	GetUserId() uint
	GetBase() string
	GetQuote() string
	GetExchangeName() string
	GetDate() time.Time
	GetType() string
	GetPrice() float64
	GetAmount() float64
	GetChartData() string
}

type Trade struct {
	Id        uint   `gorm:"primary_key"`
	ChartId   uint   `gorm:"foreign_key"`
	UserId    uint   `gorm:"foreign_key;index"`
	Base      string `gorm:"index"`
	Quote     string `gorm:"index"`
	Exchange  string `gorm:"index"`
	Date      time.Time
	Type      string
	Price     float64
	Amount    float64
	ChartData string
}

func NewTradeDAO(ctx *common.Context) TradeDAO {
	ctx.DB.AutoMigrate(&Trade{})
	return &TradeDAOImpl{ctx: ctx}
}

func (dao *TradeDAOImpl) Create(trade TradeEntity) {
	if err := dao.ctx.DB.Create(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Create] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Save(trade TradeEntity) {
	if err := dao.ctx.DB.Save(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Save] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Update(trade TradeEntity) {
	if err := dao.ctx.DB.Update(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.Update] Error:%s", err.Error())
	}
}

func (dao *TradeDAOImpl) Find(user *common.User) []Trade {
	var trades []Trade
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.GetTrades] Error: %s", err.Error())
	}
	return trades
}

func (trade *Trade) GetId() uint {
	return trade.Id
}

func (trade *Trade) GetChartId() uint {
	return trade.ChartId
}

func (trade *Trade) GetUserId() uint {
	return trade.UserId
}

func (trade *Trade) GetBase() string {
	return trade.Base
}

func (trade *Trade) GetQuote() string {
	return trade.Quote
}

func (trade *Trade) GetExchangeName() string {
	return trade.Exchange
}

func (trade *Trade) GetDate() time.Time {
	return trade.Date
}

func (trade *Trade) GetType() string {
	return trade.Type
}

func (trade *Trade) GetPrice() float64 {
	return trade.Price
}

func (trade *Trade) GetAmount() float64 {
	return trade.Amount
}

func (trade *Trade) GetChartData() string {
	return trade.ChartData
}
