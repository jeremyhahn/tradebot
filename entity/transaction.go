package entity

import "time"

type Transaction struct {
	Id               uint      `gorm:"primary_key"`
	UserId           uint      `gorm:"unique_index:idx_txhistory"`
	Date             time.Time `gorm:"unique_index:idx_txhistory"`
	Network          string    `gorm:"unique_index:idx_txhistory"`
	Type             string
	Currency         string `gorm:"unique_index:idx_txhistory"`
	Quantity         string `gorm:"type:varchar(64)"`
	QuantityCurrency string `gorm:"type:varchar(6)"`
	Price            string `gorm:"type:varchar(64)"`
	PriceCurrency    string `gorm:"type:varchar(6)"`
	Fee              string `gorm:"type:varchar(64)"`
	FeeCurrency      string `gorm:"type:varchar(6)"`
	Total            string `gorm:"type:varchar(64);unique_index:idx_txhistory"`
	TotalCurrency    string `gorm:"type:varchar(6)"`
	TransactionEntity
}

func (tx *Transaction) GetId() uint {
	return tx.Id
}

func (tx *Transaction) GetUserId() uint {
	return tx.UserId
}

func (tx *Transaction) GetDate() time.Time {
	return tx.Date
}

func (tx *Transaction) GetNetwork() string {
	return tx.Network
}

func (tx *Transaction) GetType() string {
	return tx.Type
}

func (tx *Transaction) GetCurrency() string {
	return tx.Currency
}

func (tx *Transaction) GetQuantity() string {
	return tx.Quantity
}

func (tx *Transaction) GetPrice() string {
	return tx.Price
}

func (tx *Transaction) GetPriceCurrency() string {
	return tx.PriceCurrency
}

func (tx *Transaction) GetFee() string {
	return tx.Fee
}

func (tx *Transaction) GetFeeCurrency() string {
	return tx.FeeCurrency
}

func (tx *Transaction) GetTotal() string {
	return tx.Total
}

func (tx *Transaction) GetTotalCurrency() string {
	return tx.TotalCurrency
}

func (tx *Transaction) GetQuantityCurrency() string {
	return tx.QuantityCurrency
}
