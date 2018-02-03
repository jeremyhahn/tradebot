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
	ctx    *common.Context
	Trades []Trade
	TradeDAO
}

type TradeEntity interface {
	GetId() uint
	GetBase() string
	GetQuote() string
	GetExchangeName() string
	GetDate() time.Time
	GetType() string
	GetPrice() float64
	GetAmount() float64
}

type Trade struct {
	ID        uint   `gorm:"primary_key"`
	ChartID   uint   `gorm:"foreign_key"`
	UserID    uint   `gorm:"foreign_key;index"`
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

func (dao *TradeDAOImpl) GetLastTrade(chart ChartEntity) *Trade {
	var trades []Trade
	if err := dao.ctx.DB.Order("date desc").Limit(1).Model(chart).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		dao.ctx.Logger.Warningf("[TradeDAOImpl.GetLastTrade] Invalid number of trades returned: %d", tradeLen)
		return &Trade{}
	}
	return &trades[0]
}

func (dao *TradeDAOImpl) FindByChart(chart ChartEntity) []Trade {
	var trades []Trade
	daoChart := &Chart{Id: chart.GetId()}
	if err := dao.ctx.DB.Model(daoChart).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAOImpl.GetTrades] Error: %s", err.Error())
	}
	return trades
}

func (trade *Trade) GetId() uint {
	return trade.ID
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
