package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type UserExchangeMapper interface {
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

func (em *UserExchangeMapperImpl) MapEntityToDto(entity entity.UserExchangeEntity) common.UserCryptoExchange {
	return &dto.UserCryptoExchangeDTO{
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

func (em *UserExchangeMapperImpl) MapViewModelToEntity(vm viewmodel.UserCryptoExchange) entity.UserExchangeEntity {
	return &entity.UserCryptoExchange{
		Name:  vm.Name,
		Key:   vm.Key,
		Extra: vm.Extra}
}
