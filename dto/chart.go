package dto

import (
	"encoding/json"

	"github.com/jeremyhahn/tradebot/common"
)

type ChartDTO struct {
	Id         uint                    `json:"id"`
	Base       string                  `json:"base"`
	Quote      string                  `json:"quote"`
	Exchange   string                  `json:"exchange"`
	Period     int                     `json:"period"`
	Price      float64                 `json:"price"`
	AutoTrade  uint                    `json:"autotrade"`
	Indicators []common.ChartIndicator `json:"indicators"`
	Strategies []common.ChartStrategy  `json:"strategies"`
	Trades     []common.Trade          `json:"trades"`
	common.Chart
}

func NewChartDTO() common.Chart {
	return &ChartDTO{}
}

func (chart ChartDTO) GetId() uint {
	return chart.Id
}

func (chart ChartDTO) GetBase() string {
	return chart.Base
}

func (chart ChartDTO) GetQuote() string {
	return chart.Quote
}

func (chart ChartDTO) GetExchange() string {
	return chart.Exchange
}

func (chart ChartDTO) GetPeriod() int {
	return chart.Period
}

func (chart ChartDTO) GetPrice() float64 {
	return chart.Price
}

func (chart ChartDTO) GetAutoTrade() uint {
	return chart.AutoTrade
}

func (chart ChartDTO) IsAutoTrade() bool {
	return chart.AutoTrade == 1
}

func (chart ChartDTO) GetIndicators() []common.ChartIndicator {
	return chart.Indicators
}

func (chart ChartDTO) GetStrategies() []common.ChartStrategy {
	return chart.Strategies
}

func (chart ChartDTO) GetTrades() []common.Trade {
	return chart.Trades
}

func (chart ChartDTO) ToJSON() (string, error) {
	jsonData, err := json.Marshal(chart)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
