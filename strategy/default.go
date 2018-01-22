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
	Name          string
	ctx           *common.Context
	autoTradeDAO  dao.IAutoTradeDAO
	signalLogDAO  dao.ISignalLogDAO
	autoTradeCoin dao.IAutoTradeCoin
	config        *DefaultTradingStrategyConfig
	lastTrade     *dao.Trade
	service.TradingStrategy
}

type DefaultTradingStrategyConfig struct {
	rsiOverSold            float64
	rsiOverBought          float64
	tax                    float64
	profitMarginMin        float64
	profitMarginMinPercent float64
	stopLoss               float64
	stopLossPercent        float64
	requiredBuySignals     int
	requiredSellSignals    int
	tradeSize              float64
	tradeSizePercent       float64
}

func NewDefaultTradingStrategy(ctx *common.Context, autoTradeCoin dao.IAutoTradeCoin,
	autoTradeDAO dao.IAutoTradeDAO, signalLogDAO dao.ISignalLogDAO) *DefaultTradingStrategy {
	return &DefaultTradingStrategy{
		Name:          "DefaultTradingStrategy",
		ctx:           ctx,
		autoTradeCoin: autoTradeCoin,
		autoTradeDAO:  autoTradeDAO,
		signalLogDAO:  signalLogDAO,
		config: &DefaultTradingStrategyConfig{
			rsiOverSold:            30,
			rsiOverBought:          70,
			tax:                    0,
			stopLoss:               0,
			stopLossPercent:        .20,
			profitMarginMin:        0,
			profitMarginMinPercent: .10,
			tradeSizePercent:       0,
			requiredBuySignals:     2,
			requiredSellSignals:    2}}
}

func CreateDefaultTradingStrategy(ctx *common.Context, autoTradeCoin dao.IAutoTradeCoin,
	autoTradeDAO dao.IAutoTradeDAO, signalLogDAO dao.ISignalLogDAO, config *DefaultTradingStrategyConfig) *DefaultTradingStrategy {
	return &DefaultTradingStrategy{
		Name:          "DefaultTradingStrategy",
		ctx:           ctx,
		autoTradeCoin: autoTradeCoin,
		autoTradeDAO:  autoTradeDAO,
		signalLogDAO:  signalLogDAO,
		config:        config}
}

func (strategy *DefaultTradingStrategy) OnPriceChange(chart common.ChartService) {

	data := chart.GetData()
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPriceChange] ChartData: %+v\n", data)

	buySignals, sellSignals := strategy.countSignals(data)
	strategy.lastTrade = strategy.autoTradeDAO.GetLastTrade(strategy.autoTradeCoin)

	if buySignals == strategy.config.requiredBuySignals {
		strategy.buy(chart)
		return
	}

	if sellSignals == strategy.config.requiredSellSignals {
		strategy.sell(chart)
		return
	}

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] buySignals=%d, sellSignals=%d", buySignals, sellSignals)
}

func (strategy *DefaultTradingStrategy) countSignals(data *common.ChartData) (int, int) {
	var buySignals int
	var sellSignals int
	if data.RSILive < strategy.config.rsiOverSold {
		buySignals++
		strategy.signalLogDAO.Save(&dao.SignalLog{
			UserID:     strategy.ctx.User.Id,
			Date:       time.Now(),
			Name:       "RSI",
			Type:       "buy",
			Price:      data.Price,
			SignalData: strconv.FormatFloat(data.RSILive, 'f', 8, 64)})
	} else if data.RSILive > strategy.config.rsiOverBought {
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

func (strategy *DefaultTradingStrategy) minSellPrice(currentPrice, tradingFee float64) float64 {
	var price, tax, profitMargin float64
	if strategy.lastTrade.Price > currentPrice {
		price = strategy.lastTrade.Price
	} else {
		price = currentPrice
	}
	if strategy.config.profitMarginMinPercent > 0 {
		profitMargin = price * strategy.config.profitMarginMinPercent
	} else {
		profitMargin = strategy.config.profitMarginMin
	}
	if strategy.config.tax > 0 {
		tax = (price + profitMargin) * strategy.config.tax
	}
	fee := (price + profitMargin) * tradingFee
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minSellPrice] price: %f, profitMargin: %f, fee: %f,tax: %f",
		price, profitMargin, fee, tax)
	return price + profitMargin + fee + tax
}

func (strategy *DefaultTradingStrategy) getTradeAmounts(chart common.ChartService) (float64, float64) {
	var baseAmount, quoteAmount float64
	currencyPair := chart.GetCurrencyPair()
	coins, _ := chart.GetExchange().GetBalances()
	for _, coin := range coins {
		if coin.Currency == currencyPair.Base {
			if strategy.config.tradeSizePercent > 0 {
				baseAmount = coin.Available * strategy.config.tradeSizePercent
			} else {
				baseAmount = coin.Available
			}
		}
		if coin.Currency == currencyPair.Quote {
			if strategy.config.tradeSizePercent > 0 {
				quoteAmount = coin.Available * strategy.config.tradeSizePercent
			} else {
				quoteAmount = coin.Available
			}
		}
		if baseAmount > 0 && quoteAmount > 0 {
			break
		}
	}
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.getTradeAmounts] Trading funds - baseAmount: %f, quoteAmount: %f",
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
	strategy.autoTradeCoin.AddTrade(&dao.Trade{
		UserID:    strategy.ctx.User.Id,
		Exchange:  data.Exchange,
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
		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.sell] Aborting. Buy position required.")
		return
	}
	minTradePrice := strategy.minSellPrice(data.Price, chart.GetExchange().GetTradingFee())
	if data.Price <= minTradePrice {
		strategy.ctx.Logger.Debugf(
			"[DefaultTradingStrategy.sell] Aborting. Price does not meet minimum trade requirements. Price=%d, MinRequirement=%d",
			data.Price, minTradePrice)
		return
	}
	strategy.autoTradeCoin.AddTrade(&dao.Trade{
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
