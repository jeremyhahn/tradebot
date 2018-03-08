package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type UserDAO interface {
	GetUsername() string
	GetById(userId uint) (entity.UserEntity, error)
	GetByName(username string) (entity.UserEntity, error)
	Create(user entity.UserEntity) error
	Save(user entity.UserEntity) error
	Update(user entity.UserEntity) error
	Find() ([]entity.User, error)
	GetWallets(user entity.UserEntity) []entity.UserWallet
	GetWallet(user entity.UserEntity, currency string) entity.UserWalletEntity
	GetTokens(user entity.UserEntity) []entity.UserToken
	GetToken(user entity.UserEntity, symbol string) entity.UserTokenEntity
	GetExchanges(user entity.UserEntity) []entity.UserCryptoExchange
	GetExchange(user entity.UserEntity, exchangeName string) (entity.UserExchangeEntity, error)
	CreateExchange(userExchange entity.UserExchangeEntity) error
}

type UserDAOImpl struct {
	ctx common.Context
	UserDAO
}

func NewUserDAO(ctx common.Context) UserDAO {
	return &UserDAOImpl{ctx: ctx}
}

func CreateUserDAO(ctx common.Context, user common.UserContext) UserDAO {
	//ctx.SetUser(user)
	return &UserDAOImpl{ctx: ctx}
}

func (dao *UserDAOImpl) GetById(userId uint) (entity.UserEntity, error) {
	var user entity.User
	user.Id = userId
	if err := dao.ctx.GetCoreDB().First(&user, userId).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.GetById] Error: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) GetByName(username string) (entity.UserEntity, error) {
	var user entity.User
	if err := dao.ctx.GetCoreDB().First(&user, "username = ?", username).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.GetByName] %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) Create(user entity.UserEntity) error {
	if err := dao.ctx.GetCoreDB().Create(user).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.Create] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Save(user entity.UserEntity) error {
	if err := dao.ctx.GetCoreDB().Save(user).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.Save] Error:%s", err.Error())
		return err
	}
	return nil
}

func (dao *UserDAOImpl) Find() ([]entity.User, error) {
	var users []entity.User
	if err := dao.ctx.GetCoreDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAOImpl) GetWallets(user entity.UserEntity) []entity.UserWallet {
	var wallets []entity.UserWallet
	if err := dao.ctx.GetCoreDB().Model(user).Related(&wallets).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.GetWallets] Error: %s", err.Error())
	}
	return wallets
}

func (dao *UserDAOImpl) GetWallet(user entity.UserEntity, currency string) entity.UserWalletEntity {
	wallets := dao.GetWallets(user)
	for _, w := range wallets {
		if w.Currency == currency {
			return &w
		}
	}
	return &entity.UserWallet{}
}

func (dao *UserDAOImpl) GetTokens(user entity.UserEntity) []entity.UserToken {
	var tokens []entity.UserToken
	if err := dao.ctx.GetCoreDB().Model(user).Related(&tokens).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.GetTokens] Error: %s", err.Error())
	}
	return tokens
}

func (dao *UserDAOImpl) GetToken(user entity.UserEntity, symbol string) entity.UserTokenEntity {
	tokens := dao.GetTokens(user)
	for _, t := range tokens {
		if t.Symbol == symbol {
			return &t
		}
	}
	return &entity.UserToken{}
}

func (dao *UserDAOImpl) GetExchanges(user entity.UserEntity) []entity.UserCryptoExchange {
	var exchanges []entity.UserCryptoExchange
	if err := dao.ctx.GetCoreDB().Model(user).Related(&exchanges).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.GetExchanges] Error: %s", err.Error())
	}
	return exchanges
}

func (dao *UserDAOImpl) GetExchange(user entity.UserEntity, exchangeName string) (entity.UserExchangeEntity, error) {
	exchanges := dao.GetExchanges(user)
	for _, ex := range exchanges {
		if ex.Name == exchangeName {
			return &ex, nil
		}
	}
	errmsg := fmt.Sprintf("User exchange not found: %s", exchangeName)
	return nil, errors.New(errmsg)
}

func (dao *UserDAOImpl) CreateExchange(userExchange entity.UserExchangeEntity) error {
	if err := dao.ctx.GetCoreDB().Create(userExchange).Error; err != nil {
		dao.ctx.GetLogger().Errorf("[UserDAO.CreateExchange] Error:%s", err.Error())
		return err
	}
	return nil
}
