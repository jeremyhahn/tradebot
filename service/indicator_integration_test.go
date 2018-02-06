// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

func TestGetChartIndicator_Success(t *testing.T) {
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

	rsi, err := indicatorService.GetPlatformIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_Success2(t *testing.T) {
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
	rsi, err := indicatorService.GetPlatformIndicator("RelativeStrengthIndex")
	assert.Equal(t, nil, err)
	assert.Equal(t, "RelativeStrengthIndex", rsi.GetName())
	assert.Equal(t, "rsi.so", rsi.GetFilename())
	assert.Equal(t, "0.0.1a", rsi.GetVersion())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_IndicatorDoesntExistInDatabase(t *testing.T) {
	ctx := NewIntegrationTestContext()
	indicatorDAO := dao.NewIndicatorDAO(ctx)

	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	pluginService := NewPluginService(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	indicatorService := NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)

	doesntExistIndicator, err := indicatorService.GetPlatformIndicator("IndicatorDoesntExist")
	assert.NotNil(t, err)
	assert.Equal(t, nil, doesntExistIndicator)
	assert.Equal(t, "Failed to get platform indicator: IndicatorDoesntExist", err.Error())

	CleanupIntegrationTest()
}

func TestGetChartIndicator_IndicatorDoesntExist(t *testing.T) {
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

	doesntExistIndicator, err := indicatorService.GetPlatformIndicator("IndicatorDoesntExist")
	assert.Nil(t, err)
	assert.NotNil(t, doesntExistIndicator)
	assert.Equal(t, indicator.GetName(), doesntExistIndicator.GetName())
	assert.Equal(t, indicator.GetFilename(), doesntExistIndicator.GetFilename())
	assert.Equal(t, indicator.GetVersion(), doesntExistIndicator.GetVersion())

	CleanupIntegrationTest()
}
