// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/stretchr/testify/assert"
)

func TestIndicatorDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := NewIndicatorDAO(ctx)

	indicators := []entity.Indicator{
		entity.Indicator{
			Name:     "RelativeStrengthIndex",
			Filename: "rsi.so",
			Version:  "0.0.1a"},
		entity.Indicator{
			Name:     "BollingerBands",
			Filename: "bollinger_bands.so",
			Version:  "0.0.1a"},
		entity.Indicator{
			Name:     "MovingAverageConvergenceDivergence",
			Filename: "macd.so",
			Version:  "0.0.1a"}}

	err := indicatorDAO.Create(&indicators[0])
	assert.Equal(t, nil, err)

	err = indicatorDAO.Create(&indicators[1])
	assert.Equal(t, nil, err)

	err = indicatorDAO.Create(&indicators[2])
	assert.Equal(t, nil, err)

	persisted, exErr := indicatorDAO.Find() // order by name
	assert.Equal(t, nil, exErr)

	assert.Equal(t, indicators[1].GetName(), persisted[0].GetName())
	assert.Equal(t, indicators[1].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, indicators[1].GetVersion(), persisted[0].GetVersion())

	CleanupIntegrationTest()
}

func TestIndicatorDAO_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := NewIndicatorDAO(ctx)

	indicators := []entity.Indicator{
		entity.Indicator{
			Name:     "TestIndicator",
			Filename: "test.so",
			Version:  "0.0.1a"},
		entity.Indicator{
			Name:     "FindMe",
			Filename: "findme.so",
			Version:  "0.0.2a"},
		entity.Indicator{
			Name:     "TestIndicator2",
			Filename: "test2.so",
			Version:  "0.0.3a"}}

	for _, indicator := range indicators {
		err := indicatorDAO.Create(&indicator)
		assert.Equal(t, nil, err)
	}

	persisted, exErr := indicatorDAO.Get("FindMe")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, indicators[1].GetName(), persisted.GetName())
	assert.Equal(t, indicators[1].GetFilename(), persisted.GetFilename())
	assert.Equal(t, indicators[1].GetVersion(), persisted.GetVersion())

	CleanupIntegrationTest()
}
