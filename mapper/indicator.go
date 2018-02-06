package mapper

import (
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type IndicatorMapper interface {
	MapIndicatorEntityToDto(entity dao.IndicatorEntity) dto.PlatformIndicator
	MapIndicatorDtoToEntity(dto dto.PlatformIndicator) dao.IndicatorEntity
}

type DefaultIndicatorMapper struct {
	IndicatorMapper
}

func NewIndicatorMapper() IndicatorMapper {
	return &DefaultIndicatorMapper{}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorEntityToDto(entity dao.IndicatorEntity) dto.PlatformIndicator {
	return &dto.PlatformIndicatorDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorDtoToEntity(platformIndicator dto.PlatformIndicator) dao.IndicatorEntity {
	return &dao.Indicator{
		Name:     platformIndicator.GetName(),
		Filename: platformIndicator.GetFilename(),
		Version:  platformIndicator.GetVersion()}
}
