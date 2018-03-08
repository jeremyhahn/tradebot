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

func TestGetIndicator_Success(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	pluginDAO.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	pluginMapper := mapper.NewPluginMapper()
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)

	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)
	rsi, err := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())
	assert.Equal(t, common.INDICATOR_PLUGIN_TYPE, rsi.GetType())

	CleanupIntegrationTest()
}

func TestGetIndicator_SuccessfulLoadingTwice(t *testing.T) {
	ctx := NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	assert.NotNil(t, pluginDAO)
	pluginDAO.Find(common.INDICATOR_PLUGIN_TYPE)

	pluginDAO.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	pluginMapper := mapper.NewPluginMapper()
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)

	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)

	rsi, err := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())
	assert.Equal(t, common.INDICATOR_PLUGIN_TYPE, rsi.GetType())

	rsi2, err2 := indicatorService.GetIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err2)
	assert.Equal(t, "RelativeStrengthIndex", rsi2.GetName())
	assert.Equal(t, "rsi.so", rsi2.GetFilename())
	assert.Equal(t, "0.0.1a", rsi2.GetVersion())
	assert.Equal(t, common.INDICATOR_PLUGIN_TYPE, rsi2.GetType())

	CleanupIntegrationTest()
}

func TestGetIndicator_IndicatorDoesntExistInDatabase(t *testing.T) {
	ctx := NewIntegrationTestContext()
	pluginDAO := dao.NewPluginDAO(ctx)
	pluginMapper := mapper.NewPluginMapper()
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)
	doesntExistIndicator, err := indicatorService.GetIndicator("IndicatorDoesntExist")
	assert.NotNil(t, err)
	assert.Equal(t, nil, doesntExistIndicator)
	assert.Equal(t, "IndicatorDoesntExist (indicator) plugin not found in database", err.Error())
	CleanupIntegrationTest()
}

func TestGetIndicator_IndicatorExists(t *testing.T) {
	ctx := NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	pluginMapper := mapper.NewPluginMapper()
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)

	indicator := &entity.Plugin{
		Name:     "IndicatorExists",
		Filename: "fake.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE}
	pluginDAO.Create(indicator)

	doesntExistIndicator, err := indicatorService.GetIndicator("IndicatorExists")
	assert.Nil(t, err)
	assert.NotNil(t, doesntExistIndicator)
	assert.Equal(t, indicator.GetName(), doesntExistIndicator.GetName())
	assert.Equal(t, indicator.GetFilename(), doesntExistIndicator.GetFilename())
	assert.Equal(t, indicator.GetVersion(), doesntExistIndicator.GetVersion())
	assert.Equal(t, indicator.GetType(), doesntExistIndicator.GetType())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_GetIndicator(t *testing.T) {
	ctx := NewIntegrationTestContext()

	chartDAO := dao.NewChartDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	pluginMapper := mapper.NewPluginMapper()
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, pluginMapper)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)

	pluginDAO.Create(&entity.Plugin{
		Name:     "RelativeStrengthIndex",
		Filename: "rsi.so",
		Version:  "0.0.1a",
		Type:     common.INDICATOR_PLUGIN_TYPE})

	chartEntity := createIntegrationTestChart(ctx)
	chartDAO.Create(chartEntity)

	indicators := chartEntity.GetIndicators()

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

	pluginDAO := dao.NewPluginDAO(ctx)
	pluginMapper := mapper.NewPluginMapper()

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

	pluginService := CreatePluginService(ctx, "../plugins", pluginDAO, pluginMapper)
	candles := createIntegrationTestCandles()

	chartMapper := mapper.NewChartMapper(ctx)
	chartDTO := chartMapper.MapChartEntityToDto(chartEntity)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	indicatorService := NewIndicatorService(ctx, chartIndicatorDAO, pluginService)

	financialIndicators, err := indicatorService.GetChartIndicators(chartDTO, candles)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, len(financialIndicators))
	assert.Equal(t, "RelativeStrengthIndex", financialIndicators["RelativeStrengthIndex"].GetName())
	assert.Equal(t, "BollingerBands", financialIndicators["BollingerBands"].GetName())
	assert.Equal(t, "MovingAverageConvergenceDivergence", financialIndicators["MovingAverageConvergenceDivergence"].GetName())

	CleanupIntegrationTest()
}
