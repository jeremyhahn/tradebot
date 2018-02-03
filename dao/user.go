package dao

import (
	"github.com/jeremyhahn/tradebot/common"
)

type IUser interface {
	GetUsername() string
}

type UserDAO struct {
	ctx   *common.Context
	Users []User
	IUser
}

type User struct {
	Id            uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Username      string `gorm:"type:varchar(100);unique_index"`
	LocalCurrency string `gorm:"type:varchar(5)"`
	Charts        []Chart
	Wallets       []UserWallet
	Exchanges     []UserCryptoExchange
}

type UserWallet struct {
	UserID   uint
	Currency string `gorm:"primary_key"`
	Address  string `gorm:"unique_index"`
}

type UserCryptoExchange struct {
	UserID uint
	Name   string `gorm:"primary_key"`
	URL    string `gorm:"not null" sql:"type:varchar(255)"`
	Key    string `gorm:"not null" sql:"type:varchar(255)"`
	Secret string `gorm:"not null" sql:"type:text"`
	Extra  string `gorm:"not null" sql:"type:varchar(255)"`
}

func NewUserDAO(ctx *common.Context) *UserDAO {
	ctx.DB.AutoMigrate(&User{})
	ctx.DB.AutoMigrate(&UserWallet{})
	ctx.DB.AutoMigrate(&UserCryptoExchange{})
	return &UserDAO{
		ctx:   ctx,
		Users: make([]User, 0)}
}

func CreateUserDAO(ctx *common.Context, user *common.User) *UserDAO {
	ctx.DB.AutoMigrate(&User{})
	ctx.DB.AutoMigrate(&UserWallet{})
	ctx.DB.AutoMigrate(&UserCryptoExchange{})
	ctx.User = user
	return &UserDAO{
		ctx:   ctx,
		Users: make([]User, 0)}
}

func (dao *UserDAO) GetById(userId uint) *common.User {
	var user User
	user.Id = userId
	if err := dao.ctx.DB.First(&user, userId).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetById] Error: %s", err.Error())
	}
	return &common.User{
		Id:            user.Id,
		Username:      user.Username,
		LocalCurrency: user.LocalCurrency}
}

func (dao *UserDAO) GetByName(username string) *common.User {
	var user User
	user.Username = username
	if err := dao.ctx.DB.Where("username = ?", username).First(&user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetByName] Error: %s", err.Error())
	}
	return &common.User{
		Id:            user.Id,
		Username:      user.Username,
		LocalCurrency: user.LocalCurrency}
}

func (dao *UserDAO) Create(user *User) bool {
	return dao.ctx.DB.NewRecord(user)
}

func (dao *UserDAO) Save(user *User) {
	if err := dao.ctx.DB.Save(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Save] Error:%s", err.Error())
	}
}

func (dao *UserDAO) GetWallets(user *common.User) []UserWallet {
	var wallets []UserWallet
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&wallets).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetWallets] Error: %s", err.Error())
	}
	return wallets
}

func (dao *UserDAO) GetWallet(user *common.User, currency string) *UserWallet {
	wallets := dao.GetWallets(user)
	for _, w := range wallets {
		if w.Currency == currency {
			return &w
		}
	}
	return &UserWallet{}
}

func (dao *UserDAO) GetExchanges(user *common.User) []UserCryptoExchange {
	var exchanges []UserCryptoExchange
	daoUser := &User{Id: user.Id}
	if err := dao.ctx.DB.Model(daoUser).Related(&exchanges).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetExchanges] Error: %s", err.Error())
	}
	return exchanges
}

func (dao *UserDAO) GetExchange(user *common.User, name string) *UserCryptoExchange {
	var exchange UserCryptoExchange
	exchanges := dao.GetExchanges(user)
	for _, ex := range exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
