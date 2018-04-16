package entity

import "time"

type Transaction struct {
	Id                     string `gorm:"type:varchar(200);primary_key"`
	UserId                 uint
	Date                   time.Time
	MarketPair             string `gorm:"type:varchar(10)"`
	CurrencyPair           string `gorm:"type:varchar(10)"`
	Type                   string `gorm:"type:varchar(64)"`
	Category               string `gorm:"type:varchar(200)"`
	Network                string `gorm:"type:varchar(200)"`
	NetworkDisplayName     string `gorm:"type:varchar(200)"`
	Quantity               string `gorm:"type:varchar(64)"`
	QuantityCurrency       string `gorm:"type:varchar(6)"`
	FiatQuantity           string `gorm:"type:varchar(64)"`
	FiatQuantityCurrency   string `gorm:"type:varchar(6)"`
	Price                  string `gorm:"type:varchar(64)"`
	PriceCurrency          string `gorm:"type:varchar(6)"`
	FiatPrice              string `gorm:"type:varchar(64)"`
	FiatPriceCurrency      string `gorm:"type:varchar(6)"`
	QuoteFiatPrice         string `gorm:"type:varchar(64)"`
	QuoteFiatPriceCurrency string `gorm:"type:varchar(6)"`
	Fee                    string `gorm:"type:varchar(64)"`
	FeeCurrency            string `gorm:"type:varchar(6)"`
	FiatFee                string `gorm:"type:varchar(64)"`
	FiatFeeCurrency        string `gorm:"type:varchar(6)"`
	Total                  string `gorm:"type:varchar(64)"`
	TotalCurrency          string `gorm:"type:varchar(6)"`
	FiatTotal              string `gorm:"type:varchar(64)"`
	FiatTotalCurrency      string `gorm:"type:varchar(6)"`
	Deleted                int    `gorm:"type:integer;default:0"`
	TransactionEntity
}

func (tx *Transaction) GetId() string {
	return tx.Id
}

func (tx *Transaction) GetUserId() uint {
	return tx.UserId
}

func (tx *Transaction) GetDate() time.Time {
	return tx.Date
}

func (tx *Transaction) GetMarketPair() string {
	return tx.MarketPair
}

func (tx *Transaction) GetCurrencyPair() string {
	return tx.CurrencyPair
}

func (tx *Transaction) GetType() string {
	return tx.Type
}

func (tx *Transaction) GetCategory() string {
	return tx.Category
}

func (tx *Transaction) SetCategory(category string) {
	tx.Category = category
}

func (tx *Transaction) GetNetwork() string {
	return tx.Network
}

func (tx *Transaction) GetNetworkDisplayName() string {
	return tx.NetworkDisplayName
}

func (tx *Transaction) GetQuantity() string {
	return tx.Quantity
}

func (tx *Transaction) GetQuantityCurrency() string {
	return tx.QuantityCurrency
}

func (tx *Transaction) GetFiatQuantity() string {
	return tx.FiatQuantity
}

func (tx *Transaction) GetFiatQuantityCurrency() string {
	return tx.FiatQuantityCurrency
}

func (tx *Transaction) GetPrice() string {
	return tx.Price
}

func (tx *Transaction) GetPriceCurrency() string {
	return tx.PriceCurrency
}

func (tx *Transaction) GetFiatPrice() string {
	return tx.FiatPrice
}

func (tx *Transaction) GetFiatPriceCurrency() string {
	return tx.FiatPriceCurrency
}

func (tx *Transaction) GetQuoteFiatPrice() string {
	return tx.QuoteFiatPrice
}

func (tx *Transaction) GetQuoteFiatPriceCurrency() string {
	return tx.QuoteFiatPriceCurrency
}

func (tx *Transaction) GetFee() string {
	return tx.Fee
}

func (tx *Transaction) GetFeeCurrency() string {
	return tx.FeeCurrency
}

func (tx *Transaction) GetFiatFee() string {
	return tx.FiatFee
}

func (tx *Transaction) GetFiatFeeCurrency() string {
	return tx.FiatFeeCurrency
}

func (tx *Transaction) GetTotal() string {
	return tx.Total
}

func (tx *Transaction) GetTotalCurrency() string {
	return tx.TotalCurrency
}

func (tx *Transaction) GetFiatTotal() string {
	return tx.FiatTotal
}

func (tx *Transaction) GetFiatTotalCurrency() string {
	return tx.FiatTotalCurrency
}

func (tx *Transaction) IsDeleted() bool {
	return tx.Deleted > 0
}

func (tx *Transaction) SetDeleted(value int) {
	tx.Deleted = value
}
