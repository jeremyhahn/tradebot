package dao

import "github.com/jeremyhahn/tradebot/common"

type IAutoTrade interface {
	Get(symbol string) string
}

type AutoTradeDAO struct {
	ctx   *common.Context
	Coins []AutoTradeCoin
	IAutoTrade
}

type AutoTradeCoin struct {
	UserID    uint
	Base      string
	Quote     string
	Exchange  string
	Period    int
	Positions []Position
}

func NewAutoTradeDAO(ctx *common.Context) *AutoTradeDAO {
	ctx.DB.AutoMigrate(&AutoTradeCoin{})
	return &AutoTradeDAO{ctx: ctx}
}

func (dao *AutoTradeDAO) Save(coin AutoTradeCoin) {
	if err := dao.ctx.DB.Create(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Save] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Find(user *common.User) []AutoTradeCoin {
	var coins []AutoTradeCoin
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&coins).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Find] Error: %s", err.Error())
	}
	return coins
}
