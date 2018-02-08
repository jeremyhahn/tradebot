package dto

import "github.com/jeremyhahn/tradebot/common"

type ChartIndicatorDTO struct {
	Id         uint   `json:"id"`
	ChartId    uint   `json:"chart_id"`
	Name       string `json:"name"`
	Parameters string `json:"parameters"`
	Filename   string `json:"filename"`
	common.ChartIndicator
}

func NewChartIndicatorDTO() common.ChartIndicator {
	return &ChartIndicatorDTO{}
}

func (chartIndicator *ChartIndicatorDTO) GetId() uint {
	return chartIndicator.Id
}

func (chartIndicator *ChartIndicatorDTO) GetChartId() uint {
	return chartIndicator.ChartId
}

func (chartIndicator *ChartIndicatorDTO) GetName() string {
	return chartIndicator.Name
}

func (chartIndicator *ChartIndicatorDTO) GetParameters() string {
	return chartIndicator.Parameters
}

func (chartIndicator *ChartIndicatorDTO) GetFilename() string {
	return chartIndicator.Filename
}
