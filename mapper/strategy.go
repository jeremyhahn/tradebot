package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type StrategyMapper interface {
	MapStrategyEntityToDto(entity dao.StrategyEntity) common.Strategy
	MapStrategyDtoToEntity(dto common.Strategy) dao.StrategyEntity
}

type DefaultStrategyMapper struct {
	StrategyMapper
}

func NewStrategyMapper() StrategyMapper {
	return &DefaultStrategyMapper{}
}

func (mapper *DefaultStrategyMapper) MapStrategyEntityToDto(entity dao.StrategyEntity) common.Strategy {
	return &dto.StrategyDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultStrategyMapper) MapStrategyDtoToEntity(platformStrategy common.Strategy) dao.StrategyEntity {
	return &dao.Strategy{
		Name:     platformStrategy.GetName(),
		Filename: platformStrategy.GetFilename(),
		Version:  platformStrategy.GetVersion()}
}
