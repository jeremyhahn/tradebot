package mapper

import (
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type UserExchangeMapper interface {
	MapEntityToViewModel(entity entity.UserExchangeEntity) *viewmodel.UserCryptoExchange
	MapViewModelToEntity(vm viewmodel.UserCryptoExchange) entity.UserExchangeEntity
}

type UserExchangeMapperImpl struct {
	UserExchangeMapper
}

func NewExchangeMapper() UserExchangeMapper {
	return &UserExchangeMapperImpl{}
}

func (em *UserExchangeMapperImpl) MapEntityToViewModel(entity entity.UserExchangeEntity) *viewmodel.UserCryptoExchange {
	return &viewmodel.UserCryptoExchange{
		Id:    entity.GetName(),
		Name:  entity.GetName(),
		Key:   entity.GetKey(),
		URL:   entity.GetURL(),
		Extra: entity.GetExtra()}
}

func (em *UserExchangeMapperImpl) MapViewModelToEntity(vm viewmodel.UserCryptoExchange) entity.UserExchangeEntity {
	return &entity.UserCryptoExchange{
		Name:  vm.Name,
		Key:   vm.Key,
		URL:   vm.URL,
		Extra: vm.Extra}
}
