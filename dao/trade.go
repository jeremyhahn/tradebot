package dao

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type ITrade interface {
	Get(symbol string) string
}

type TradeDAO struct {
	ctx *common.Context
	IAutoTrade
}

type Trade struct {
	ID          uint   `gorm:"primary_key"`
	AutoTradeID uint   `gorm:"foreign_key;"`
	UserID      uint   `gorm:"index"`
	Base        string `gorm:"index"`
	Quote       string `gorm:"index"`
	Exchange    string `gorm:"index"`
	Date        time.Time
	Type        string
	Price       float64
	Amount      float64
	ChartData   string
}

func NewTradeDAO(ctx *common.Context) *TradeDAO {
	ctx.DB.AutoMigrate(&Trade{})
	return &TradeDAO{ctx: ctx}
}

func (dao *TradeDAO) Save(trade *Trade) {
	if err := dao.ctx.DB.Create(trade).Error; err != nil {
		dao.ctx.Logger.Errorf("[TradeDAO.Save] Error:%s", err.Error())
	}
}
