package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
)

type IndicatorService interface {
	GetPlatformIndicator(name string) (dto.PlatformIndicator, error)
	GetChartIndicator(chart *common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error)
	GetChartIndicators(chart *common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error)
}

type IndicatorServiceImpl struct {
	ctx               *common.Context
	indicatorDAO      dao.IndicatorDAO
	chartIndicatorDAO dao.ChartIndicatorDAO
	pluginService     PluginService
	indicatorMapper   mapper.IndicatorMapper
	IndicatorService
}

func NewIndicatorService(ctx *common.Context, indicatorDAO dao.IndicatorDAO,
	chartIndicatorDAO dao.ChartIndicatorDAO, pluginService PluginService, indicatorMapper mapper.IndicatorMapper) IndicatorService {
	return &IndicatorServiceImpl{
		ctx:               ctx,
		indicatorDAO:      indicatorDAO,
		chartIndicatorDAO: chartIndicatorDAO,
		pluginService:     pluginService,
		indicatorMapper:   indicatorMapper}
}

func (service *IndicatorServiceImpl) GetPlatformIndicator(name string) (dto.PlatformIndicator, error) {
	entity, err := service.indicatorDAO.Get(name)
	if err != nil {
		return nil, err
	}
	return service.indicatorMapper.MapIndicatorEntityToDto(entity), nil
}

func (service *IndicatorServiceImpl) GetChartIndicator(chart *common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error) {
	daoChart := &dao.Chart{Id: chart.Id}
	chartIndicator, err := service.chartIndicatorDAO.Get(daoChart, name)
	if err != nil {
		return nil, err
	}
	indicator, err := service.GetPlatformIndicator(name)
	if err != nil {
		return nil, err
	}
	constructor, err := service.pluginService.GetIndicator(indicator.GetFilename(), name)
	if err != nil {
		return nil, err
	}
	params := strings.Split(chartIndicator.GetParameters(), ",")
	return constructor(candles, params)
}

func (service *IndicatorServiceImpl) GetChartIndicators(chart *common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error) {
	chartFinancialIndicators := make(map[string]common.FinancialIndicator, len(chart.Indicators))
	daoChart := &dao.Chart{Id: chart.Id}
	chartIndicators, err := service.chartIndicatorDAO.Find(daoChart)
	if err != nil {
		return nil, err
	}
	for _, ci := range chartIndicators {
		indicator, err := service.GetPlatformIndicator(ci.GetName())
		if err != nil {
			return nil, err
		}
		constructor, err := service.pluginService.GetIndicator(indicator.GetFilename(), ci.GetName())
		if err != nil {
			return nil, err
		}
		params := strings.Split(ci.GetParameters(), ",")
		FinancialIndicator, err := constructor(candles, params)
		if err != nil {
			return nil, err
		}
		chartFinancialIndicators[FinancialIndicator.GetName()] = FinancialIndicator
	}
	return chartFinancialIndicators, nil
}
