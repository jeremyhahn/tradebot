package dto

import "github.com/jeremyhahn/tradebot/common"

type StrategyDTO struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	common.Strategy
}

func NewStrategyDTO() common.Strategy {
	return &StrategyDTO{}
}

func CreateStrategyDTO(name, filename, version string) common.Strategy {
	return &StrategyDTO{
		Name:     name,
		Filename: filename,
		Version:  version}
}

func (dto *StrategyDTO) GetName() string {
	return dto.Name
}

func (dto *StrategyDTO) GetFilename() string {
	return dto.Filename
}

func (dto *StrategyDTO) GetVersion() string {
	return dto.Version
}
