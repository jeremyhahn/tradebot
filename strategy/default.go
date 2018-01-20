package strategy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/service"
)

type DefaultTradingStrategy struct {
	ctx                 *common.Context
	autoTradeDAO        *dao.AutoTradeDAO
	signalLogDAO        *dao.SignalLogDAO
	autoTradeCoin       *dao.AutoTradeCoin
	lastTrade           *dao.Trade
	Name                string
	rsiOverSold         float64
	rsiOverBought       float64
	tax                 float64
	fees                float64
	profitMargin        float64
	stopLoss            float64
	requiredBuySignals  int
	requiredSellSignals int
	service.TradingStrategy
}

func NewDefaultTradingStrategy(ctx *common.Context, autoTradeCoin *dao.AutoTradeCoin,
	autoTradeDAO *dao.AutoTradeDAO, signalLogDAO *dao.SignalLogDAO) *DefaultTradingStrategy {
	return &DefaultTradingStrategy{
		Name:                "DefaultTradingStrategy",
		ctx:                 ctx,
		autoTradeCoin:       autoTradeCoin,
		autoTradeDAO:        autoTradeDAO,
		signalLogDAO:        signalLogDAO,
		rsiOverSold:         30,
		rsiOverBought:       70,
		tax:                 .40,
		fees:                .10,
		profitMargin:        .20,
		stopLoss:            .10,
		requiredBuySignals:  2,
		requiredSellSignals: 2}
}

func (strategy *DefaultTradingStrategy) OnPriceChange(chart common.ChartService) {

	data := chart.GetData()
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPriceChange] ChartData: %+v\n", data)

	buySignals, sellSignals := strategy.countSignals(data)
	strategy.lastTrade = strategy.autoTradeDAO.GetLastTrade(&dao.AutoTradeCoin{})

	if buySignals == strategy.requiredBuySignals {
		strategy.buy(chart)
		return
	}

	if sellSignals == strategy.requiredSellSignals {
		strategy.sell(chart)
		return
	}

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] buySignals=%d, sellSignals=%d", buySignals, sellSignals)
}

func (strategy *DefaultTradingStrategy) countSignals(data *common.ChartData) (int, int) {
	var buySignals int
	var sellSignals int
	if data.RSILive < strategy.rsiOverSold {
		buySignals++
		strategy.signalLogDAO.Save(&dao.SignalLog{
			UserID:     strategy.ctx.User.Id,
			Date:       time.Now(),
			Name:       "RSI",
			Type:       "buy",
			Price:      data.Price,
			SignalData: strconv.FormatFloat(data.RSILive, 'f', 8, 64)})
	} else if data.RSILive > strategy.rsiOverBought {
		sellSignals++
		strategy.signalLogDAO.Save(&dao.SignalLog{
			UserID:     strategy.ctx.User.Id,
			Date:       time.Now(),
			Name:       "RSI",
			Type:       "sell",
			Price:      data.Price,
			SignalData: strconv.FormatFloat(data.RSILive, 'f', 8, 64)})
	}
	if data.Price > data.BollingerUpperLive {
		sellSignals++
		strategy.signalLogDAO.Save(&dao.SignalLog{
			UserID:     strategy.ctx.User.Id,
			Date:       time.Now(),
			Name:       "Bollinger",
			Type:       "sell",
			Price:      data.Price,
			SignalData: fmt.Sprintf("%f,%f,%f", data.BollingerUpperLive, data.BollingerMiddleLive, data.BollingerLowerLive)})
	} else if data.Price < data.BollingerLowerLive {
		buySignals++
		strategy.signalLogDAO.Save(&dao.SignalLog{
			UserID:     strategy.ctx.User.Id,
			Date:       time.Now(),
			Name:       "Bollinger",
			Type:       "buy",
			Price:      data.Price,
			SignalData: fmt.Sprintf("%f,%f,%f", data.BollingerUpperLive, data.BollingerMiddleLive, data.BollingerLowerLive)})
	}
	return buySignals, sellSignals
}

func (strategy *DefaultTradingStrategy) minBuyPrice() float64 {
	tax := strategy.lastTrade.Price * strategy.tax
	profitMargin := strategy.lastTrade.Price * strategy.profitMargin
	tradeFees := strategy.lastTrade.Price * strategy.fees
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minBuyPrice] tax: %f, profitMargin: %f, tradeFees: %f",
		tax, profitMargin, tradeFees)
	return strategy.lastTrade.Price + tax + profitMargin + tradeFees
}

func (strategy *DefaultTradingStrategy) minSellPrice(currentPrice float64) float64 {
	tax := strategy.lastTrade.Price * strategy.tax
	profitMargin := strategy.lastTrade.Price * strategy.profitMargin
	tradeFees := strategy.lastTrade.Price * strategy.fees
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minSellPrice] tax: %f, profitMargin: %f, tradeFees: %f",
		tax, profitMargin, tradeFees)
	return currentPrice + tax + profitMargin + tradeFees
}

func (strategy *DefaultTradingStrategy) getTradeAmounts(chart common.ChartService) (float64, float64) {
	baseAmount := 0.0
	quoteAmount := 0.0
	currencyPair := chart.GetCurrencyPair()
	coins, _ := chart.GetExchange().GetBalances()
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
	return baseAmount, quoteAmount
}

func (strategy *DefaultTradingStrategy) buy(chart common.ChartService) {
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.buy] $$$ BUY SIGNAL $$$")
	data := chart.GetData()
	currencyPair := chart.GetCurrencyPair()
	baseAmount, quoteAmount := strategy.getTradeAmounts(chart)
	if quoteAmount <= 0 {
		strategy.ctx.Logger.Errorf("[DefaultTradingStrategy.buy] Aborting. Out of %s funding!", currencyPair.Quote)
		return
	}
	if strategy.lastTrade.Type == "buy" {
		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.buy] Aborting. Already in a buy position.")
		return
	}
	minTradePrice := strategy.minBuyPrice()
	if data.Price <= minTradePrice {
		strategy.ctx.Logger.Debugf(
			"[DefaultTradingStrategy.buy] Aborting. Price does not meet minimum trade requirements. Price=%d, MinRequirement=%d",
			data.Price, minTradePrice)
		return
	}
	strategy.autoTradeCoin.Trades = append(strategy.autoTradeCoin.Trades, dao.Trade{
		UserID:    strategy.ctx.User.Id,
		Base:      currencyPair.Base,
		Quote:     currencyPair.Quote,
		Date:      time.Now(),
		Type:      "buy",
		Price:     data.Price,
		Amount:    baseAmount,
		ChartData: strategy.chartJSON(data)})
	strategy.autoTradeDAO.Save(strategy.autoTradeCoin)
}

func (strategy *DefaultTradingStrategy) sell(chart common.ChartService) {
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.sell] $$$ SELL SIGNAL $$$")
	data := chart.GetData()
	_, quoteAmount := strategy.getTradeAmounts(chart)
	if strategy.lastTrade.Type == "sell" {
		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.sell] Aborting. Already in a sell position.")
		return
	}
	minTradePrice := strategy.minSellPrice(data.Price)
	if data.Price <= minTradePrice {
		strategy.ctx.Logger.Debugf(
			"[DefaultTradingStrategy.sell] Aborting. Price does not meet minimum trade requirements. Price=%d, MinRequirement=%d",
			data.Price, minTradePrice)
		return
	}
	strategy.autoTradeCoin.Trades = append(strategy.autoTradeCoin.Trades, dao.Trade{
		UserID:    strategy.ctx.User.Id,
		Exchange:  data.Exchange,
		Base:      data.CurrencyPair.Base,
		Quote:     data.CurrencyPair.Quote,
		Date:      time.Now(),
		Type:      "sell",
		Price:     data.Price,
		Amount:    quoteAmount,
		ChartData: strategy.chartJSON(data)})
	strategy.autoTradeDAO.Save(strategy.autoTradeCoin)
}

func (strategy *DefaultTradingStrategy) chartJSON(data *common.ChartData) string {
	jsonChart, err := json.Marshal(data)
	if err != nil {
		strategy.ctx.Logger.Errorf("[DefaultTradingStrategy.chartJSON] Error marshalling chart state: %s", err.Error())
	}
	return string(jsonChart)
}
