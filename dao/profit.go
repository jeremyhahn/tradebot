package dao

import "github.com/jeremyhahn/tradebot/common"

type ProfitDAO interface {
	Create(profit *Profit)
	Save(profit *Profit)
}

type ProfitDAOImpl struct {
	ctx   *common.Context
	Items []Profit
	ProfitDAO
}

type Profit struct {
	ID       uint `gorm:"primary_key"`
	UserID   uint `gorm:"unique_index:idx_profit"`
	TradeID  uint `gorm:"foreign_key;unique_index:idx_profit"`
	Quantity float64
	Bought   float64
	Sold     float64
	Fee      float64
	Tax      float64
	Total    float64
}

func NewProfitDAO(ctx *common.Context) ProfitDAO {
	var profits []Profit
	ctx.DB.AutoMigrate(&Profit{})
	if err := ctx.DB.Find(&profits).Error; err != nil {
		ctx.Logger.Error(err)
	}
	return &ProfitDAOImpl{ctx: ctx, Items: profits}
}

func (dao *ProfitDAOImpl) Create(profit *Profit) {
	if err := dao.ctx.DB.Create(profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[ProfitDAOImpl.Create] Error:%s", err.Error())
	}
}

func (dao *ProfitDAOImpl) Save(profit *Profit) {
	if err := dao.ctx.DB.Save(profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[ProfitDAOImpl.Save] Error:%s", err.Error())
	}
}

func (dao *ProfitDAOImpl) GetByTrade(trade *Trade) *Profit {
	var profit Profit
	if err := dao.ctx.DB.Model(trade).Related(&profit).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetTrades] Error: %s", err.Error())
	}
	return &profit
}
