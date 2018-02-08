package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type DefaultAutoTradeService struct {
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
	return &DefaultAutoTradeService{
		ctx:             ctx,
		exchangeService: exchangeService,
		chartService:    chartService,
		tradeService:    tradeService,
		profitService:   profitService,
		strategyService: strategyService}
}

func (ats *DefaultAutoTradeService) EndWorldHunger() error {
	charts, err := ats.chartService.GetCharts()
	if err != nil {
		return err
	}
	for _, autoTradeChart := range charts {
		ats.ctx.Logger.Debugf("[AutoTradeService.EndWorldHunger] Loading chart %s-%s\n",
			autoTradeChart.GetBase(), autoTradeChart.GetQuote())

		exchange := ats.exchangeService.CreateExchange(ats.ctx.User, autoTradeChart.GetExchange())
		candlesticks := ats.chartService.LoadCandlesticks(autoTradeChart, exchange)

		currencyPair := &common.CurrencyPair{
			Base:          autoTradeChart.GetBase(),
			Quote:         autoTradeChart.GetQuote(),
			LocalCurrency: ats.ctx.User.LocalCurrency}

		indicators, err := ats.chartService.GetIndicators(autoTradeChart, candlesticks)
		if err != nil {
			return err
		}

		coins, _ := exchange.GetBalances()
		lastTrade, err := ats.chartService.GetLastTrade(autoTradeChart)
		if err != nil {
			return err
		}

		go func(chart common.Chart) {

			streamErr := ats.chartService.Stream(chart, candlesticks, func(currentPrice float64) error {

				params := common.TradingStrategyParams{
					CurrencyPair: currencyPair,
					Balances:     coins,
					NewPrice:     currentPrice,
					LastTrade:    lastTrade,
					Indicators:   indicators}

				strategies, err := ats.strategyService.GetChartStrategies(chart, &params, candlesticks)
				if err != nil {
					return err
				}

				for _, strategy := range strategies {

					buy, sell, data, err := strategy.Analyze()
					ats.ctx.Logger.Debugf("[DefaultAutoTradeService.EndWorldHunger] Indicator data: %+v\n", data)
					if err != nil {
						return err
					}

					if buy || sell {
						var tradeType string
						if buy {
							ats.ctx.Logger.Debug("[DefaultAutoTradeService.EndWorldHunger] $$$ BUY SIGNAL $$$")
							tradeType = "buy"
						} else if sell {
							ats.ctx.Logger.Debug("[DefaultAutoTradeService.EndWorldHunger] $$$ SELL SIGNAL $$$")
							tradeType = "sell"
						}
						_, quoteAmount := strategy.GetTradeAmounts()
						fee, tax := strategy.CalculateFeeAndTax(currentPrice)
						chartJSON, err := chart.ToJSON()
						if err != nil {
							return err
						}
						thisTrade := &dto.TradeDTO{
							UserId:    ats.ctx.User.Id,
							Exchange:  exchange.GetName(),
							Base:      chart.GetBase(),
							Quote:     chart.GetQuote(),
							Date:      time.Now(),
							Type:      tradeType,
							Price:     currentPrice,
							Amount:    quoteAmount,
							ChartData: chartJSON}
						thisProfit := &dao.Profit{
							UserId:   ats.ctx.User.Id,
							TradeId:  thisTrade.GetId(),
							Quantity: quoteAmount,
							Bought:   lastTrade.GetPrice(),
							Sold:     currentPrice,
							Fee:      fee,
							Tax:      tax,
							Total:    currentPrice - lastTrade.GetPrice() - fee - tax}
						ats.tradeService.Save(thisTrade)
						ats.profitService.Save(thisProfit)
					}
				}
				return nil
			})
			if streamErr != nil {
				ats.ctx.Logger.Error(streamErr.Error())
			}
		}(autoTradeChart)

	}
	return nil
}
