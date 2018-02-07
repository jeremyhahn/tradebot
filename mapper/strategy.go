package mapper

import (
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type StrategyMapper interface {
	MapStrategyEntityToDto(entity dao.StrategyEntity) dto.PlatformStrategy
	MapStrategyDtoToEntity(dto dto.PlatformStrategy) dao.StrategyEntity
}

type DefaultStrategyMapper struct {
	StrategyMapper
}

func NewStrategyMapper() StrategyMapper {
	return &DefaultStrategyMapper{}
}

func (mapper *DefaultStrategyMapper) MapStrategyEntityToDto(entity dao.StrategyEntity) dto.PlatformStrategy {
	return &dto.PlatformStrategyDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultStrategyMapper) MapStrategyDtoToEntity(platformStrategy dto.PlatformStrategy) dao.StrategyEntity {
	return &dao.Strategy{
		Name:     platformStrategy.GetName(),
		Filename: platformStrategy.GetFilename(),
		Version:  platformStrategy.GetVersion()}
}
