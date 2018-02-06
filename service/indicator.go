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
	userIndicator, err := service.chartIndicatorDAO.Get(daoChart, name)
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
	params := strings.Split(userIndicator.GetParameters(), ",")
	return constructor(candles, params)
}
