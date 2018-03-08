// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/stretchr/testify/assert"
)

func TestPluginService_Indicator_ExampleIndicator(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "ExampleIndicator",
		Filename: "example.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("ExampleIndicator")
	assert.Equal(t, nil, err)
	assert.NotNil(t, constructor)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"4", "5", "6"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "ExampleIndicator", indicator.GetName())
	assert.Equal(t, "Example Indicator®", indicator.GetDisplayName())
	assert.Equal(t, []string{"4", "5", "6"}, indicator.GetParameters())
	assert.Equal(t, []string{"1", "2", "3"}, indicator.GetDefaultParameters())
	example := indicator.(example.ExampleIndicator)
	assert.Equal(t, 6.0, example.Calculate(5))
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_SimpleMovingAverage(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "SimpleMovingAverage",
		Filename: "sma.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("SimpleMovingAverage")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "SimpleMovingAverage", indicator.GetName())
	assert.Equal(t, "Simple Moving Average (SMA)", indicator.GetDisplayName())
	assert.Equal(t, []string{"5"}, indicator.GetParameters())
	assert.Equal(t, []string{"20"}, indicator.GetDefaultParameters())
	sma := indicator.(indicators.SimpleMovingAverage)
	assert.Equal(t, 15780.0, sma.GetAverage())
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_ExponentialMovingAverage(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "ExponentialMovingAverage",
		Filename: "ema.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("ExponentialMovingAverage")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "ExponentialMovingAverage", indicator.GetName())
	assert.Equal(t, "Exponential Moving Average (EMA)", indicator.GetDisplayName())
	assert.Equal(t, []string{"5"}, indicator.GetParameters())
	assert.Equal(t, []string{"20"}, indicator.GetDefaultParameters())
	ema := indicator.(indicators.ExponentialMovingAverage)
	assert.Equal(t, 15780.0, ema.GetAverage())
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_RelativeStrengthIndex(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"14", "80", "20"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", indicator.GetName())
	assert.Equal(t, "Relative Strength Index (RSI)", indicator.GetDisplayName())
	assert.Equal(t, []string{"14", "80.000000", "20.000000"}, indicator.GetParameters())
	assert.Equal(t, []string{"14", "70", "30"}, indicator.GetDefaultParameters())
	rsi := indicator.(indicators.RelativeStrengthIndex)
	assert.Equal(t, 35.830474730988755, rsi.Calculate(2000))
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_BollingerBands(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("BollingerBands")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"15", "2"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "BollingerBands", indicator.GetName())
	assert.Equal(t, "Bollinger Bands®", indicator.GetDisplayName())
	assert.Equal(t, []string{"15", "2.000000"}, indicator.GetParameters())
	assert.Equal(t, []string{"20", "2"}, indicator.GetDefaultParameters())
	bollinger := indicator.(indicators.BollingerBands)
	upper, middle, lower := bollinger.Calculate(1000)
	assert.Equal(t, 4588.21, upper)
	assert.Equal(t, 3113.33, middle)
	assert.Equal(t, 1638.45, lower)
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_MACD(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("MovingAverageConvergenceDivergence")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"10", "24", "9"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "MovingAverageConvergenceDivergence", indicator.GetName())
	assert.Equal(t, "Moving Average Convergence Divergence (MACD)", indicator.GetDisplayName())
	assert.Equal(t, []string{"10", "24", "9"}, indicator.GetParameters())
	assert.Equal(t, []string{"12", "26", "9"}, indicator.GetDefaultParameters())
	macd := indicator.(indicators.MovingAverageConvergenceDivergence)
	macdValue, signal, histogram := macd.Calculate(1000)
	assert.Equal(t, 730.7857256592952, macdValue)
	assert.Equal(t, 828.5704569933756, signal)
	assert.Equal(t, -97.78473133408045, histogram)
	CleanupIntegrationTest()
}

func TestPluginService_Indicator_OnBalanceVolume(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "OnBalanceVolume",
		Filename: "obv.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateIndicator("OnBalanceVolume")
	assert.Equal(t, nil, err)
	indicator, err := constructor(createIntegrationTestCandles(), []string{"5"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "OnBalanceVolume", indicator.GetName())
	assert.Equal(t, "On Balance Volume (OBV)", indicator.GetDisplayName())
	assert.Equal(t, []string{}, indicator.GetParameters())
	assert.Equal(t, []string{}, indicator.GetDefaultParameters())
	obv := indicator.(indicators.OnBalanceVolume)
	assert.Equal(t, 1.0, obv.Calculate(12000))
	CleanupIntegrationTest()
}

func TestPluginService_Strategy_DefaultTradingStrategy(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)

	dao.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	dao.Create(&entity.Plugin{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	dao.Create(&entity.Plugin{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	dao.Create(&entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE})

	constructor, err := pluginService.CreateStrategy("DefaultTradingStrategy")
	assert.Equal(t, nil, err)

	rsiConstructor, err := pluginService.CreateIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	rsiIndicator, err := rsiConstructor(createIntegrationTestCandles(), nil)
	assert.Equal(t, nil, err)
	rsi := rsiIndicator.(indicators.RelativeStrengthIndex)

	bollingerConstructor, err := pluginService.CreateIndicator("BollingerBands")
	assert.Equal(t, nil, err)
	bollingerIndicator, err := bollingerConstructor(createIntegrationTestCandles(), nil)
	assert.Equal(t, nil, err)
	bbands := bollingerIndicator.(indicators.BollingerBands)

	macdConstructor, err := pluginService.CreateIndicator("MovingAverageConvergenceDivergence")
	assert.Equal(t, nil, err)
	macdIndicator, err := macdConstructor(createIntegrationTestCandles(), nil)
	assert.Equal(t, nil, err)
	macd := macdIndicator.(indicators.MovingAverageConvergenceDivergence)

	indicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              rsi,
		"BollingerBands":                     bbands,
		"MovingAverageConvergenceDivergence": macd}

	strategyParams := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"},
		Indicators: indicators}
	strategy, err := constructor(strategyParams)
	assert.Equal(t, nil, err)
	assert.NotNilf(t, strategy, "Failed to instantiate strategy: %s", "DefaultTradingStrategy")
	requiredIndicators := []string{
		"RelativeStrengthIndex",
		"BollingerBands",
		"MovingAverageConvergenceDivergence"}

	assert.Equal(t, requiredIndicators, strategy.GetRequiredIndicators())
	assert.Equal(t, strategyParams, strategy.GetParameters())
	CleanupIntegrationTest()
}

func TestPluginService_Strategy_DefaultTradingStrategy_MissingRsiIndicators(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateStrategy("DefaultTradingStrategy")
	assert.Equal(t, nil, err)
	strategyParams := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"}}
	strategy, err := constructor(strategyParams)
	assert.Equal(t, "Strategy requires missing indicator: RelativeStrengthIndex", err.Error())
	assert.Equal(t, nil, strategy)
	CleanupIntegrationTest()
}

func TestPluginService_Strategy_DefaultTradingStrategy_InvalidConfiguration(t *testing.T) {
	ctx := NewIntegrationTestContext()
	dao := dao.NewPluginDAO(ctx)
	dao.Create(&entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", dao, mapper)
	constructor, err := pluginService.CreateStrategy("DefaultTradingStrategy")
	assert.Equal(t, nil, err)
	strategyParams := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"},
		Config: []string{"foo", "bar"}}
	strategy, err := constructor(strategyParams)
	assert.Equal(t, "Invalid configuration. Expected 8 items, received 2 (foo,bar)", err.Error())
	assert.Equal(t, nil, strategy)
	CleanupIntegrationTest()
}
