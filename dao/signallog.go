package dao

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type ISignalLogDAO interface {
	Save(log *SignalLog)
	Update(log *SignalLog)
	Find(user *common.User) []SignalLog
}

type SignalLogDAO struct {
	ctx *common.Context
	ISignalLogDAO
}

type SignalLog struct {
	ID         uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID     uint
	Date       time.Time
	Name       string
	Type       string
	Price      float64
	SignalData string
}

func NewSignalLogDAO(ctx *common.Context) *SignalLogDAO {
	ctx.DB.AutoMigrate(&SignalLog{})
	return &SignalLogDAO{ctx: ctx}
}

func (dao *SignalLogDAO) Save(log *SignalLog) {
	if err := dao.ctx.DB.Create(log).Error; err != nil {
		dao.ctx.Logger.Errorf("[SignalLogDAO.Save] Error:%s", err.Error())
	}
}

func (dao *SignalLogDAO) Update(log *SignalLog) {
	if err := dao.ctx.DB.Update(log).Error; err != nil {
		dao.ctx.Logger.Errorf("[SignalLogDAO.Update] Error:%s", err.Error())
	}
}

func (dao *SignalLogDAO) Find(user *common.User) []SignalLog {
	var entries []SignalLog
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&entries).Error; err != nil {
		dao.ctx.Logger.Errorf("[SignalLogDAO.Find] Error: %s", err.Error())
	}
	return entries
}
