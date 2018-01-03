package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/util"
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
	Id       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"type:varchar(100);unique_index"`
	Wallets  []Wallet
}

type Wallet struct {
	UserID   uint
	Currency string `gorm:"primary_key"`
	Address  string `gorm:"unique_index"`
}

func NewUserDAO(ctx *common.Context) *UserDAO {
	ctx.DB.AutoMigrate(&User{})
	ctx.DB.AutoMigrate(&Wallet{})
	return &UserDAO{
		ctx:   ctx,
		Users: make([]User, 0)}
}

func (dao *UserDAO) Get(userId uint) *common.User {
	var user User
	user.Id = userId
	if err := dao.ctx.DB.First(&user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Get] Error: %s", err.Error())
	}
	return &common.User{
		Id:       user.Id,
		Username: user.Username}
}

func (dao *UserDAO) Create(user *User) {
	if err := dao.ctx.DB.Create(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Create] Error:%s", err.Error())
	}
}

func (dao *UserDAO) GetWallets(user *common.User) []common.CryptoWallet {
	var wallets []Wallet
	var walletList []common.CryptoWallet
	if err := dao.ctx.DB.Preload("Users").Find(&wallets).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetWallets] Error: %s", err.Error())
	}
	for _, wallet := range wallets {
		balance := dao.getBalance(wallet.Currency, wallet.Address)
		walletList = append(walletList, common.CryptoWallet{
			Address:  wallet.Address,
			Currency: wallet.Currency,
			Balance:  balance,
			NetWorth: dao.getPrice(wallet.Currency, balance)})
	}
	return walletList
}

func (dao *UserDAO) GetWallet(user *common.User, currency string) *common.CryptoWallet {
	var wallets []Wallet
	if err := dao.ctx.DB.Preload("Users").Find(&wallets).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetWallet] Error: %s", err.Error())
	}
	for _, w := range wallets {
		if w.Currency == currency {
			balance := dao.getBalance(currency, w.Address)
			return &common.CryptoWallet{
				Address:  w.Address,
				Currency: w.Currency,
				Balance:  balance,
				NetWorth: dao.getPrice(currency, balance)}
		}
	}
	return &common.CryptoWallet{}
}

func (dao *UserDAO) getBalance(currency, address string) float64 {
	dao.ctx.Logger.Debugf("[UserDAO.getBalance] currency=%s, address=%s", currency, address)
	if currency == "XRP" {
		return service.NewRipple(dao.ctx).GetBalance(address).Balance
	} else if currency == "BTC" {
		return service.NewBlockchainInfo(dao.ctx).GetBalance(address).Balance
	}
	return 0.0
}

func (dao *UserDAO) getPrice(currency string, amt float64) float64 {
	dao.ctx.Logger.Debugf("[UserDAO.getPrice] currency=%s, amt=%.8f", currency, amt)
	if currency == "BTC" {
		return util.TruncateFloat(service.NewBlockchainInfo(dao.ctx).GetPrice()*amt, 8)
	}
	return 0.0
}
