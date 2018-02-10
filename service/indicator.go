package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type IndicatorService interface {
	GetIndicator(name string) (common.Indicator, error)
	GetChartIndicator(chart common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error)
	GetChartIndicators(chart common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error)
}

type DefaultIndicatorService struct {
	ctx               *common.Context
	indicatorDAO      dao.IndicatorDAO
	chartIndicatorDAO dao.ChartIndicatorDAO
	pluginService     PluginService
	indicatorMapper   mapper.IndicatorMapper
	IndicatorService
}

func NewIndicatorService(ctx *common.Context, indicatorDAO dao.IndicatorDAO,
	chartIndicatorDAO dao.ChartIndicatorDAO, pluginService PluginService, indicatorMapper mapper.IndicatorMapper) IndicatorService {
	return &DefaultIndicatorService{
		ctx:               ctx,
		indicatorDAO:      indicatorDAO,
		chartIndicatorDAO: chartIndicatorDAO,
		pluginService:     pluginService,
		indicatorMapper:   indicatorMapper}
}

func (service *DefaultIndicatorService) GetIndicator(name string) (common.Indicator, error) {
	entity, err := service.indicatorDAO.Get(name)
	if err != nil {
		return nil, err
	}
	return service.indicatorMapper.MapIndicatorEntityToDto(entity), nil
}

func (service *DefaultIndicatorService) GetChartIndicator(chart common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error) {
	daoChart := &entity.Chart{Id: chart.GetId()}
	chartIndicator, err := service.chartIndicatorDAO.Get(daoChart, name)
	if err != nil {
		return nil, err
	}
	indicator, err := service.GetIndicator(name)
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

func (service *DefaultIndicatorService) GetChartIndicators(chart common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error) {
	chartFinancialIndicators := make(map[string]common.FinancialIndicator, len(chart.GetIndicators()))
	daoChart := &entity.Chart{Id: chart.GetId()}
	chartIndicators, err := service.chartIndicatorDAO.Find(daoChart)
	if err != nil {
		return nil, err
	}
	for _, ci := range chartIndicators {
		indicator, err := service.GetIndicator(ci.GetName())
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
