package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
)

type StrategyService interface {
	GetStrategy(name string) (common.Strategy, error)
	GetChartStrategy(chart common.Chart, name string, candles []common.Candlestick) (common.TradingStrategy, error)
	GetChartStrategies(chart common.Chart, params *common.TradingStrategyParams, candles []common.Candlestick) ([]common.TradingStrategy, error)
}

type DefaultStrategyService struct {
	ctx              *common.Context
	strategyDAO      dao.StrategyDAO
	chartStrategyDAO dao.ChartStrategyDAO
	pluginService    PluginService
	indicatorService IndicatorService
	chartMapper      mapper.ChartMapper
	strategyMapper   mapper.StrategyMapper
	StrategyService
}

func NewStrategyService(ctx *common.Context, strategyDAO dao.StrategyDAO,
	chartStrategyDAO dao.ChartStrategyDAO, pluginService PluginService,
	indicatorService IndicatorService, chartMapper mapper.ChartMapper,
	strategyMapper mapper.StrategyMapper) StrategyService {
	return &DefaultStrategyService{
		ctx:              ctx,
		strategyDAO:      strategyDAO,
		chartStrategyDAO: chartStrategyDAO,
		pluginService:    pluginService,
		indicatorService: indicatorService,
		chartMapper:      chartMapper,
		strategyMapper:   strategyMapper}
}

func (service *DefaultStrategyService) GetStrategy(name string) (common.Strategy, error) {
	entity, err := service.strategyDAO.Get(name)
	if err != nil {
		return nil, err
	}
	return service.strategyMapper.MapStrategyEntityToDto(entity), nil
}

func (service *DefaultStrategyService) GetChartStrategy(chart common.Chart, name string, candles []common.Candlestick) (common.TradingStrategy, error) {
	financialIndicators, err := service.indicatorService.GetChartIndicators(chart, candles)
	if err != nil {
		return nil, err
	}
	strategy, err := service.GetStrategy(name)
	if err != nil {
		return nil, err
	}
	constructor, err := service.pluginService.GetStrategy(strategy.GetFilename(), name)
	if err != nil {
		return nil, err
	}
	trades := chart.GetTrades()
	tradeLen := len(trades)
	lastTrade := trades[tradeLen-1]
	params := common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{
			Base:          chart.GetBase(),
			Quote:         chart.GetQuote(),
			LocalCurrency: service.ctx.User.GetLocalCurrency()},
		LastTrade:  lastTrade,
		Indicators: financialIndicators}
	return constructor(&params)
}

func (service *DefaultStrategyService) GetChartStrategies(chart common.Chart, params *common.TradingStrategyParams,
	candles []common.Candlestick) ([]common.TradingStrategy, error) {
	var strategies []common.TradingStrategy
	daoChart := service.chartMapper.MapChartDtoToEntity(chart)
	strategyEntities, err := service.chartStrategyDAO.Find(daoChart)
	if err != nil {
		return nil, err
	}
	for _, strategyEntity := range strategyEntities {
		platformStrategy, err := service.GetStrategy(strategyEntity.GetName())
		if err != nil {
			return nil, err
		}
		constructor, err := service.pluginService.GetStrategy(platformStrategy.GetFilename(), strategyEntity.GetName())
		if err != nil {
			return nil, err
		}
		TradingStrategy, err := constructor(params)
		if err != nil {
			return nil, err
		}
		strategies = append(strategies, TradingStrategy)
	}
	return strategies, nil
}
