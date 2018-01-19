package strategy

import (
	"encoding/json"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/service"
)

type DefaultTradingStrategy struct {
	ctx           *common.Context
	Name          string
	autoTradeCoin *dao.AutoTradeCoin
	lastTrade     *dao.Trade
	service.TradingStrategy
}

func NewDefaultTradingStrategy(ctx *common.Context, autoTradeCoin *dao.AutoTradeCoin) *DefaultTradingStrategy {
	return &DefaultTradingStrategy{
		ctx:           ctx,
		Name:          "DefaultTradingStrategy",
		autoTradeCoin: autoTradeCoin}
}

func (strategy *DefaultTradingStrategy) OnPriceChange(chart *service.ChartService) {

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPriceChange] ChartData: %+v\n", chart.Data)

	buySignals := 0
	sellSignals := 0

	if chart.RSI.IsOverSold(chart.Data.RSILive) {
		buySignals++
		strategy.ctx.Logger.Debug("[DefaultTradingStrategy.OnPeriodChange] RSI buy signal!")
	} else if chart.RSI.IsOverBought(chart.Data.RSILive) {
		sellSignals++
		strategy.ctx.Logger.Debug("[DefaultTradingStrategy.OnPeriodChange] RSI sell signal!")
	}
	if chart.Data.Price >= chart.Data.BollingerUpperLive {
		sellSignals++
		strategy.ctx.Logger.Debug("[DefaultTradingStrategy.OnPeriodChange] Bollinger sell signal!")
	} else if chart.Data.Price <= chart.Data.BollingerLowerLive {
		buySignals++
		strategy.ctx.Logger.Debug("[DefaultTradingStrategy.OnPeriodChange] Bollinger buy signal!")
	}

	currencyPair := chart.Exchange.GetCurrencyPair()
	baseAmount := 0.0
	quoteAmount := 0.0
	coins, _ := chart.Exchange.GetBalances()
	for _, coin := range coins {
		if coin.Currency == currencyPair.Base {
			baseAmount = coin.Available * 0.10 // Only trade with 10% of current balance
		}
		if coin.Currency == currencyPair.Quote {
			quoteAmount = coin.Available * 0.10
		}
		if baseAmount > 0 && quoteAmount > 0 {
			break
		}
	}
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] Trading funds - baseAmount: %f, quoteAmount: %f",
		baseAmount, quoteAmount)

	autoTradeDAO := dao.NewAutoTradeDAO(strategy.ctx)
	strategy.lastTrade = autoTradeDAO.GetLastTrade(&dao.AutoTradeCoin{})

	jsonChart, err := json.Marshal(chart.Data)
	if err != nil {
		strategy.ctx.Logger.Errorf("[DefaultTradingStrategy.OnPeriodChange] Error marshalling chart state: %s", err.Error())
	}

	if buySignals == 2 {

		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] $$$ BUY SIGNAL $$$")
		if quoteAmount <= 0 {
			strategy.ctx.Logger.Errorf("[DefaultTradingStrategy.OnPeriodChange] Aborting. Out of %s funding!", chart.Data.CurrencyPair.Quote)
			return
		}
		if strategy.lastTrade.Type == "buy" {
			strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] Aborting. Already in a buy position.")
			return
		}
		minTradePrice := strategy.minBuyPrice()
		if chart.Data.Price <= minTradePrice {
			strategy.ctx.Logger.Debugf(
				"[DefaultTradingStrategy.OnPeriodChange] Aborting. Price does not meet minimum trade requirements. Price=%d, MinRequirement=%d",
				chart.Data.Price, minTradePrice)
			return
		}
		strategy.autoTradeCoin.Trades = append(strategy.autoTradeCoin.Trades, dao.Trade{
			UserID:    strategy.ctx.User.Id,
			Base:      chart.Data.CurrencyPair.Base,
			Quote:     chart.Data.CurrencyPair.Quote,
			Date:      time.Now(),
			Type:      "buy",
			Price:     chart.Data.Price,
			Amount:    baseAmount,
			ChartData: string(jsonChart)})
		autoTradeDAO.Save(strategy.autoTradeCoin)
		return
	}

	if sellSignals == 2 {

		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] $$$ SELL SIGNAL $$$")

		if strategy.lastTrade.Type == "sell" {
			strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] Aborting. Already in a sell position.")
			return
		}
		minTradePrice := strategy.minSellPrice(chart.Data.Price)
		if chart.Data.Price <= minTradePrice {
			strategy.ctx.Logger.Debugf(
				"[DefaultTradingStrategy.OnPeriodChange] Aborting. Price does not meet minimum trade requirements. Price=%d, MinRequirement=%d",
				chart.Data.Price, minTradePrice)
			return
		}
		strategy.autoTradeCoin.Trades = append(strategy.autoTradeCoin.Trades, dao.Trade{
			UserID:    strategy.ctx.User.Id,
			Base:      chart.Data.CurrencyPair.Base,
			Quote:     chart.Data.CurrencyPair.Quote,
			Date:      time.Now(),
			Type:      "sell",
			Price:     chart.Data.Price,
			Amount:    quoteAmount,
			ChartData: string(jsonChart)})
		autoTradeDAO.Save(strategy.autoTradeCoin)
		return
	}

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] buySignals=%d, sellSignals=%d", buySignals, sellSignals)
}

func (strategy *DefaultTradingStrategy) minBuyPrice() float64 {
	tax := strategy.lastTrade.Price * .40
	profitMargin := strategy.lastTrade.Price * .20
	tradeFees := strategy.lastTrade.Price * .10
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minBuyPrice] tax: %f, profitMargin: %f, tradeFees: %f",
		tax, profitMargin, tradeFees)
	return strategy.lastTrade.Price + tax + profitMargin + tradeFees
}

func (strategy *DefaultTradingStrategy) minSellPrice(currentPrice float64) float64 {
	tax := strategy.lastTrade.Price * .40
	profitMargin := strategy.lastTrade.Price * .20
	tradeFees := strategy.lastTrade.Price * .10
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minSellPrice] tax: %f, profitMargin: %f, tradeFees: %f",
		tax, profitMargin, tradeFees)
	return currentPrice + tax + profitMargin + tradeFees
}
