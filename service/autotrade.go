package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/strategy"
)

type AutoTradeServiceImpl struct {
	ctx             *common.Context
	exchangeService ExchangeService
	chartService    ChartService
	tradeService    TradeService
	profitService   ProfitService
	AutoTradeService
}

func NewAutoTradeService(ctx *common.Context, exchangeService ExchangeService, chartService ChartService,
	profitService ProfitService, tradeService TradeService) AutoTradeService {
	return &AutoTradeServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService,
		chartService:    chartService,
		tradeService:    tradeService,
		profitService:   profitService}
}

func (ats *AutoTradeServiceImpl) EndWorldHunger() {

	for _, chart := range ats.chartService.GetCharts() {

		ats.ctx.Logger.Debugf("[AutoTradeService.EndWorldHunger] Loading chart %s-%s\n", chart.Base, chart.Quote)

		exchange := ats.exchangeService.CreateExchange(ats.ctx.User, chart.Exchange)

		coins, _ := exchange.GetBalances()
		params := &strategy.TradingStrategyParams{
			CurrencyPair: &common.CurrencyPair{
				Base:          chart.Base,
				Quote:         chart.Quote,
				LocalCurrency: ats.ctx.User.LocalCurrency},
			Balances:   coins,
			Indicators: ats.chartService.GetIndicators(&chart)}

		s, err := strategy.CreateDefaultTradingStrategy(params)
		if err != nil {
			ats.ctx.Logger.Errorf("[AutoTradeServiceImpl.EndWorldHunger] %s", err.Error())
			continue
		}

		go func() {
			ats.chartService.Stream(&chart, func(newPrice float64) {
				buy, sell, err := s.GetBuySellSignals()
				if err != nil {
					ats.ctx.Logger.Errorf("[AutoTradeServiceImpl.EndWorldHunger] %s", err.Error())
					return
				}
				if buy {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ BUY SIGNAL $$$")
				} else if sell {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ SELL SIGNAL $$$")
				}
			})
		}()
	}
}
