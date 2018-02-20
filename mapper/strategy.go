package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type StrategyMapper interface {
	MapStrategyEntityToDto(entity entity.StrategyEntity) common.Strategy
	MapStrategyDtoToEntity(dto common.Strategy) entity.StrategyEntity
}

type DefaultStrategyMapper struct {
	StrategyMapper
}

func NewStrategyMapper() StrategyMapper {
	return &DefaultStrategyMapper{}
}

func (mapper *DefaultStrategyMapper) MapStrategyEntityToDto(entity entity.StrategyEntity) common.Strategy {
	return &dto.StrategyDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultStrategyMapper) MapStrategyDtoToEntity(platformStrategy common.Strategy) entity.StrategyEntity {
	return &entity.Strategy{
		Name:     platformStrategy.GetName(),
		Filename: platformStrategy.GetFilename(),
		Version:  platformStrategy.GetVersion()}
}
