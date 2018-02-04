package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/stretchr/testify/assert"
)

func TestExampleIndicator(t *testing.T) {
	var candles []common.Candlestick
	params := []string{"4", "5", "6"}
	exampleIndicator := CreateExampleIndicator(candles, params)
	assert.Equal(t, "ExampleIndicator", exampleIndicator.GetName())
	assert.Equal(t, "Example Indicator®", exampleIndicator.GetDisplayName())
	assert.Equal(t, params, exampleIndicator.GetParameters())
	assert.Equal(t, []string{"1", "2", "3"}, exampleIndicator.GetDefaultParameters())
	assert.Equal(t, 2.0, exampleIndicator.Calculate(1))
}
