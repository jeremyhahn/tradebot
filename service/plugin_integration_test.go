// +build integration

package service

import (
	"os"
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPluginService_GetWallets(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	wallets, err := pluginService.GetPlugins(common.WALLET_PLUGIN_TYPE)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(wallets))
	CleanupIntegrationTest()
}

func TestPluginService_GetBtcWallet(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	constructor, err := pluginService.CreateWallet("BTC")
	assert.Equal(t, nil, err)
	wallet := constructor(&common.WalletParams{
		Context: ctx,
		Address: os.Getenv("BTC_ADDRESS")})
	price := wallet.GetPrice()
	//util.DUMP(price)
	assert.Equal(t, true, price.GreaterThan(decimal.NewFromFloat(0)))
	CleanupIntegrationTest()
}

func TestPluginService_GetXrpWallet(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	constructor, err := pluginService.CreateWallet("XRP")
	assert.Equal(t, nil, err)
	wallet := constructor(&common.WalletParams{
		Context:          ctx,
		Address:          os.Getenv("XRP_ADDRESS"),
		MarketCapService: NewMarketCapService(ctx)})
	price := wallet.GetPrice()
	assert.Equal(t, true, price.GreaterThan(decimal.NewFromFloat(0)))
	CleanupIntegrationTest()
}

func TestPluginService_GetPlugins(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	exchanges, err := pluginService.GetPlugins(common.EXCHANGE_PLUGIN_TYPE)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(exchanges) == 4)
	CleanupIntegrationTest()
}

func TestPluginService_CreatePlugins(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	exchanges, err := pluginService.GetPlugins(common.EXCHANGE_PLUGIN_TYPE)
	for _, ex := range exchanges {
		constructor, err := pluginService.CreateExchange(ex)
		Exchange := constructor(ctx, &entity.UserCryptoExchange{
			Key:    "abc123",
			Secret: "$ecret!"})
		assert.Nil(t, err)
		assert.Equal(t, ex, Exchange.GetName())
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(exchanges) == 4)
	CleanupIntegrationTest()
}

func TestPluginService_Exchange_CreateCoinbase(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	pluginDAO.Create(&entity.Plugin{
		Name:     "Coinbase",
		Filename: "coinbase.so",
		Version:  "0.0.1a",
		Type:     common.EXCHANGE_PLUGIN_TYPE})
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	constructor, err := pluginService.CreateExchange("Coinbase")
	assert.Equal(t, nil, err)
	assert.NotNil(t, constructor)

	userExchangeEntity := &entity.UserCryptoExchange{
		Name:   "coinbase",
		Key:    "abc123",
		Secret: "$ecret!"}

	exchange := constructor(ctx, userExchangeEntity)
	assert.Equal(t, "Coinbase", exchange.GetDisplayName())

	CleanupIntegrationTest()
}

func TestPluginService_Exchange_ListExchanges(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	mapper := mapper.NewPluginMapper()
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper)
	plugins, err := pluginService.ListPlugins(common.EXCHANGE_PLUGIN_TYPE)
	assert.Nil(t, err)
	assert.Equal(t, true, len(plugins) == 4)
	assert.Equal(t, "binance", plugins[0])
	assert.Equal(t, "bittrex", plugins[1])
	assert.Equal(t, "coinbase", plugins[2])
	assert.Equal(t, "gdax", plugins[3])
	CleanupIntegrationTest()
}

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
	assert.Equal(t, decimal.NewFromFloat(15780.0).String(), sma.GetAverage().String())
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
	assert.Equal(t, decimal.NewFromFloat(15780.0).String(), ema.GetAverage().String())
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
	//	rsi := indicator.(indicators.RelativeStrengthIndex)
	//	assert.Equal(t, decimal.NewFromFloat(35.830474730988755).String(), rsi.Calculate(decimal.NewFromFloat(2000)).String())
	CleanupIntegrationTest()
}

/* BROKEN
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
	upper, middle, lower := bollinger.Calculate(decimal.NewFromFloat(1000))
	assert.Equal(t, decimal.NewFromFloat(4588.21).String(), upper.String())
	assert.Equal(t, decimal.NewFromFloat(3113.33).String(), middle.String())
	assert.Equal(t, decimal.NewFromFloat(1638.45).String(), lower.String())
	CleanupIntegrationTest()
}*/

/*
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
	macdValue, signal, histogram := macd.Calculate(decimal.NewFromFloat(1000))
	assert.Equal(t, decimal.NewFromFloat(730.7857256592952).String(), macdValue.String())
	assert.Equal(t, decimal.NewFromFloat(828.5704569933756).String(), signal.String())
	assert.Equal(t, decimal.NewFromFloat(-97.78473133408045).String(), histogram.String())
	CleanupIntegrationTest()
}*/

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
	assert.Equal(t, decimal.NewFromFloat(1.0).String(), obv.Calculate(decimal.NewFromFloat(12000)).String())
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
