package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/shopspring/decimal"
)

type DefaultAutoTradeService struct {
	ctx             common.Context
	exchangeService ExchangeService
	chartService    ChartService
	tradeService    TradeService
	profitService   ProfitService
	strategyService StrategyService
	userMapper      mapper.UserMapper
	AutoTradeService
}

func NewAutoTradeService(ctx common.Context, exchangeService ExchangeService, chartService ChartService,
	profitService ProfitService, tradeService TradeService, strategyService StrategyService,
	userMapper mapper.UserMapper) AutoTradeService {
	return &DefaultAutoTradeService{
		ctx:             ctx,
		exchangeService: exchangeService,
		chartService:    chartService,
		tradeService:    tradeService,
		profitService:   profitService,
		strategyService: strategyService,
		userMapper:      userMapper}
}

func (ats *DefaultAutoTradeService) EndWorldHunger() error {
	if ats.ctx.GetUser() == nil {
		ats.ctx.GetLogger().Warningf("[DefaultAutoTradeService.EndWorldHunger] No users configured")
		return nil
	}
	charts, err := ats.chartService.GetCharts(true)
	if err != nil {
		return err
	}
	for _, autoTradeChart := range charts {
		ats.ctx.GetLogger().Debugf("[AutoTradeService.EndWorldHunger] Loading chart %s-%s\n",
			autoTradeChart.GetBase(), autoTradeChart.GetQuote())

		//userEntity := ats.userMapper.MapUserDtoToEntity(ats.ctx.GetUser())
		exchange, err := ats.exchangeService.CreateExchange(autoTradeChart.GetExchange())
		if err != nil {
			return err
		}

		candlesticks := ats.chartService.LoadCandlesticks(autoTradeChart, exchange)

		currencyPair := &common.CurrencyPair{
			Base:          autoTradeChart.GetBase(),
			Quote:         autoTradeChart.GetQuote(),
			LocalCurrency: ats.ctx.GetUser().GetLocalCurrency()}

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

			streamErr := ats.chartService.Stream(chart, candlesticks, func(currentPrice decimal.Decimal) error {

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
					ats.ctx.GetLogger().Debugf("[DefaultAutoTradeService.EndWorldHunger] Indicator data: %+v\n", data)
					if err != nil {
						return err
					}

					if buy || sell {
						var tradeType string
						if buy {
							ats.ctx.GetLogger().Debug("[DefaultAutoTradeService.EndWorldHunger] $$$ BUY SIGNAL $$$")
							tradeType = "buy"
						} else if sell {
							ats.ctx.GetLogger().Debug("[DefaultAutoTradeService.EndWorldHunger] $$$ SELL SIGNAL $$$")
							tradeType = "sell"
						}
						_, quoteAmount := strategy.GetTradeAmounts()
						fee, tax := strategy.CalculateFeeAndTax(currentPrice)
						chartJSON, err := chart.ToJSON()
						if err != nil {
							return err
						}
						thisTrade := &dto.TradeDTO{
							UserId:    ats.ctx.GetUser().GetId(),
							Exchange:  exchange.GetName(),
							Base:      chart.GetBase(),
							Quote:     chart.GetQuote(),
							Date:      time.Now(),
							Type:      tradeType,
							Price:     currentPrice,
							Amount:    quoteAmount,
							ChartData: chartJSON}
						thisProfit := &dto.ProfitDTO{
							UserId:   ats.ctx.GetUser().GetId(),
							TradeId:  thisTrade.GetId(),
							Quantity: quoteAmount,
							Bought:   lastTrade.GetPrice(),
							Sold:     currentPrice,
							Fee:      fee,
							Tax:      tax,
							Total:    currentPrice.Sub(lastTrade.GetPrice()).Sub(fee).Sub(tax)}
						ats.tradeService.Save(thisTrade)
						ats.profitService.Save(thisProfit)
					}
				}
				return nil
			})
			if streamErr != nil {
				ats.ctx.GetLogger().Error(streamErr.Error())
			}
		}(autoTradeChart)

	}
	return nil
}
