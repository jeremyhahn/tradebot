// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestGetPlatformStrategy_GetPlatformStrategy(t *testing.T) {
	ctx := NewIntegrationTestContext()
	strategyDAO := dao.NewStrategyDAO(ctx)
	strategyEntity := &dao.Strategy{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a"}
	strategyDAO.Create(strategyEntity)

	pluginService := NewPluginService(ctx)

	indicatorDAO := dao.NewIndicatorDAO(ctx)
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	strategyMapper := mapper.NewStrategyMapper()
	strategyService := NewStrategyService(ctx, strategyDAO, chartStrategyDAO, pluginService, indicatorService, chartMapper, strategyMapper)

	defaultTradingStrategy, err := strategyService.GetPlatformStrategy("DefaultTradingStrategy")
	assert.Equal(t, nil, err)
	assert.Equal(t, strategyEntity.GetName(), defaultTradingStrategy.GetName())
	assert.Equal(t, strategyEntity.GetFilename(), defaultTradingStrategy.GetFilename())
	assert.Equal(t, strategyEntity.GetVersion(), defaultTradingStrategy.GetVersion())

	CleanupIntegrationTest()
}

func TestGetChartStrategy_GetChartStrategy(t *testing.T) {
	ctx := NewIntegrationTestContext()

	indicatorDAO := dao.NewIndicatorDAO(ctx)
	indicatorDAO.Create(&dao.Indicator{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a"})
	indicatorDAO.Create(&dao.Indicator{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a"})
	indicatorDAO.Create(&dao.Indicator{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a"})

	strategyDAO := dao.NewStrategyDAO(ctx)
	strategyEntity := &dao.Strategy{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a"}
	strategyDAO.Create(strategyEntity)

	chartDAO := dao.NewChartDAO(ctx)
	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	pluginService := CreatePluginService(ctx, "../plugins")
	candles := createIntegrationTestCandles()

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)
	financialIndicators, err := indicatorService.GetChartIndicators(&chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))

	strategyMapper := mapper.NewStrategyMapper()
	strategyService := NewStrategyService(ctx, strategyDAO, chartStrategyDAO, pluginService, indicatorService, chartMapper, strategyMapper)

	defaultTradingStrategy, err := strategyService.GetChartStrategy(&chartDTO, "DefaultTradingStrategy", candles)
	assert.Equal(t, nil, err)
	assert.Equal(t, defaultTradingStrategy.GetRequiredIndicators(),
		[]string{"RelativeStrengthIndex", "BollingerBands", "MovingAverageConvergenceDivergence"})

	CleanupIntegrationTest()
}

func TestGetChartStrategy_GetChartStrategies(t *testing.T) {
	ctx := NewIntegrationTestContext()

	indicatorDAO := dao.NewIndicatorDAO(ctx)
	indicatorDAO.Create(&dao.Indicator{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a"})
	indicatorDAO.Create(&dao.Indicator{
		Name:     "BollingerBands",
		Filename: "bollinger_bands.so",
		Version:  "0.0.1a"})
	indicatorDAO.Create(&dao.Indicator{
		Name:     "MovingAverageConvergenceDivergence",
		Filename: "macd.so",
		Version:  "0.0.1a"})

	strategyDAO := dao.NewStrategyDAO(ctx)
	strategyEntity := &dao.Strategy{
		Name:     "DefaultTradingStrategy",
		Filename: "default.so",
		Version:  "0.0.1a"}
	strategyDAO.Create(strategyEntity)

	chartDAO := dao.NewChartDAO(ctx)
	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	pluginService := CreatePluginService(ctx, "../plugins")
	candles := createIntegrationTestCandles()

	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)
	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)
	financialIndicators, err := indicatorService.GetChartIndicators(&chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))

	strategyMapper := mapper.NewStrategyMapper()
	strategyService := NewStrategyService(ctx, strategyDAO, chartStrategyDAO, pluginService, indicatorService, chartMapper, strategyMapper)

	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          chartEntity.GetBase(),
			Quote:         chartEntity.GetQuote(),
			LocalCurrency: ctx.User.LocalCurrency},
		Indicators: financialIndicators}

	tradingStrategies, err := strategyService.GetChartStrategies(&chartDTO, params, candles)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(tradingStrategies))
	assert.Equal(t, tradingStrategies[0].GetRequiredIndicators(),
		[]string{"RelativeStrengthIndex", "BollingerBands", "MovingAverageConvergenceDivergence"})

	CleanupIntegrationTest()
}
