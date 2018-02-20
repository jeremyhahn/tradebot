package entity

type CryptoExchange struct {
	UserID uint
	Name   string `gorm:"primary_key" sql:"type:varchar(255)"`
	URL    string `gorm:"not null" sql:"type:varchar(255)"`
	Key    string `gorm:"not null" sql:"type:varchar(255)"`
	Secret string `gorm:"not null" sql:"type:text"`
	Extra  string `gorm:"not null" sql:"type:varchar(255)"`
}
