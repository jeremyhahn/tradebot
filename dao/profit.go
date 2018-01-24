package dao

import "github.com/jeremyhahn/tradebot/common"

type IProfitDAO interface {
	Create(profit *Profit)
	Save(profit *Profit)
}

type ProfitDAO struct {
	ctx   *common.Context
	Items []Profit
	IProfitDAO
}

type Profit struct {
	UserID   uint
	TradeID  uint `gorm:"foreign_key"`
	Quantity float64
	Bought   float64
	Sold     float64
	Fee      float64
	Tax      float64
	Total    float64
}

func NewProfitDAO(ctx *common.Context) *ProfitDAO {
	var profits []Profit
	ctx.DB.AutoMigrate(&Profit{})
	if err := ctx.DB.Find(&profits).Error; err != nil {
		ctx.Logger.Error(err)
	}
	return &ProfitDAO{ctx: ctx, Items: profits}
}

func (dao *ProfitDAO) Create(profit *Profit) {
	if err := dao.ctx.DB.Create(profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[ProfitDAO.Create] Error:%s", err.Error())
	}
}

func (dao *ProfitDAO) Save(profit *Profit) {
	if err := dao.ctx.DB.Save(profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[ProfitDAO.Save] Error:%s", err.Error())
	}
}

func (dao *ProfitDAO) GetByTrade(trade *Trade) *Profit {
	var profit Profit
	if err := dao.ctx.DB.Model(trade).Related(&profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetTrades] Error: %s", err.Error())
	}
	return &profit
}
