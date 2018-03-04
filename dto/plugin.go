package dto

import "github.com/jeremyhahn/tradebot/common"

type PluginDTO struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	Type     string `json:"type"`
	common.Plugin
}

func NewPluginDTO() common.Plugin {
	return &PluginDTO{}
}

func CreatePluginDTO(name, filename, version, pluginType string) common.Plugin {
	return &PluginDTO{
		Name:     name,
		Filename: filename,
		Version:  version,
		Type:     pluginType}
}

func (dto *PluginDTO) GetName() string {
	return dto.Name
}

func (dto *PluginDTO) GetFilename() string {
	return dto.Filename
}

func (dto *PluginDTO) GetVersion() string {
	return dto.Version
}

func (dto *PluginDTO) GetType() string {
	return dto.Type
}
