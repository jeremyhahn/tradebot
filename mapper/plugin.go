package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type PluginMapper interface {
	MapPluginEntityToDto(entity entity.PluginEntity) common.Plugin
	MapPluginDtoToEntity(dto common.Plugin) entity.PluginEntity
}

type DefaultPluginMapper struct {
	PluginMapper
}

func NewPluginMapper() PluginMapper {
	return &DefaultPluginMapper{}
}

func (mapper *DefaultPluginMapper) MapPluginEntityToDto(entity entity.PluginEntity) common.Plugin {
	return &dto.PluginDTO{
		Name:     entity.GetName(),
		Filename: entity.GetFilename(),
		Version:  entity.GetVersion(),
		Type:     entity.GetType()}
}

func (mapper *DefaultPluginMapper) MapPluginDtoToEntity(platformPlugin common.Plugin) entity.PluginEntity {
	return &entity.Plugin{
		Name:     platformPlugin.GetName(),
		Filename: platformPlugin.GetFilename(),
		Version:  platformPlugin.GetVersion(),
		Type:     platformPlugin.GetType()}
}
