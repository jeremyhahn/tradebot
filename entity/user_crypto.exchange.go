package entity

type UserCryptoExchange struct {
	UserID  uint   `gorm:"type:integer;primary_key"`
	Name    string `gorm:"type:varchar(200);primary_key"`
	Key     string `gorm:"not null" sql:"type:varchar(255)"`
	Secret  string `gorm:"not null" sql:"type:text"`
	Markets string `sql:"type:text"`
	Extra   string `gorm:"not null" sql:"type:varchar(255)"`
	UserExchangeEntity
}

func (entity *UserCryptoExchange) GetUserID() uint {
	return entity.UserID
}

func (entity *UserCryptoExchange) GetName() string {
	return entity.Name
}

func (entity *UserCryptoExchange) GetKey() string {
	return entity.Key
}

func (entity *UserCryptoExchange) GetSecret() string {
	return entity.Secret
}

func (entity *UserCryptoExchange) GetMarkets() string {
	return entity.Markets
}

func (entity *UserCryptoExchange) GetExtra() string {
	return entity.Extra
}
