// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestGetIndicator_Success(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	indicatorDAO.Create(&dao.Indicator{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a"})

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	rsi, err := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())

	CleanupIntegrationTest()
}

func TestGetIndicator_SuccessfulLoadingTwice(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	assert.NotNil(t, indicatorDAO)
	indicatorDAO.Find()

	indicatorDAO.Create(&dao.Indicator{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a"})

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	rsi, err := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())

	rsi2, err2 := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err2)
	assert.Equal(t, "RelativeStrengthIndex", rsi2.GetName())
	assert.Equal(t, "rsi.so", rsi2.GetFilename())
	assert.Equal(t, "0.0.1a", rsi2.GetVersion())

	CleanupIntegrationTest()
}

func TestGetIndicator_IndicatorDoesntExistInDatabase(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := dao.NewIndicatorDAO(ctx)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	doesntExistIndicator, err := indicatorService.GetIndicator("IndicatorDoesntExist")
	assert.NotNil(t, err)
	assert.Equal(t, nil, doesntExistIndicator)
	assert.Equal(t, "Failed to get platform indicator: IndicatorDoesntExist", err.Error())

	CleanupIntegrationTest()
}

func TestGetIndicator_IndicatorDoesntExist(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	indicator := &dao.Indicator{
		Name:     "IndicatorDoesntExist",
		Filename: "fake.so",
		Version:  "0.0.1a"}
	indicatorDAO.Create(indicator)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	doesntExistIndicator, err := indicatorService.GetIndicator("IndicatorDoesntExist")
	assert.Nil(t, err)
	assert.NotNil(t, doesntExistIndicator)
	assert.Equal(t, indicator.GetName(), doesntExistIndicator.GetName())
	assert.Equal(t, indicator.GetFilename(), doesntExistIndicator.GetFilename())
	assert.Equal(t, indicator.GetVersion(), doesntExistIndicator.GetVersion())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_GetIndicator(t *testing.T) {
	ctx := NewIntegrationTestContext()

	chartDAO := dao.NewChartDAO(ctx)
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	indicatorDAO.Create(&dao.Indicator{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a"})

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	indicators := chartEntity.GetIndicators()

	pluginService := CreatePluginService(ctx, "../plugins")
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	candles := createIntegrationTestCandles()
	chartIndicator, err := indicatorService.GetChartIndicator(chartDTO, "RelativeStrengthIndex", candles)

	assert.Nil(t, err)
	assert.NotNil(t, chartIndicator)
	assert.Equal(t, indicators[0].GetName(), chartIndicator.GetName())
	//assert.Equal(t, indicators[0].GetParameters(), strings.Join(chartIndicator.GetParameters(), ","))
	assert.Equal(t, indicators[0].GetChartId(), chartEntity.GetId())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_GetIndicators(t *testing.T) {
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

	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	financialIndicators, err := indicatorService.GetChartIndicators(chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))
	assert.Equal(t, "RelativeStrengthIndex", financialIndicators["RelativeStrengthIndex"].GetName())
	assert.Equal(t, "BollingerBands", financialIndicators["BollingerBands"].GetName())
	assert.Equal(t, "MovingAverageConvergenceDivergence", financialIndicators["MovingAverageConvergenceDivergence"].GetName())

	CleanupIntegrationTest()
}
