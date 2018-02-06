// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChartStrategyDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()

	chartDAO := NewChartDAO(ctx)
	chartStrategyDAO := NewChartStrategyDAO(ctx)

	chartStrategies := []ChartStrategy{
		ChartStrategy{
			ChartId:    1,
			Name:       "DefaultTradingStrategy",
			Parameters: "1,2,3"}}

	chart := &Chart{
		UserId:    ctx.User.Id,
		Base:      "BTC",
		Quote:     "USD",
		Exchange:  "gdax",
		Period:    900,
		AutoTrade: 1}

	chartDAO.Create(chart)

	err := chartStrategyDAO.Create(&chartStrategies[0])
	persistedStrategies, err := chartStrategyDAO.Find(chart)
	assert.Equal(t, nil, err)

	assert.Equal(t, chartStrategies[0].GetId(), persistedStrategies[0].GetId())
	assert.Equal(t, chartStrategies[0].GetChartId(), persistedStrategies[0].GetChartId())
	assert.Equal(t, chartStrategies[0].GetName(), persistedStrategies[0].GetName())
	assert.Equal(t, chartStrategies[0].GetParameters(), persistedStrategies[0].GetParameters())

	CleanupIntegrationTest()
}

func TestChartStrategyDAO_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)
	userStrategyDAO := NewChartStrategyDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	indicators := chart.GetStrategies()

	err := chartDAO.Create(chart)
	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), chart.GetId())

	persisted, exErr := userStrategyDAO.Get(chart, "DefaultTradingStrategy")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, indicators[0].GetId(), persisted.GetId())
	assert.Equal(t, indicators[0].GetChartId(), persisted.GetChartId())
	assert.Equal(t, indicators[0].GetName(), persisted.GetName())
	assert.Equal(t, indicators[0].GetParameters(), persisted.GetParameters())

	CleanupIntegrationTest()
}
