package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
	"github.com/stretchr/testify/assert"
)

func TestExampleIndicator(t *testing.T) {
	var candles []common.Candlestick
	params := []string{"4", "5", "6"}
	exampleIndicator, err := CreateExampleIndicator(candles, params)
	assert.Equal(t, nil, err)
	example := exampleIndicator.(example.ExampleIndicator)
	assert.Equal(t, "ExampleIndicator", example.GetName())
	assert.Equal(t, "Example Indicator®", example.GetDisplayName())
	assert.Equal(t, params, example.GetParameters())
	assert.Equal(t, []string{"1", "2", "3"}, example.GetDefaultParameters())
	assert.Equal(t, 2.0, example.Calculate(1))
}
