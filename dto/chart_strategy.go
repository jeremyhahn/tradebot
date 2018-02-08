package dto

import "github.com/jeremyhahn/tradebot/common"

type ChartStrategyDTO struct {
	Id         uint   `json:"id"`
	ChartId    uint   `json:"chart_id"`
	Name       string `json:"name"`
	Parameters string `json:"parameters"`
	Filename   string `json:"filename"`
	common.ChartStrategy
}

func NewChartStrategyDTO() common.ChartStrategy {
	return &ChartStrategyDTO{}
}

func (chartStrategy *ChartStrategyDTO) GetId() uint {
	return chartStrategy.Id
}

func (chartStrategy *ChartStrategyDTO) GetChartId() uint {
	return chartStrategy.ChartId
}

func (chartStrategy *ChartStrategyDTO) GetName() string {
	return chartStrategy.Name
}

func (chartStrategy *ChartStrategyDTO) GetParameters() string {
	return chartStrategy.Parameters
}

func (chartStrategy *ChartStrategyDTO) GetFilename() string {
	return chartStrategy.Filename
}
