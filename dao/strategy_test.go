// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/stretchr/testify/assert"
)

func TestStrategyDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	strategyDAO := NewStrategyDAO(ctx)

	strategies := []entity.Strategy{
		entity.Strategy{
			Name:     "DefaultTradingStrategy",
			Filename: "default.so",
			Version:  "0.0.1a"}}

	err := strategyDAO.Create(&strategies[0])
	assert.Equal(t, nil, err)

	persisted, exErr := strategyDAO.Find()
	assert.Equal(t, nil, exErr)

	assert.Equal(t, strategies[0].GetName(), persisted[0].GetName())
	assert.Equal(t, strategies[0].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, strategies[0].GetVersion(), persisted[0].GetVersion())

	CleanupIntegrationTest()
}

func TestStrategyDAO_Get(t *testing.T) {
	ctx := NewIntegrationTestContext()
	strategyDAO := NewStrategyDAO(ctx)

	strategies := []entity.Strategy{
		entity.Strategy{
			Name:     "TestStrategy1",
			Filename: "test.so",
			Version:  "0.0.1a"},
		entity.Strategy{
			Name:     "StrategyFindMe",
			Filename: "findme.so",
			Version:  "0.0.2a"},
		entity.Strategy{
			Name:     "TestStrategy2",
			Filename: "test2.so",
			Version:  "0.0.3a"}}

	for _, stratey := range strategies {
		err := strategyDAO.Create(&stratey)
		assert.Equal(t, nil, err)
	}

	persisted, exErr := strategyDAO.Get("StrategyFindMe")
	assert.Equal(t, nil, exErr)
	assert.NotNil(t, persisted)

	assert.Equal(t, strategies[1].GetName(), persisted.GetName())
	assert.Equal(t, strategies[1].GetFilename(), persisted.GetFilename())
	assert.Equal(t, strategies[1].GetVersion(), persisted.GetVersion())

	CleanupIntegrationTest()
}
