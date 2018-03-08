// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestGetPlatformStrategy_GetPlatformStrategy(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	strategyEntity := &entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE}
	pluginDAO.Create(strategyEntity)

	pluginService := NewPluginService(ctx, pluginDAO, mapper.NewPluginMapper())

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	strategyService := NewStrategyService(ctx, chartStrategyDAO, pluginService, indicatorService, chartMapper)

	defaultTradingStrategy, err := strategyService.GetStrategy("DefaultTradingStrategy")
	assert.Equal(t, nil, err)
	assert.Equal(t, strategyEntity.GetName(), defaultTradingStrategy.GetName())
	assert.Equal(t, strategyEntity.GetFilename(), defaultTradingStrategy.GetFilename())
	assert.Equal(t, strategyEntity.GetVersion(), defaultTradingStrategy.GetVersion())

	CleanupIntegrationTest()
}

func TestGetChartStrategy_GetChartStrategy(t *testing.T) {
	ctx := NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	pluginDAO.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	strategyEntity := &entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE}
	pluginDAO.Create(strategyEntity)

	chartDAO := dao.NewChartDAO(ctx)
	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper.NewPluginMapper())
	candles := createIntegrationTestCandles()

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)
	financialIndicators, err := indicatorService.GetChartIndicators(chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))

	strategyService := NewStrategyService(ctx, chartStrategyDAO, pluginService, indicatorService, chartMapper)

	defaultTradingStrategy, err := strategyService.GetChartStrategy(chartDTO, "DefaultTradingStrategy", candles)
	assert.Equal(t, nil, err)
	assert.Equal(t, defaultTradingStrategy.GetRequiredIndicators(),
		[]string{"RelativeStrengthIndex", "BollingerBands", "MovingAverageConvergenceDivergence"})

	CleanupIntegrationTest()
}

func TestGetChartStrategy_GetChartStrategies(t *testing.T) {
	ctx := NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	pluginDAO.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	strategyEntity := &entity.Plugin{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a",
		Type:     common.STRATEGY_PLUGIN_TYPE}
	pluginDAO.Create(strategyEntity)

	chartDAO := dao.NewChartDAO(ctx)
	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, mapper.NewPluginMapper())
	candles := createIntegrationTestCandles()

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)
	financialIndicators, err := indicatorService.GetChartIndicators(chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))

	strategyService := NewStrategyService(ctx, chartStrategyDAO, pluginService, indicatorService, chartMapper)

	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          chartEntity.GetBase(),
			Quote:         chartEntity.GetQuote(),
			LocalCurrency: ctx.GetUser().GetLocalCurrency()},
		Indicators: financialIndicators}

	tradingStrategies, err := strategyService.GetChartStrategies(chartDTO, params, candles)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(tradingStrategies))
	assert.Equal(t, tradingStrategies[0].GetRequiredIndicators(),
		[]string{"RelativeStrengthIndex", "BollingerBands", "MovingAverageConvergenceDivergence"})

	CleanupIntegrationTest()
}
