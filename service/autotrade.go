package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/strategy"
)

type AutoTradeService interface {
	Trade()
}

type AutoTradeServiceImpl struct {
	ctx           *common.Context
	chartServices []common.ChartService
	chartDAO      dao.ChartDAO
	tradeService  common.TradeService
	profitService common.ProfitService
	AutoTradeService
}

func NewAutoTradeService(ctx *common.Context, exchangeService ExchangeService, chartDAO dao.ChartDAO,
	tradeService common.TradeService, profitService common.ProfitService) AutoTradeService {

	var chartServices []common.ChartService
	for _, chart := range chartDAO.Find(ctx.User) {
		ctx.Logger.Debugf("[NewAutoTradeService] Loading chart currency pair: %s-%s\n", chart.GetBase(), chart.GetQuote())
		currencyPair := &common.CurrencyPair{
			Base:          chart.GetBase(),
			Quote:         chart.GetQuote(),
			LocalCurrency: ctx.User.LocalCurrency}
		exchange := exchangeService.NewExchange(ctx.User, chart.GetExchangeName(), currencyPair)
		chartService := NewChartService(ctx, chartDAO, &chart, exchange)

		ctx.Logger.Debugf("[NewAutoTradeService] ChartService: %+v\n", chartService)
		chartServices = append(chartServices, chartService)
	}
	return &AutoTradeServiceImpl{
		ctx:           ctx,
		chartServices: chartServices,
		tradeService:  tradeService,
		profitService: profitService}
}

func (ats *AutoTradeServiceImpl) Trade() {
	for _, chartService := range ts.chartServices {
		strategy := strategy.NewDefaultTradingStrategy(ats.ctx, chartService, ats.tradeService, ats.profitService)
		chartService.Stream(ats.ChartTick(chartService))
	}
}

func (ats *AutoTradeServiceImpl) ChartTick(chartService common.ChartService) {
	strategy.OnPriceChange(chartService)
}
