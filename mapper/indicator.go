package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type IndicatorMapper interface {
	MapIndicatorEntityToDto(entity dao.IndicatorEntity) common.Indicator
	MapIndicatorDtoToEntity(dto common.Indicator) dao.IndicatorEntity
}

type DefaultIndicatorMapper struct {
	IndicatorMapper
}

func NewIndicatorMapper() IndicatorMapper {
	return &DefaultIndicatorMapper{}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorEntityToDto(entity dao.IndicatorEntity) common.Indicator {
	return &dto.IndicatorDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion()}
}

func (mapper *DefaultIndicatorMapper) MapIndicatorDtoToEntity(platformIndicator common.Indicator) dao.IndicatorEntity {
	return &dao.Indicator{
		Name:     platformIndicator.GetName(),
		Filename: platformIndicator.GetFilename(),
		Version:  platformIndicator.GetVersion()}
}
