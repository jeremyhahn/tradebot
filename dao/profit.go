package dao

import "github.com/jeremyhahn/tradebot/common"

type ProfitDAO interface {
	Create(profit ProfitEntity) error
	Save(profit ProfitEntity) error
	Find() ([]Profit, error)
	GetByTrade(trade TradeEntity) (ProfitEntity, error)
}

type ProfitDAOImpl struct {
	ctx   *common.Context
	Items []Profit
	ProfitDAO
}

type ProfitEntity interface {
	GetId() uint
	GetUserId() uint
	GetTradeId() uint
	GetQuantity() float64
	GetBought() float64
	GetSold() float64
	GetFee() float64
	GetTax() float64
	GetTotal() float64
}

type Profit struct {
	Id       uint `gorm:"primary_key"`
	UserId   uint `gorm:"unique_index:idx_profit"`
	TradeId  uint `gorm:"foreign_key;unique_index:idx_profit"`
	Quantity float64
	Bought   float64
	Sold     float64
	Fee      float64
	Tax      float64
	Total    float64
	ProfitEntity
}

func NewProfitDAO(ctx *common.Context) ProfitDAO {
	ctx.DB.AutoMigrate(&Profit{})
	return &ProfitDAOImpl{ctx: ctx}
}

func (dao *ProfitDAOImpl) Create(profit ProfitEntity) error {
	return dao.ctx.DB.Create(profit).Error
}

func (dao *ProfitDAOImpl) Save(profit ProfitEntity) error {
	return dao.ctx.DB.Save(profit).Error
}

func (dao *ProfitDAOImpl) Find() ([]Profit, error) {
	var profits []Profit
	daoUser := &User{Id: dao.ctx.User.Id}
	if err := dao.ctx.DB.Model(daoUser).Related(&profits).Error; err != nil {
		return nil, err
	}
	return profits, nil
}

func (dao *ProfitDAOImpl) GetByTrade(trade TradeEntity) (ProfitEntity, error) {
	var profit Profit
	if err := dao.ctx.DB.Model(trade).Related(&profit).Error; err != nil {
		return nil, err
	}
	return &profit, nil
}

func (entity *Profit) GetId() uint {
	return entity.Id
}

func (entity *Profit) GetUserId() uint {
	return entity.UserId
}

func (entity *Profit) GetTradeId() uint {
	return entity.TradeId
}

func (entity *Profit) GetQuantity() float64 {
	return entity.Quantity
}

func (entity *Profit) GetBought() float64 {
	return entity.Bought
}

func (entity *Profit) GetSold() float64 {
	return entity.Sold
}

func (entity *Profit) GetFee() float64 {
	return entity.Fee
}

func (entity *Profit) GetTax() float64 {
	return entity.Tax
}

func (entity *Profit) GetTotal() float64 {
	return entity.Total
}
