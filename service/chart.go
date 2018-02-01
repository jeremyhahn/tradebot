package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ChartServiceImpl struct {
	ctx          *common.Context
	dao          dao.ChartDAO
	entity       dao.IChart
	price        float64
	priceStream  *PriceStream
	priceChannel chan common.PriceChange
	closeChan    chan bool
	period       int
	candlesticks []common.Candlestick
	Exchange     common.Exchange
	Indicators   map[string]common.Indicator
	common.ChartService
}

func NewChartService(ctx *common.Context, chartDAO dao.ChartDAO, entity dao.IChart, exchange common.Exchange) common.ChartService {
	indicators := chartDAO.GetIndicators(entity)
	service := &ChartServiceImpl{
		ctx:         ctx,
		dao:         chartDAO,
		entity:      entity,
		Exchange:    exchange,
		Indicators:  make(map[string]common.Indicator, len(indicators)),
		priceStream: NewPriceStream(entity.GetPeriod()),
		closeChan:   make(chan bool, 1)}
	service.load()
	for _, indicator := range indicators {
		service.Indicators[indicator.Name] = service.createIndicator(&indicator)
	}
	return service
}

func (service *ChartServiceImpl) Find() []common.Chart {
	var charts []common.Chart
	_charts := service.dao.Find(service.ctx.User)
	for _, chart := range _charts {
		var indicators []common.Indicator
		var trades []common.Trade
		for _, i := range chart.GetIndicators() {
			indicators = append(indicators, common.Indicator{
				Id:         i.Id,
				ChartID:    i.ChartID,
				Name:       i.Name,
				Parameters: i.Parameters})
		}
		for _, trade := range chart.GetTrades() {
			trades = append(trades, common.Trade{
				ID:        trade.ID,
				UserID:    trade.UserID,
				ChartID:   chart.ID,
				Date:      trade.Date,
				Exchange:  trade.Exchange,
				Type:      trade.Type,
				Base:      trade.Base,
				Quote:     trade.Quote,
				Amount:    trade.Amount,
				Price:     trade.Price,
				ChartData: trade.ChartData})
		}
		charts = append(charts, common.Chart{
			ID:         chart.ID,
			Base:       chart.Base,
			Quote:      chart.Quote,
			Exchange:   chart.Exchange,
			Period:     chart.Period,
			Indicators: indicators,
			Trades:     trades})
	}
	return charts
}

func (service *ChartServiceImpl) GetExchange() common.Exchange {
	return service.Exchange
}

func (service *ChartServiceImpl) GetCurrencyPair() common.CurrencyPair {
	return service.Exchange.GetCurrencyPair()
}

func (service *ChartServiceImpl) GetPrice() float64 {
	return service.price
}

func (service *ChartServiceImpl) GetChart() *common.Chart {
	var trades []common.Trade
	var indicators []common.Indicator
	for _, entity := range service.entity.GetTrades() {
		trades = append(trades, common.Trade{
			ID:        entity.ID,
			UserID:    service.ctx.User.Id,
			ChartID:   entity.ChartID,
			Date:      entity.Date,
			Exchange:  entity.Exchange,
			Type:      entity.Type,
			Base:      entity.Base,
			Quote:     entity.Quote,
			Amount:    entity.Amount,
			Price:     entity.Price,
			ChartData: entity.ChartData})
	}
	for _, i := range service.entity.GetIndicators() {
		indicators = append(indicators, common.Indicator{
			Id:         i.Id,
			ChartID:    i.ChartID,
			Name:       i.Name,
			Parameters: i.Parameters})
	}
	return &common.Chart{
		ID:         service.entity.GetId(),
		Exchange:   service.entity.GetExchangeName(),
		Base:       service.entity.GetBase(),
		Quote:      service.entity.GetQuote(),
		Period:     service.entity.GetPeriod(),
		Indicators: indicators,
		Trades:     trades}
}

func (service *ChartServiceImpl) GetIndicators() []common.FinancialIndicator {
	var indicators []common.FinancialIndicator
	daoChart := &dao.Chart{ID: service.entity.GetId()}
	for _, indicator := range service.dao.GetIndicators(daoChart) {
		indicators = append(indicators, service.createIndicator(&indicator))
	}
	return indicators
}

func (service *ChartServiceImpl) GetIndicator(name string) common.FinancialIndicator {
	var indicator common.FinancialIndicator
	if i, ok := service.Indicators[name]; ok {
		indicator = i
	}
	return indicator
}

func (service *ChartServiceImpl) Stream(strategyHandler func(chart common.ChartService)) {
	service.ctx.Logger.Infof("[ChartService.Stream] Streaming %s %s chart data.",
		service.Exchange.GetName(), service.Exchange.FormattedCurrencyPair())
	service.load()
	for _, indicator := range service.Indicators {
		service.priceStream.SubscribeToPeriod(indicator)
	}
	priceChange := make(chan common.PriceChange)
	go service.Exchange.SubscribeToLiveFeed(priceChange)
	for {
		select {
		case <-service.closeChan:
			service.ctx.Logger.Debug("[ChartServiceImpl.Stream] Closing stream")
			return
		default:
			priceChange := service.priceStream.Listen(priceChange)
			service.price = priceChange.Price
			strategyHandler(service)
		}
	}
}

func (service *ChartServiceImpl) StopStream() {
	service.ctx.Logger.Debugf("[ChartServiceImpl.StopStream]")
	service.closeChan <- true
}

func (service *ChartServiceImpl) ToJSON() string {
	service.ctx.Logger.Debugf("[ChartServiceImpl.ToJSON]")
	data, _ := json.Marshal(service.GetChart())
	return string(data)
}

func (service *ChartServiceImpl) load() {

	t := time.Now()
	year, month, day := t.Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	now := time.Now()

	service.ctx.Logger.Debugf("[ChartService.Stream] Getting %s %s trade history from %s - %s ",
		service.Exchange.GetName(), service.Exchange.FormattedCurrencyPair(), yesterday, now)

	candlesticks := service.Exchange.GetPriceHistory(yesterday, now, service.period)
	if len(candlesticks) < 20 {
		service.ctx.Logger.Errorf("[ChartService.Stream] Failed to load initial candlesticks from %s. Total records: %d",
			service.Exchange.GetName(), len(candlesticks))
		return
	}

	var reversed []common.Candlestick
	for i := len(candlesticks) - 1; i > 0; i-- {
		reversed = append(reversed, candlesticks[i])
	}

	service.candlesticks = reversed
}

func (service *ChartServiceImpl) createIndicator(dao *dao.Indicator) common.Indicator {
	fqcn := fmt.Sprintf("indicators.%s", dao.Name)
	service.ctx.Logger.Debugf("[ChartServiceImpl.createIndicator] Creating indicator: %s", fqcn)
	return reflect.New(reflect.TypeOf(fqcn)).Elem().Interface().(common.Indicator)
}
