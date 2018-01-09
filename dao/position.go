package dao

import "github.com/jeremyhahn/tradebot/common"

type IPosition interface {
	Get(id uint) string
}

type PositionDAO struct {
	ctx *common.Context
	IAutoTrade
}

type Position struct {
	ID     uint `gorm:"primary_key;AUTOINCREMENT"`
	Symbol string
}

func NewPositionDAO(ctx *common.Context) *PositionDAO {
	ctx.DB.AutoMigrate(&Position{})
	return &PositionDAO{ctx: ctx}
}

func (dao *PositionDAO) Save(pos Position) {
	if err := dao.ctx.DB.Create(pos).Error; err != nil {
		dao.ctx.Logger.Errorf("[PositionDAO.Save] Error:%s", err.Error())
	}
}

func (dao *PositionDAO) Find(user *common.User) []Position {
	var positions []Position
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&positions).Error; err != nil {
		dao.ctx.Logger.Errorf("[Position.Find] Error: %s", err.Error())
	}
	return positions
}

func (dao *PositionDAO) GetCoins(user *common.User, coin *AutoTradeCoin) []Position {
	var positions []Position
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&coin).Related(&positions).Error; err != nil {
		dao.ctx.Logger.Errorf("[Position.GetCoins] Error: %s", err.Error())
	}
	return positions
}
