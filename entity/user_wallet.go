package entity

type UserWallet struct {
	UserId   uint
	Currency string `gorm:"primary_key"`
	Address  string `gorm:"primary_key"`
	Native   uint
	UserWalletEntity
}

func (entity *UserWallet) GetUserId() uint {
	return entity.UserId
}

func (entity *UserWallet) GetCurrency() string {
	return entity.Currency
}

func (entity *UserWallet) GetAddress() string {
	return entity.Address
}

func (entity *UserWallet) IsNative() bool {
	return entity.Native > 0
}
