package strategy

import (
	"encoding/json"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
)

type DefaultTradingStrategy struct {
	Name          string
	ctx           *common.Context
	chartService  common.ChartService
	tradeService  common.TradeService
	profitService common.ProfitService
	config        *DefaultTradingStrategyConfig
	lastTrade     *common.Trade
	Chart         chan common.ChartService
	common.TradingStrategy
}

type DefaultTradingStrategyConfig struct {
	config                 *DefaultTradingStrategyConfig
	tradeSize              float64
	tax                    float64
	profitMarginMin        float64
	profitMarginMinPercent float64
	stopLoss               float64
	stopLossPercent        float64
	requiredBuySignals     int
	requiredSellSignals    int
}

func NewDefaultTradingStrategy(ctx *common.Context, chartService common.ChartService,
	tradeService common.TradeService, profitService common.ProfitService) common.TradingStrategy {
	return &DefaultTradingStrategy{
		Name:          "DefaultTradingStrategy",
		Chart:         make(chan common.ChartService),
		ctx:           ctx,
		chartService:  chartService,
		tradeService:  tradeService,
		profitService: profitService,
		config: &DefaultTradingStrategyConfig{
			tax:                    .40,
			stopLoss:               0,
			stopLossPercent:        .20,
			profitMarginMin:        0,
			profitMarginMinPercent: .10,
			tradeSize:              1,
			requiredBuySignals:     2,
			requiredSellSignals:    2}}
}

func CreateDefaultTradingStrategy(ctx *common.Context, chartService common.ChartService,
	tradeService common.TradeService, profitService common.ProfitService,
	config *DefaultTradingStrategyConfig) common.TradingStrategy {
	return &DefaultTradingStrategy{
		Name:          "DefaultTradingStrategy",
		ctx:           ctx,
		chartService:  chartService,
		tradeService:  tradeService,
		profitService: profitService,
		config:        config}
}

func (strategy *DefaultTradingStrategy) OnPriceChange(chart common.ChartService) {

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPriceChange] ChartService: %+v\n", chart)

	buySignals, sellSignals := strategy.countSignals(chart.GetPrice())
	strategy.lastTrade = strategy.tradeService.GetLastTrade(strategy.chartService.GetChart())

	if buySignals == strategy.config.requiredBuySignals {
		strategy.buy(chart)
		return
	}

	if sellSignals == strategy.config.requiredSellSignals {
		strategy.sell(chart)
		return
	}

	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.OnPeriodChange] buySignals=%d, sellSignals=%d", buySignals, sellSignals)
	return
}

func (strategy *DefaultTradingStrategy) GetRequiredIndicators() []string {
	return []string{"RSI", "BollingerBands", "MACD"}
}

func (strategy *DefaultTradingStrategy) countSignals(price float64) (int, int) {
	var buySignals int
	var sellSignals int
	rsi := strategy.chartService.GetIndicator("RSI").(*indicators.RelativeStrengthIndex)
	rsiValue := rsi.Calculate(price)
	if rsi.IsOverBought(rsiValue) {
		sellSignals++
	} else if rsi.IsOverSold(rsiValue) {
		buySignals++
	}
	bollinger := strategy.chartService.GetIndicator("BBands").(*indicators.BollingerBands)
	upper, _, lower := bollinger.Calculate(price)
	if price > upper {
		sellSignals++
	} else if price < lower {
		buySignals++
	}
	return buySignals, sellSignals
}

func (strategy *DefaultTradingStrategy) calculateFeeAndTax(price, tradingFee float64) (float64, float64) {
	var tax float64
	diff := price - strategy.lastTrade.Price
	if strategy.config.tax > 0 && diff > 0 {
		tax = diff * strategy.config.tax
	}
	fee := price * tradingFee
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.calculateFeeAndTax] lastTradePrice: %f, price: %f, fee: %f,tax: %f",
		strategy.lastTrade.Price, price, fee, tax)
	return fee, tax
}

func (strategy *DefaultTradingStrategy) minSellPrice(tradingFee float64) float64 {
	var price, profitMargin, fee, tax float64
	if strategy.config.profitMarginMinPercent > 0 {
		profitMargin = strategy.lastTrade.Price * strategy.config.profitMarginMinPercent
	} else {
		profitMargin = strategy.config.profitMarginMin
	}
	price = strategy.lastTrade.Price + profitMargin
	fee, tax = strategy.calculateFeeAndTax(price, tradingFee)
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.minSellPrice] lastTradePrice: %f, price: %f, fee: %f,tax: %f",
		strategy.lastTrade.Price, price, fee, tax)
	return price + fee + tax
}

func (strategy *DefaultTradingStrategy) getTradeAmounts(chart common.ChartService) (float64, float64) {
	var baseAmount, quoteAmount float64
	currencyPair := chart.GetCurrencyPair()
	coins, _ := chart.GetExchange().GetBalances()
	for _, coin := range coins {
		if coin.Currency == currencyPair.Base {
			if strategy.config.tradeSize > 0 {
				baseAmount = coin.Available * strategy.config.tradeSize
			} else {
				baseAmount = coin.Available
			}
		}
		if coin.Currency == currencyPair.Quote {
			if strategy.config.tradeSize > 0 {
				quoteAmount = coin.Available * strategy.config.tradeSize
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
	strategy.tradeService.Save(&common.Trade{
		UserID:    strategy.ctx.User.Id,
		Exchange:  strategy.chartService.GetExchange().GetName(),
		Base:      currencyPair.Base,
		Quote:     currencyPair.Quote,
		Date:      time.Now(),
		Type:      "buy",
		Price:     chart.GetPrice(),
		Amount:    baseAmount,
		ChartData: chart.ToJson()})
}

func (strategy *DefaultTradingStrategy) sell(chart common.ChartService) {
	strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.sell] $$$ SELL SIGNAL $$$")
	price := chart.GetPrice()
	baseAmount, _ := strategy.getTradeAmounts(chart)
	if strategy.lastTrade.Type == "sell" {
		strategy.ctx.Logger.Debugf("[DefaultTradingStrategy.sell] Aborting. Buy position required.")
		return
	}
	tradeFee := chart.GetExchange().GetTradingFee()
	minPrice := strategy.minSellPrice(tradeFee)
	if price <= minPrice {
		strategy.ctx.Logger.Debugf(
			"[DefaultTradingStrategy.sell] Aborting. Does not meet minimum trade requirements. Price=%f, MinPrice=%f",
			price, minPrice)
		return
	}
	fee, tax := strategy.calculateFeeAndTax(price, tradeFee)
	trade := &common.Trade{
		UserID:    strategy.ctx.User.Id,
		Exchange:  chart.GetExchange().GetName(),
		Base:      chart.GetCurrencyPair().Base,
		Quote:     chart.GetCurrencyPair().Quote,
		Date:      time.Now(),
		Type:      "sell",
		Price:     price,
		Amount:    baseAmount,
		ChartData: chart.ToJson()}
	strategy.tradeService.Save(trade)

	strategy.profitService.Save(&common.Profit{
		UserID:   strategy.ctx.User.Id,
		TradeID:  trade.ID,
		Quantity: baseAmount,
		Bought:   strategy.lastTrade.Price,
		Sold:     price,
		Fee:      fee,
		Tax:      tax,
		Total:    price - strategy.lastTrade.Price - fee - tax})
}

func (strategy *DefaultTradingStrategy) chartJSON(data *common.ChartData) string {
	jsonChart, err := json.Marshal(data)
	if err != nil {
		strategy.ctx.Logger.Errorf("[DefaultTradingStrategy.chartJSON] Error marshalling chart state: %s", err.Error())
	}
	return string(jsonChart)
}
