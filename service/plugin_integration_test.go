package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestPluginService_Indicator_ExampleIndicator(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("example.so", "ExampleIndicator")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"4", "5", "6"})
	assert.Equal(t, "ExampleIndicator", indicator.GetName())
	assert.Equal(t, "Example Indicator®", indicator.GetDisplayName())
	assert.Equal(t, []string{"4", "5", "6"}, indicator.GetParameters())
	assert.Equal(t, []string{"1", "2", "3"}, indicator.GetDefaultParameters())
	example := indicator.(example.ExampleIndicator)
	assert.Equal(t, 6.0, example.Calculate(5))
}

func TestPluginService_Indicator_SimpleMovingAverage(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("sma.so", "SimpleMovingAverage")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, "SimpleMovingAverage", indicator.GetName())
	assert.Equal(t, "Simple Moving Average (SMA)", indicator.GetDisplayName())
	assert.Equal(t, []string{"5"}, indicator.GetParameters())
	assert.Equal(t, []string{"20"}, indicator.GetDefaultParameters())
	sma := indicator.(indicators.SimpleMovingAverage)
	assert.Equal(t, 15780.0, sma.GetAverage())
}

func TestPluginService_Indicator_ExponentialMovingAverage(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("ema.so", "ExponentialMovingAverage")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, "ExponentialMovingAverage", indicator.GetName())
	assert.Equal(t, "Exponential Moving Average (EMA)", indicator.GetDisplayName())
	assert.Equal(t, []string{"5"}, indicator.GetParameters())
	assert.Equal(t, []string{"20"}, indicator.GetDefaultParameters())
	ema := indicator.(indicators.ExponentialMovingAverage)
	assert.Equal(t, 15780.0, ema.GetAverage())
}

func TestPluginService_Indicator_RelativeStrengthIndex(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("rsi.so", "RelativeStrengthIndex")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"1", "2", "3"})
	assert.Equal(t, "RelativeStrengthIndex", indicator.GetName())
	assert.Equal(t, "Relative Strength Index (RSI)", indicator.GetDisplayName())
	assert.Equal(t, []string{"1", "2.000000", "3.000000"}, indicator.GetParameters())
	assert.Equal(t, []string{"14", "70", "30"}, indicator.GetDefaultParameters())
	rsi := indicator.(indicators.RelativeStrengthIndex)
	assert.Equal(t, 56.52173913043478, rsi.Calculate(1000))
}

func TestPluginService_Indicator_BollingerBands(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("bollinger_bands.so", "BollingerBands")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"15", "2"})
	assert.Equal(t, "BollingerBands", indicator.GetName())
	assert.Equal(t, "Bollinger Bands®", indicator.GetDisplayName())
	assert.Equal(t, []string{"15", "2.000000"}, indicator.GetParameters())
	assert.Equal(t, []string{"20", "2"}, indicator.GetDefaultParameters())
	bollinger := indicator.(indicators.BollingerBands)
	upper, middle, lower := bollinger.Calculate(1000)
	assert.Equal(t, 1642.48, upper)
	assert.Equal(t, 860.0, middle)
	assert.Equal(t, 77.52, lower)
}

func TestPluginService_Indicator_MACD(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("macd.so", "MovingAverageConvergenceDivergence")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"10", "24", "9"})
	assert.Equal(t, "MovingAverageConvergenceDivergence", indicator.GetName())
	assert.Equal(t, "Moving Average Convergence Divergence (MACD)", indicator.GetDisplayName())
	assert.Equal(t, []string{"10", "24", "9"}, indicator.GetParameters())
	assert.Equal(t, []string{"12", "26", "9"}, indicator.GetDefaultParameters())
	macd := indicator.(indicators.MovingAverageConvergenceDivergence)
	macdValue, signal, histogram := macd.Calculate(1000)
	assert.Equal(t, 730.7857256592952, macdValue)
	assert.Equal(t, 828.5704569933756, signal)
	assert.Equal(t, -97.78473133408045, histogram)
}

func TestPluginService_Indicator_OnBalanceVolume(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetIndicator("obv.so", "OnBalanceVolume")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	indicator := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, "OnBalanceVolume", indicator.GetName())
	assert.Equal(t, "On Balance Volume (OBV)", indicator.GetDisplayName())
	assert.Equal(t, []string{}, indicator.GetParameters())
	assert.Equal(t, []string{}, indicator.GetDefaultParameters())
	obv := indicator.(indicators.OnBalanceVolume)
	assert.Equal(t, 1.0, obv.Calculate(12000))
}

func TestPluginService_Strategy_DefaultTradingStrategy(t *testing.T) {
	ctx := test.NewUnitTestContext()
	pluginService := CreatePluginService(ctx, "../plugins")
	constructor, err := pluginService.GetStrategy("default.so", "DefaultTradingStrategy")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	strategy := constructor(createIntegrationTestCandles(), []string{"5"})
	requiredIndicators := []string{
		"RelativeStrengthIndex",
		"BollingerBands",
		"MovingAverageConvergenceDivergence"}
	assert.Equal(t, requiredIndicators, strategy.GetRequiredIndicators())
	assert.Equal(t, "On Balance Volume (OBV)", indicator.GetDisplayName())
	assert.Equal(t, []string{}, indicator.GetParameters())
	assert.Equal(t, []string{}, indicator.GetDefaultParameters())
	obv := indicator.(indicators.OnBalanceVolume)
	assert.Equal(t, 1.0, obv.Calculate(12000))
}
