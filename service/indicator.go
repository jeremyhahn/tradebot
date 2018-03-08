package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
)

type IndicatorService interface {
	GetIndicator(name string) (common.Plugin, error)
	GetChartIndicator(chart common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error)
	GetChartIndicators(chart common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error)
}

type DefaultIndicatorService struct {
	ctx               common.Context
	pluginDAO         dao.PluginDAO
	chartIndicatorDAO dao.ChartIndicatorDAO
	pluginService     PluginService
	IndicatorService
}

func NewIndicatorService(ctx common.Context, chartIndicatorDAO dao.ChartIndicatorDAO, pluginService PluginService) IndicatorService {
	return &DefaultIndicatorService{
		ctx:               ctx,
		chartIndicatorDAO: chartIndicatorDAO,
		pluginService:     pluginService}
}

func (service *DefaultIndicatorService) GetIndicator(name string) (common.Plugin, error) {
	entity, err := service.pluginService.GetPlugin(name, common.INDICATOR_PLUGIN_TYPE)
	if err != nil {
		return nil, err
	}
	return service.pluginService.GetMapper().MapPluginEntityToDto(entity), nil
}

func (service *DefaultIndicatorService) GetChartIndicator(chart common.Chart, name string,
	candles []common.Candlestick) (common.FinancialIndicator, error) {
	daoChart := &entity.Chart{Id: chart.GetId()}
	chartIndicator, err := service.chartIndicatorDAO.Get(daoChart, name)
	if err != nil {
		return nil, err
	}
	constructor, err := service.pluginService.CreateIndicator(name)
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
		constructor, err := service.pluginService.CreateIndicator(ci.GetName())
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
