package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultChartService struct {
	ctx              *common.Context
	chartDAO         dao.ChartDAO
	charts           map[uint]common.Chart
	priceStreams     map[uint]PriceStream
	closeChans       map[uint]chan bool
	exchangeService  ExchangeService
	indicatorService IndicatorService
	ChartService
}

func NewChartService(ctx *common.Context, chartDAO dao.ChartDAO, exchangeService ExchangeService,
	indicatorService IndicatorService) ChartService {
	service := &DefaultChartService{
		ctx:              ctx,
		chartDAO:         chartDAO,
		charts:           make(map[uint]common.Chart),
		priceStreams:     make(map[uint]PriceStream),
		closeChans:       make(map[uint]chan bool),
		exchangeService:  exchangeService,
		indicatorService: indicatorService}
	return service
}

func (service *DefaultChartService) GetCurrencyPair(chart common.Chart) *common.CurrencyPair {
	return &common.CurrencyPair{
		Base:          chart.GetBase(),
		Quote:         chart.GetQuote(),
		LocalCurrency: service.ctx.User.GetLocalCurrency()}
}

func (service *DefaultChartService) GetExchange(chart common.Chart) common.Exchange {
	return service.exchangeService.GetExchange(service.ctx.User, chart.GetExchange())
}

func (service *DefaultChartService) Stream(chart common.Chart,
	candlesticks []common.Candlestick, strategyHandler func(price float64) error) error {

	chartId := chart.GetId()

	if _, ok := service.charts[chartId]; ok {
		return errors.New(fmt.Sprintf("Already streaming chart %s-%s", chart.GetBase(), chart.GetQuote()))
	}

	service.charts[chartId] = chart
	service.closeChans[chartId] = make(chan bool)

	currencyPair := service.GetCurrencyPair(chart)
	exchange := service.GetExchange(chart)

	service.ctx.Logger.Infof("[DefaultChartService.Stream] Streaming %s %s chart data.",
		exchange.GetName(), exchange.FormattedCurrencyPair(currencyPair))

	indicators, err := service.GetIndicators(chart, candlesticks)
	if err != nil {
		return err
	}

	service.priceStreams[chartId] = NewPriceStream(chart.GetPeriod())
	for _, indicator := range indicators {
		service.priceStreams[chartId].SubscribeToPeriod(indicator)
	}

	priceChange := make(chan common.PriceChange)
	go exchange.SubscribeToLiveFeed(currencyPair, priceChange)

	for {
		select {
		case <-service.closeChans[chartId]:
			service.ctx.Logger.Debug("[DefaultChartService.Stream] Closing stream")
			delete(service.charts, chartId)
			delete(service.closeChans, chartId)
			return nil
		default:
			priceChange := service.priceStreams[chartId].Listen(priceChange)
			strategyErr := strategyHandler(priceChange.Price)
			if strategyErr != nil {
				return strategyErr
			}
		}
	}
}

func (service *DefaultChartService) StopStream(chart common.Chart) {
	service.ctx.Logger.Debugf("[DefaultChartService.StopStream]")
	service.closeChans[chart.GetId()] <- true
}

func (service *DefaultChartService) GetCharts(autoTradeOnly bool) ([]common.Chart, error) {
	var charts []common.Chart
	_charts, err := service.chartDAO.Find(service.ctx.User, autoTradeOnly)
	if err != nil {
		return nil, err
	}
	mapper := mapper.NewChartMapper(service.ctx)
	for _, chart := range _charts {
		var indicators []common.ChartIndicator
		var trades []common.Trade
		for _, indicator := range chart.GetIndicators() {
			indicators = append(indicators, mapper.MapIndicatorEntityToDto(indicator))
		}
		for _, trade := range chart.GetTrades() {
			trades = append(trades, mapper.MapTradeEntityToDto(&trade))
		}
		charts = append(charts, &dto.ChartDTO{
			Id:         chart.Id,
			Base:       chart.Base,
			Quote:      chart.Quote,
			Exchange:   chart.Exchange,
			Period:     chart.Period,
			AutoTrade:  chart.AutoTrade,
			Indicators: indicators,
			Trades:     trades})
	}
	return charts, nil
}

func (service *DefaultChartService) GetTrades(chart common.Chart) ([]common.Trade, error) {
	var trades []common.Trade
	entities, err := service.chartDAO.GetTrades(service.ctx.User)
	if err != nil {
		return nil, err
	}
	mapper := mapper.NewChartMapper(service.ctx)
	for _, entity := range entities {
		trades = append(trades, mapper.MapTradeEntityToDto(&entity))
	}
	return trades, nil
}

func (service *DefaultChartService) GetLastTrade(chart common.Chart) (common.Trade, error) {
	daoChart := &entity.Chart{Id: chart.GetId()}
	entity, err := service.chartDAO.GetLastTrade(daoChart)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return &dto.TradeDTO{}, nil
	}
	mapper := mapper.NewChartMapper(service.ctx)
	tradeDTO := mapper.MapTradeEntityToDto(entity)
	return tradeDTO, nil
}

func (service *DefaultChartService) GetChart(id uint) (common.Chart, error) {
	var trades []common.Trade
	var indicators []common.ChartIndicator
	entity, err := service.chartDAO.Get(id)
	if err != nil {
		return nil, err
	}
	mapper := mapper.NewChartMapper(service.ctx)
	for _, entity := range entity.GetTrades() {
		trades = append(trades, mapper.MapTradeEntityToDto(&entity))
	}
	for _, indicator := range entity.GetIndicators() {
		indicators = append(indicators, mapper.MapIndicatorEntityToDto(indicator))
	}
	return &dto.ChartDTO{
		Id:         entity.GetId(),
		Exchange:   entity.GetExchangeName(),
		Base:       entity.GetBase(),
		Quote:      entity.GetQuote(),
		Period:     entity.GetPeriod(),
		Indicators: indicators,
		Trades:     trades}, nil
}

func (service *DefaultChartService) GetIndicator(chart common.Chart, name string, candles []common.Candlestick) (common.FinancialIndicator, error) {
	indicators, err := service.GetIndicators(chart, candles)
	if err != nil {
		return nil, err
	}
	for _, indicator := range indicators {
		if indicator.GetName() == name {
			return indicator, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to locate indicator: %s", name))
}

func (service *DefaultChartService) GetIndicators(chart common.Chart, candles []common.Candlestick) (map[string]common.FinancialIndicator, error) {
	indicators := make(map[string]common.FinancialIndicator, len(chart.GetIndicators()))
	entity := &entity.Chart{Id: chart.GetId()}
	daoIndicators, err := service.chartDAO.GetIndicators(entity)
	if err != nil {
		return nil, err
	}
	for _, daoIndicator := range daoIndicators {
		indicator, err := service.indicatorService.GetChartIndicator(chart, daoIndicator.GetName(), candles)
		if err != nil {
			return nil, err
		}
		if indicator == nil {
			return nil, errors.New(fmt.Sprintf("Unable to create indicator instance: %s", indicator))
		}
		indicators[daoIndicator.Name] = indicator
	}
	return indicators, nil
}

func (service *DefaultChartService) LoadCandlesticks(chart common.Chart, exchange common.Exchange) []common.Candlestick {
	var candles []common.Candlestick
	t := time.Now()
	year, month, day := t.Date()
	//yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	lastWeek := time.Date(year, month, day-7, 0, 0, 0, 0, t.Location())
	now := time.Now()
	currencyPair := &common.CurrencyPair{
		Base:          chart.GetBase(),
		Quote:         chart.GetQuote(),
		LocalCurrency: service.ctx.User.GetLocalCurrency()}
	service.ctx.Logger.Debugf("[DefaultChartService.LoadCandlesticks] Getting %s %s trade history from %s - %s ",
		exchange.GetName(), exchange.FormattedCurrencyPair(currencyPair), lastWeek, now)
	candles = exchange.GetPriceHistory(currencyPair, lastWeek, now, chart.GetPeriod())
	if service.ctx.Debug {
		for _, c := range candles {
			service.ctx.Logger.Debugf("Prewarming kline: %s", c.ToString())
		}
	}
	if len(candles) < 35 {
		service.ctx.Logger.Errorf("[DefaultChartService.LoadCandlesticks] Failed to load initial price history from %s. Total candlesticks: %d",
			exchange.GetName(), len(candles))
		return candles
	}
	return candles
}
