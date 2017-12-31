package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type IUser interface {
	GetUsername() string
}

type UserDAO struct {
	db     *gorm.DB
	logger *logging.Logger
	Users  []User
	IUser
}

type User struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"type:varchar(100);unique_index"`
}

func NewUserDAO(db *gorm.DB, logger *logging.Logger) *UserDAO {
	db.AutoMigrate(&User{})
	return &UserDAO{
		db:     db,
		logger: logger,
		Users:  make([]User, 0)}
}

func (dao *UserDAO) Get(userId int64) *common.User {
	var user User
	if err := dao.db.First(&user).Error; err != nil {
		dao.logger.Errorf("[UserDAO.Get] Error: %s", err.Error())
	}
	return &common.User{
		Id:       user.Id,
		Username: user.Username}
}

func (dao *UserDAO) Create(user *User) {
	if err := dao.db.Create(user).Error; err != nil {
		dao.logger.Errorf("[UserDAO.Create] Error:%s", err.Error())
	}
}
