// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestPluginService_Indicator(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins", INDICATOR_PLUGIN)
	constructor, err := pluginService.GetIndicator("ExampleIndicator.so")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createChartCandles(), []string{"1", "2", "3"})
	assert.Equal(t, "ExampleIndicator", indicator.GetName())
	assert.Equal(t, "Example Indicator®", indicator.GetDisplayName())
	assert.Equal(t, []string{"1", "2", "3"}, indicator.GetDefaultParameters())
	assert.Equal(t, []string{"1", "2", "3"}, indicator.GetParameters())
	example := indicator.(example.ExampleIndicator)
	assert.Equal(t, 6.0, example.Calculate(5))
}
