package service

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ChartServiceImpl struct {
	ctx             *common.Context
	chartDAO        dao.ChartDAO
	charts          map[uint]*common.Chart
	priceStreams    map[uint]PriceStream
	closeChans      map[uint]chan bool
	exchangeService ExchangeService
	ChartService
}

func NewChartService(ctx *common.Context, chartDAO dao.ChartDAO, exchangeService ExchangeService) ChartService {
	service := &ChartServiceImpl{
		ctx:             ctx,
		chartDAO:        chartDAO,
		charts:          make(map[uint]*common.Chart),
		priceStreams:    make(map[uint]PriceStream),
		closeChans:      make(map[uint]chan bool),
		exchangeService: exchangeService}
	return service
}

func (service *ChartServiceImpl) GetCurrencyPair(chart *common.Chart) *common.CurrencyPair {
	return &common.CurrencyPair{
		Base:          chart.Base,
		Quote:         chart.Quote,
		LocalCurrency: service.ctx.User.LocalCurrency}
}

func (service *ChartServiceImpl) GetExchange(chart *common.Chart) common.Exchange {
	return service.exchangeService.GetExchange(service.ctx.User, chart.Exchange)
}

func (service *ChartServiceImpl) Stream(chart *common.Chart, strategyHandler func(price float64) error) error {

	if _, ok := service.charts[chart.Id]; !ok {
		return errors.New(fmt.Sprintf("Already streaming chart %s-%s", chart.Base, chart.Quote))
	}

	service.charts[chart.Id] = chart
	service.closeChans[chart.Id] = make(chan bool)

	currencyPair := service.GetCurrencyPair(chart)
	exchange := service.GetExchange(chart)

	service.ctx.Logger.Infof("[ChartServiceImpl.Stream] Streaming %s %s chart data.",
		exchange.GetName(), exchange.FormattedCurrencyPair(currencyPair))

	indicators, err := service.GetIndicators(chart)
	if err != nil {
		return err
	}
	for _, indicator := range indicators {
		service.priceStreams[chart.Id] = NewPriceStream(chart.Period)
		service.priceStreams[chart.Id].SubscribeToPeriod(indicator)
	}

	priceChange := make(chan common.PriceChange)
	go exchange.SubscribeToLiveFeed(currencyPair, priceChange)
	for {
		select {
		case <-service.closeChans[chart.Id]:
			service.ctx.Logger.Debug("[ChartServiceImpl.Stream] Closing stream")
			delete(service.charts, chart.Id)
			delete(service.closeChans, chart.Id)
			return nil
		default:
			priceChange := service.priceStreams[chart.Id].Listen(priceChange)
			chart.Price = priceChange.Price
			strategyErr := strategyHandler(chart.Price)
			if strategyErr != nil {
				return strategyErr
			}
		}
	}
}

func (service *ChartServiceImpl) StopStream(chart *common.Chart) {
	service.ctx.Logger.Debugf("[ChartServiceImpl.StopStream]")
	service.closeChans[chart.Id] <- true
}

func (service *ChartServiceImpl) GetCharts() ([]common.Chart, error) {
	var charts []common.Chart
	_charts, err := service.chartDAO.Find(service.ctx.User)
	if err != nil {
		return nil, err
	}
	mapper := mapper.NewChartMapper(service.ctx)
	for _, chart := range _charts {
		var indicators []common.Indicator
		var trades []common.Trade
		for _, indicator := range chart.GetIndicators() {
			indicators = append(indicators, mapper.MapIndicatorEntityToDto(indicator))
		}
		for _, trade := range chart.GetTrades() {
			trades = append(trades, mapper.MapTradeEntityToDto(&trade))
		}
		charts = append(charts, common.Chart{
			Id:         chart.Id,
			Base:       chart.Base,
			Quote:      chart.Quote,
			Exchange:   chart.Exchange,
			Period:     chart.Period,
			Indicators: indicators,
			Trades:     trades})
	}
	return charts, nil
}

func (service *ChartServiceImpl) GetTrades(chart *common.Chart) ([]common.Trade, error) {
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

func (service *ChartServiceImpl) GetLastTrade(chart common.Chart) (*common.Trade, error) {
	daoChart := &dao.Chart{Id: chart.Id}
	entity, err := service.chartDAO.GetLastTrade(daoChart)
	if err != nil {
		return nil, err
	}
	mapper := mapper.NewChartMapper(service.ctx)
	tradeDTO := mapper.MapTradeEntityToDto(entity)
	return &tradeDTO, nil
}

func (service *ChartServiceImpl) GetChart(id uint) (*common.Chart, error) {
	var trades []common.Trade
	var indicators []common.Indicator
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
	return &common.Chart{
		Id:         entity.GetId(),
		Exchange:   entity.GetExchangeName(),
		Base:       entity.GetBase(),
		Quote:      entity.GetQuote(),
		Period:     entity.GetPeriod(),
		Indicators: indicators,
		Trades:     trades}, nil
}

func (service *ChartServiceImpl) GetIndicator(chart *common.Chart, name string) (common.FinancialIndicator, error) {
	indicators, err := service.GetIndicators(chart)
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

func (service *ChartServiceImpl) GetIndicators(chart *common.Chart) (map[string]common.FinancialIndicator, error) {
	var indicators map[string]common.FinancialIndicator
	entity := &dao.Chart{Id: chart.Id}
	daoIndicators, err := service.chartDAO.GetIndicators(entity)
	if err != nil {
		return nil, err
	}
	for _, daoIndicator := range daoIndicators {
		indicator := service.CreateIndicator(&daoIndicator)
		if indicator != nil {
			return nil, errors.New(fmt.Sprintf("Unable to create indicator instace: %s", indicator))
		}
		indicators[daoIndicator.Name] = indicator
	}
	return indicators, nil
}

func (service *ChartServiceImpl) loadCandlesticks(chart *common.Chart, exchange common.Exchange) []common.Candlestick {
	var candles []common.Candlestick
	t := time.Now()
	year, month, day := t.Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	now := time.Now()

	currencyPair := &common.CurrencyPair{
		Base:          chart.Base,
		Quote:         chart.Quote,
		LocalCurrency: service.ctx.User.LocalCurrency}

	service.ctx.Logger.Debugf("[ChartServiceImpl.loadCandlesticks] Getting %s %s trade history from %s - %s ",
		exchange.GetName(), exchange.FormattedCurrencyPair(currencyPair), yesterday, now)

	candles = exchange.GetPriceHistory(currencyPair, yesterday, now, chart.Period)
	if len(candles) < 20 {
		service.ctx.Logger.Errorf("[ChartServiceImpl.loadCandlesticks] Failed to load initial candlesticks from %s. Total records: %d",
			exchange.GetName(), len(candles))
		return candles
	}
	var reversed []common.Candlestick
	for i := len(candles) - 1; i > 0; i-- {
		reversed = append(reversed, candles[i])
	}
	return reversed
}

func (service *ChartServiceImpl) CreateIndicator(dao *dao.Indicator) common.FinancialIndicator {
	//var candles []common.Candlestick
	//rsi := indicators.NewRelativeStrengthIndex(candles)

	fqcn := fmt.Sprintf("*indicators.%s", dao.Name)
	service.ctx.Logger.Debugf("[ChartServiceImpl.createIndicator] Creating indicator: %s", fqcn)
	elem := reflect.New(reflect.TypeOf(fqcn)).Elem()
	return elem.Interface().(common.FinancialIndicator)
}
