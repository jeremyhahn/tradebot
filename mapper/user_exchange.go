package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type UserExchangeMapper interface {
	MapDtoToViewModel(dto common.UserCryptoExchange) *viewmodel.UserCryptoExchange
	MapDtoToEntity(userCryptoExchange common.UserCryptoExchange) entity.UserExchangeEntity
	MapEntityToDto(entity entity.UserExchangeEntity) common.UserCryptoExchange
	MapEntityToViewModel(entity entity.UserExchangeEntity) *viewmodel.UserCryptoExchange
	MapViewModelToEntity(vm viewmodel.UserCryptoExchange) entity.UserExchangeEntity
}

type UserExchangeMapperImpl struct {
	UserExchangeMapper
}

func NewUserExchangeMapper() UserExchangeMapper {
	return &UserExchangeMapperImpl{}
}

func (em *UserExchangeMapperImpl) MapDtoToEntity(userCryptoExchange common.UserCryptoExchange) entity.UserExchangeEntity {
	return &entity.UserCryptoExchange{
		UserID: userCryptoExchange.GetUserID(),
		Name:   userCryptoExchange.GetName(),
		Key:    userCryptoExchange.GetKey(),
		Secret: userCryptoExchange.GetSecret(),
		Extra:  userCryptoExchange.GetExtra()}
}

func (em *UserExchangeMapperImpl) MapEntityToDto(entity entity.UserExchangeEntity) common.UserCryptoExchange {
	return &dto.UserCryptoExchangeDTO{
		UserID: entity.GetUserID(),
		Name:   entity.GetName(),
		Key:    entity.GetKey(),
		Secret: entity.GetSecret(),
		Extra:  entity.GetExtra()}
}

func (em *UserExchangeMapperImpl) MapEntityToViewModel(entity entity.UserExchangeEntity) *viewmodel.UserCryptoExchange {
	return &viewmodel.UserCryptoExchange{
		Id:    entity.GetName(),
		Name:  entity.GetName(),
		Key:   entity.GetKey(),
		Extra: entity.GetExtra()}
}

func (em *UserExchangeMapperImpl) MapDtoToViewModel(dto common.UserCryptoExchange) *viewmodel.UserCryptoExchange {
	return &viewmodel.UserCryptoExchange{
		Id:    dto.GetName(),
		Name:  dto.GetName(),
		Key:   dto.GetKey(),
		Extra: dto.GetExtra()}
}

func (em *UserExchangeMapperImpl) MapViewModelToEntity(vm viewmodel.UserCryptoExchange) entity.UserExchangeEntity {
	return &entity.UserCryptoExchange{
		Name:  vm.Name,
		Key:   vm.Key,
		Extra: vm.Extra}
}
