// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategyDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	strategyDAO := NewStrategyDAO(ctx)

	strategies := []Strategy{
		Strategy{
			Name:     "DefaultTradingStrategy",
			Filename: "default.so",
			Version:  "0.0.1a"}}

	err := strategyDAO.Create(&strategies[0])
	assert.Equal(t, nil, err)

	persisted, exErr := strategyDAO.Find()
	assert.Equal(t, nil, exErr)

	assert.Equal(t, strategies[0].GetId(), persisted[0].GetId())
	assert.Equal(t, strategies[0].GetName(), persisted[0].GetName())
	assert.Equal(t, strategies[0].GetFilename(), persisted[0].GetFilename())
	assert.Equal(t, strategies[0].GetVersion(), persisted[0].GetVersion())

	CleanupIntegrationTest()
}
