// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/stretchr/testify/assert"
)

func TestPluginDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := NewPluginDAO(ctx)

	pluginDAO.Create(&entity.Plugin{
		Name:     "ExampleIndicator2",
		Filename: "example2.so",
		Version:  "0.0.1",
		Type:     "indicator"})

	persisted, err := pluginDAO.Get("ExampleIndicator2", "indicator")
	assert.Nil(t, err)
	assert.NotNil(t, persisted)

	CleanupIntegrationTest()
}

func TestPluginDAO_Strategy(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := NewPluginDAO(ctx)

	strategies := []entity.Plugin{
		entity.Plugin{
			Name:     "DefaultTradingStrategy2",
			Filename: "default2.so",
			Version:  "0.0.1a",
			Type:     "strategy"}}

	err := pluginDAO.Create(&strategies[0])
	assert.Equal(t, nil, err)

	persisted, exErr := pluginDAO.Find("strategy")
	assert.Equal(t, nil, exErr)
	assert.Equal(t, 1, len(persisted))

	assert.Equal(t, strategies[0].GetName(), persisted[0].GetName())
	assert.Equal(t, strategies[0].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, strategies[0].GetVersion(), persisted[0].GetVersion())
	assert.Equal(t, strategies[0].GetType(), persisted[0].GetType())

	CleanupIntegrationTest()
}

func TestPluginDAO_Strategy_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := NewPluginDAO(ctx)

	plugins := []entity.Plugin{
		entity.Plugin{
			Name:     "TestStrategy1",
			Filename: "test.so",
			Version:  "0.0.1a",
			Type:     "strategy"},
		entity.Plugin{
			Name:     "StrategyFindMe",
			Filename: "findme.so",
			Version:  "0.0.2a",
			Type:     "strategy"},
		entity.Plugin{
			Name:     "TestStrategy2",
			Filename: "test2.so",
			Version:  "0.0.3a",
			Type:     "strategy"}}

	for _, plugin := range plugins {
		err := pluginDAO.Create(&plugin)
		assert.Equal(t, nil, err)
	}

	persisted, exErr := pluginDAO.Get("StrategyFindMe", "strategy")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, plugins[1].GetName(), persisted.GetName())
	assert.Equal(t, plugins[1].GetFilename(), persisted.GetFilename())
	assert.Equal(t, plugins[1].GetVersion(), persisted.GetVersion())
	assert.Equal(t, plugins[1].GetType(), persisted.GetType())

	CleanupIntegrationTest()
}

func TestPluginDAO_Indicator(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := NewPluginDAO(ctx)

	indicators := []entity.Plugin{
		entity.Plugin{
			Name:     "RelativeStrengthIndex",
			Filename: "rsi.so",
			Version:  "0.0.1a",
			Type:     "indicator"},
		entity.Plugin{
			Name:     "BollingerBands",
			Filename: "bollinger_bands.so",
			Version:  "0.0.1a",
			Type:     "indicator"},
		entity.Plugin{
			Name:     "MovingAverageConvergenceDivergence",
			Filename: "macd.so",
			Version:  "0.0.1a",
			Type:     "indicator"}}

	err := pluginDAO.Create(&indicators[0])
	assert.Equal(t, nil, err)

	err = pluginDAO.Create(&indicators[1])
	assert.Equal(t, nil, err)

	err = pluginDAO.Create(&indicators[2])
	assert.Equal(t, nil, err)

	persisted, exErr := pluginDAO.Find("indicator") // order by name
	assert.Equal(t, nil, exErr)

	assert.Equal(t, indicators[1].GetName(), persisted[0].GetName())
	assert.Equal(t, indicators[1].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, indicators[1].GetVersion(), persisted[0].GetVersion())
	assert.Equal(t, indicators[1].GetType(), persisted[0].GetType())

	CleanupIntegrationTest()
}

func TestPluginDAO_Indicator_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := NewPluginDAO(ctx)

	indicators := []entity.Plugin{
		entity.Plugin{
			Name:     "TestIndicator",
			Filename: "test.so",
			Version:  "0.0.1a",
			Type:     "indicator"},
		entity.Plugin{
			Name:     "FindMe",
			Filename: "findme.so",
			Version:  "0.0.2a",
			Type:     "indicator"},
		entity.Plugin{
			Name:     "TestIndicator2",
			Filename: "test2.so",
			Version:  "0.0.3a",
			Type:     "indicator"}}

	for _, indicator := range indicators {
		err := pluginDAO.Create(&indicator)
		assert.Equal(t, nil, err)
	}

	persisted, exErr := pluginDAO.Get("FindMe", "indicator")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, indicators[1].GetName(), persisted.GetName())
	assert.Equal(t, indicators[1].GetFilename(), persisted.GetFilename())
	assert.Equal(t, indicators[1].GetVersion(), persisted.GetVersion())
	assert.Equal(t, indicators[1].GetType(), persisted.GetType())

	CleanupIntegrationTest()
}
