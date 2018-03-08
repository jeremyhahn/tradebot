package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type UserMapper interface {
	MapUserEntityToDto(entity entity.UserEntity) common.UserContext
	MapUserDtoToEntity(dto common.UserContext) entity.UserEntity
}

type DefaultUserMapper struct {
	ctx common.Context
}

func NewUserMapper() UserMapper {
	return &DefaultUserMapper{}
}

func (mapper *DefaultUserMapper) MapUserEntityToDto(entity entity.UserEntity) common.UserContext {
	return &dto.UserContextDTO{
		Id:            entity.GetId(),
		Username:      entity.GetUsername(),
		LocalCurrency: entity.GetLocalCurrency(),
		Etherbase:     entity.GetEtherbase(),
		Keystore:      entity.GetKeystore()}
}

func (mapper *DefaultUserMapper) MapUserDtoToEntity(dto common.UserContext) entity.UserEntity {
	return &entity.User{
		Id:            dto.GetId(),
		Username:      dto.GetUsername(),
		LocalCurrency: dto.GetLocalCurrency(),
		Etherbase:     dto.GetEtherbase(),
		Keystore:      dto.GetKeystore()}
}
