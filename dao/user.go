package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type UserDAO interface {
	GetUsername() string
	GetById(userId uint) (entity.UserEntity, error)
	GetByName(username string) (entity.UserEntity, error)
	Create(user *entity.User) error
	Save(user *entity.User) error
	Update(user *entity.User) error
	Find() ([]entity.User, error)
	GetWallets(user *entity.User) []entity.UserWallet
	GetWallet(user *entity.User, currency string) entity.UserWalletEntity
	GetExchanges(user *entity.User) []entity.UserCryptoExchange
	GetExchange(user *entity.User, name string) *entity.UserCryptoExchange
}

type UserDAOImpl struct {
	ctx *common.Context
	UserDAO
}

func NewUserDAO(ctx *common.Context) UserDAO {
	ctx.CoreDB.AutoMigrate(&entity.User{})
	ctx.CoreDB.AutoMigrate(&entity.UserWallet{})
	ctx.CoreDB.AutoMigrate(&entity.UserCryptoExchange{})
	return &UserDAOImpl{ctx: ctx}
}

func CreateUserDAO(ctx *common.Context, user common.User) UserDAO {
	ctx.CoreDB.AutoMigrate(&entity.User{})
	ctx.CoreDB.AutoMigrate(&entity.UserWallet{})
	ctx.CoreDB.AutoMigrate(&entity.UserCryptoExchange{})
	ctx.SetUser(user)
	return &UserDAOImpl{ctx: ctx}
}

func (dao *UserDAOImpl) GetById(userId uint) (entity.UserEntity, error) {
	var user entity.User
	user.Id = userId
	if err := dao.ctx.CoreDB.First(&user, userId).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetById] Error: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) GetByName(username string) (entity.UserEntity, error) {
	var user entity.User
	if err := dao.ctx.CoreDB.First(&user, "username = ?", username).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetByName] %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) Create(user *entity.User) error {
	if err := dao.ctx.CoreDB.Create(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Create] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Save(user *entity.User) error {
	if err := dao.ctx.CoreDB.Save(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Save] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Find() ([]entity.User, error) {
	var users []entity.User
	if err := dao.ctx.CoreDB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAOImpl) GetWallets(user *entity.User) []entity.UserWallet {
	var wallets []entity.UserWallet
	if err := dao.ctx.CoreDB.Model(user).Related(&wallets).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetWallets] Error: %s", err.Error())
	}
	return wallets
}

func (dao *UserDAOImpl) GetWallet(user *entity.User, currency string) entity.UserWalletEntity {
	wallets := dao.GetWallets(user)
	for _, w := range wallets {
		if w.Currency == currency {
			return &w
		}
	}
	return &entity.UserWallet{}
}

func (dao *UserDAOImpl) GetExchanges(user *entity.User) []entity.UserCryptoExchange {
	var exchanges []entity.UserCryptoExchange
	if err := dao.ctx.CoreDB.Model(user).Related(&exchanges).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetExchanges] Error: %s", err.Error())
	}
	return exchanges
}

func (dao *UserDAOImpl) GetExchange(user *entity.User, name string) *entity.UserCryptoExchange {
	var exchange entity.UserCryptoExchange
	exchanges := dao.GetExchanges(user)
	for _, ex := range exchanges {
		if ex.Name == name {
			return &ex
		}
	}
	return &exchange
}
