package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type IndicatorMapper interface {
	MapIndicatorEntityToDto(entity entity.IndicatorEntity) common.Indicator
	MapIndicatorDtoToEntity(dto common.Indicator) entity.IndicatorEntity
}

type DefaultIndicatorMapper struct {
	IndicatorMapper
}

func NewIndicatorMapper() IndicatorMapper {
	return &DefaultIndicatorMapper{}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorEntityToDto(entity entity.IndicatorEntity) common.Indicator {
	return &dto.IndicatorDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorDtoToEntity(platformIndicator common.Indicator) entity.IndicatorEntity {
	return &entity.Indicator{
		Name:     platformIndicator.GetName(),
		Filename: platformIndicator.GetFilename(),
		Version:  platformIndicator.GetVersion()}
}
