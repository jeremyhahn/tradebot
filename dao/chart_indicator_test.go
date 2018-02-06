// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChartIndicatorDAO_CreateFind(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)
	userIndicatorDAO := NewChartIndicatorDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	indicators := chart.GetIndicators()

	err := chartDAO.Create(chart)
	assert.Equal(t, nil, err)

	persisted, exErr := userIndicatorDAO.Find(chart)
	assert.Equal(t, nil, exErr)

	assert.Equal(t, indicators[0].GetId(), persisted[0].GetId())
	assert.Equal(t, indicators[0].GetChartId(), persisted[0].GetChartId())
	assert.Equal(t, indicators[0].GetName(), persisted[0].GetName())
	assert.Equal(t, indicators[0].GetParameters(), persisted[0].GetParameters())

	CleanupIntegrationTest()
}

func TestChartIndicatorDAO_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)
	userIndicatorDAO := NewChartIndicatorDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	indicators := chart.GetIndicators()

	err := chartDAO.Create(chart)
	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), chart.GetId())

	persisted, exErr := userIndicatorDAO.Get(chart, "RelativeStrengthIndex")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, indicators[0].GetId(), persisted.GetId())
	assert.Equal(t, indicators[0].GetChartId(), persisted.GetChartId())
	assert.Equal(t, indicators[0].GetName(), persisted.GetName())
	assert.Equal(t, indicators[0].GetParameters(), persisted.GetParameters())

	CleanupIntegrationTest()
}
