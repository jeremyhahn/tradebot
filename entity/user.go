package entity

type User struct {
	Id            uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Username      string `gorm:"type:varchar(100);unique_index"`
	LocalCurrency string `gorm:"type:varchar(5)"`
	FiatExchange  string
	Etherbase     string `gorm:"type:varchar(160)"`
	Keystore      string
	Charts        []Chart
	Wallets       []UserWallet
	Tokens        []UserToken
	Exchanges     []UserCryptoExchange
	UserEntity
}

func (entity *User) GetId() uint {
	return entity.Id
}

func (entity *User) GetUsername() string {
	return entity.Username
}

func (entity *User) GetLocalCurrency() string {
	return entity.LocalCurrency
}

func (entity *User) GetFiatExchange() string {
	return entity.FiatExchange
}

func (entity *User) GetEtherbase() string {
	return entity.Etherbase
}

func (entity *User) GetKeystore() string {
	return entity.Keystore
}

func (entity *User) GetCharts() []Chart {
	return entity.Charts
}

func (entity *User) GetWallets() []UserWallet {
	return entity.Wallets
}

func (entity *User) GetTokens() []UserToken {
	return entity.Tokens
}

func (entity *User) GetExchanges() []UserCryptoExchange {
	return entity.Exchanges
}
