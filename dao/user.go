package dao

import (
	"errors"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
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
	ctx.DB.AutoMigrate(&entity.User{})
	ctx.DB.AutoMigrate(&entity.UserWallet{})
	ctx.DB.AutoMigrate(&entity.UserCryptoExchange{})
	return &UserDAOImpl{ctx: ctx}
}

func CreateUserDAO(ctx *common.Context, user common.User) UserDAO {
	ctx.DB.AutoMigrate(&entity.User{})
	ctx.DB.AutoMigrate(&entity.UserWallet{})
	ctx.DB.AutoMigrate(&entity.UserCryptoExchange{})
	ctx.SetUser(user)
	return &UserDAOImpl{ctx: ctx}
}

func (dao *UserDAOImpl) GetById(userId uint) (entity.UserEntity, error) {
	var user entity.User
	user.Id = userId
	if err := dao.ctx.DB.First(&user, userId).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetById] Error: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) GetByName(username string) (entity.UserEntity, error) {
	var user entity.User
	if err := dao.ctx.DB.First(&user, "username = ?", username).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.GetByName] Error: %s", err.Error())
		return nil, err
	}
	if user.GetId() == 0 {
		util.DUMP(user)
		dao.ctx.Logger.Warningf("[UserDAO.GetByName] Unable to locate user: %s", username)
		return nil, errors.New("User not found")
	}
	return &user, nil
}

func (dao *UserDAOImpl) Create(user *entity.User) error {
	if err := dao.ctx.DB.Create(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Create] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Save(user *entity.User) error {
	if err := dao.ctx.DB.Save(user).Error; err != nil {
		dao.ctx.Logger.Errorf("[UserDAO.Save] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Find() ([]entity.User, error) {
	var users []entity.User
	if err := dao.ctx.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAOImpl) GetWallets(user *entity.User) []entity.UserWallet {
	var wallets []entity.UserWallet
	if err := dao.ctx.DB.Model(user).Related(&wallets).Error; err != nil {
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
	if err := dao.ctx.DB.Model(user).Related(&exchanges).Error; err != nil {
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
