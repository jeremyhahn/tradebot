package entity

type UserToken struct {
	UserId          uint
	Symbol          string `gorm:"primary_key"`
	ContractAddress string `gorm:"unique_index"`
	WalletAddress   string
	UserTokenEntity
}

func (entity *UserToken) GetUserId() uint {
	return entity.UserId
}

func (entity *UserToken) GetSymbol() string {
	return entity.Symbol
}

func (entity *UserToken) GetContractAddress() string {
	return entity.ContractAddress
}

func (entity *UserToken) GetWalletAddress() string {
	return entity.WalletAddress
}
