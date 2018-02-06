package service

import (
	"github.com/jeremyhahn/tradebot/common"
)

type AutoTradeServiceImpl struct {
	ctx             *common.Context
	exchangeService ExchangeService
	chartService    ChartService
	tradeService    TradeService
	profitService   ProfitService
	pluginService   PluginService
	AutoTradeService
}

func NewAutoTradeService(ctx *common.Context, exchangeService ExchangeService, chartService ChartService,
	profitService ProfitService, tradeService TradeService, pluginService PluginService) AutoTradeService {
	return &AutoTradeServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService,
		chartService:    chartService,
		tradeService:    tradeService,
		profitService:   profitService,
		pluginService:   pluginService}
}

func (ats *AutoTradeServiceImpl) EndWorldHunger() error {
	charts, err := ats.chartService.GetCharts()
	if err != nil {
		return err
	}
	for _, chart := range charts {
		ats.ctx.Logger.Debugf("[AutoTradeService.EndWorldHunger] Loading chart %s-%s\n", chart.Base, chart.Quote)
		/*
			indicators, err2 := ats.chartService.GetIndicators(&chart)
			if err2 != nil {
				return err2
			}
			exchange := ats.exchangeService.CreateExchange(ats.ctx.User, chart.Exchange)
			coins, _ := exchange.GetBalances()
			lastTrade, err3 := ats.chartService.GetLastTrade(chart)
			if err3 != nil {
				return err3
			}
			//go func() {
			streamErr := ats.chartService.Stream(&chart, func(newPrice float64) error {
				params := &common.TradingStrategyParams{
					CurrencyPair: &common.CurrencyPair{
						Base:          chart.Base,
						Quote:         chart.Quote,
						LocalCurrency: ats.ctx.User.LocalCurrency},
					Balances:   coins,
					NewPrice:   newPrice,
					LastTrade:  lastTrade,
					Indicators: indicators}

				strategy, err4 := strategies.CreateDefaultTradingStrategy(params)
				if err4 != nil {
					return err4
				}
				buy, sell, err5 := strategy.GetBuySellSignals()
				if err5 != nil {
					return err5
				}
				if buy {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ BUY SIGNAL $$$")
				} else if sell {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ SELL SIGNAL $$$")
				}
				return nil
			})
			if streamErr != nil {
				return streamErr
			}
			return nil
			//}()
		*/
	}
	return nil
}
