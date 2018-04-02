package service

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"golang.org/x/crypto/bcrypt"
)

type LocalAuthService struct {
	ctx    common.Context
	dao    dao.UserDAO
	mapper mapper.UserMapper
}

func NewLocalAuthService(ctx common.Context, userDAO dao.UserDAO, userMapper mapper.UserMapper) AuthService {
	return &LocalAuthService{
		ctx:    ctx,
		dao:    userDAO,
		mapper: userMapper}
}

func (service *LocalAuthService) Login(username, password string) (common.UserContext, error) {
	userEntity, err := service.dao.GetByName(username)
	if err != nil && err.Error() != "record not found" {
		return nil, errors.New("Invalid username/password")
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.GetKeystore()), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid username/password")
	}
	return service.mapper.MapUserEntityToDto(userEntity), nil
}

func (service *LocalAuthService) Register(username, password string) error {
	_, err := service.dao.GetByName(username)
	if err != nil && err.Error() != "record not found" {
		service.ctx.GetLogger().Errorf("[LocalAuthService.Register] %s", err.Error())
		return errors.New(fmt.Sprintf("Unexpected error: %s", err.Error()))
	}
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return service.dao.Save(&entity.User{
		Username:      username,
		LocalCurrency: "USD",
		Etherbase:     "etherscan",
		Keystore:      string(encrypted)})
}
