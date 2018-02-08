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
	strategyService StrategyService
	AutoTradeService
}

func NewAutoTradeService(ctx *common.Context, exchangeService ExchangeService, chartService ChartService,
	profitService ProfitService, tradeService TradeService, strategyService StrategyService) AutoTradeService {
	return &AutoTradeServiceImpl{
		ctx:             ctx,
		exchangeService: exchangeService,
		chartService:    chartService,
		tradeService:    tradeService,
		profitService:   profitService,
		strategyService: strategyService}
}

func (ats *AutoTradeServiceImpl) EndWorldHunger() error {
	charts, err := ats.chartService.GetCharts()
	if err != nil {
		return err
	}
	for _, chart := range charts {
		ats.ctx.Logger.Debugf("[AutoTradeService.EndWorldHunger] Loading chart %s-%s\n", chart.GetBase(), chart.GetQuote())

		exchange := ats.exchangeService.CreateExchange(ats.ctx.User, chart.GetExchange())
		candlesticks := ats.chartService.LoadCandlesticks(chart, exchange)

		currencyPair := &common.CurrencyPair{
			Base:          chart.GetBase(),
			Quote:         chart.GetQuote(),
			LocalCurrency: ats.ctx.User.LocalCurrency}

		indicators, err := ats.chartService.GetIndicators(chart, candlesticks)
		if err != nil {
			return err
		}

		coins, _ := exchange.GetBalances()
		lastTrade, err := ats.chartService.GetLastTrade(chart)
		if err != nil {
			return err
		}

		//go func() {

		streamErr := ats.chartService.Stream(chart, candlesticks, func(newPrice float64) error {

			params := common.TradingStrategyParams{
				CurrencyPair: currencyPair,
				Balances:     coins,
				NewPrice:     newPrice,
				LastTrade:    lastTrade,
				Indicators:   indicators}

			strategies, err := ats.strategyService.GetChartStrategies(chart, &params, candlesticks)
			if err != nil {
				return err
			}

			for _, strategy := range strategies {

				buy, sell, data, err := strategy.Analyze()
				ats.ctx.Logger.Debugf("[AutoTradeServiceImpl.EndWorldHunger] Indicator data: %+v\n", data)
				if err != nil {
					return err
				}

				if buy {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ BUY SIGNAL $$$")
				} else if sell {
					ats.ctx.Logger.Debug("[AutoTradeServiceImpl.EndWorldHunger] $$$ SELL SIGNAL $$$")
				}
			}
			return nil
		})
		if streamErr != nil {
			ats.ctx.Logger.Error(streamErr.Error())
		}
		//}()

	}
	return nil
}
