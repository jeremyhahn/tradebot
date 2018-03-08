package entity

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
