// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndicatorDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)
	indicatorDAO := NewIndicatorDAO(ctx)

	chart, indicators, _ := createIntegrationTestChart(ctx)

	err := chartDAO.Create(chart)
	assert.Equal(t, nil, err)

	persisted, exErr := indicatorDAO.Find(chart)
	assert.Equal(t, nil, exErr)

	assert.Equal(t, indicators[0].GetId(), persisted[0].GetId())
	assert.Equal(t, indicators[0].GetChartId(), persisted[0].GetChartId())
	assert.Equal(t, indicators[0].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, indicators[0].GetName(), persisted[0].GetName())
	assert.Equal(t, indicators[0].GetParameters(), persisted[0].GetParameters())

	CleanupIntegrationTest()
}
