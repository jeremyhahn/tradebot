package entity

type Profit struct {
	Id       uint `gorm:"primary_key"`
	UserId   uint `gorm:"unique_index:idx_profit"`
	TradeId  uint `gorm:"foreign_key;unique_index:idx_profit"`
	Quantity string
	Bought   string
	Sold     string
	Fee      string
	Tax      string
	Total    string
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

func (entity *Profit) GetQuantity() string {
	return entity.Quantity
}

func (entity *Profit) GetBought() string {
	return entity.Bought
}

func (entity *Profit) GetSold() string {
	return entity.Sold
}

func (entity *Profit) GetFee() string {
	return entity.Fee
}

func (entity *Profit) GetTax() string {
	return entity.Tax
}

func (entity *Profit) GetTotal() string {
	return entity.Total
}
