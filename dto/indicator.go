package dto

import "github.com/jeremyhahn/tradebot/common"

type IndicatorDTO struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	common.Indicator
}

func NewIndicatorDTO() common.Indicator {
	return &IndicatorDTO{}
}

func CreateIndicatorDTO(name, filename, version string) common.Indicator {
	return &IndicatorDTO{
		Name:     name,
		Filename: filename,
		Version:  version}
}

func (dto *IndicatorDTO) GetName() string {
	return dto.Name
}

func (dto *IndicatorDTO) GetFilename() string {
	return dto.Filename
}

func (dto *IndicatorDTO) GetVersion() string {
	return dto.Version
}
