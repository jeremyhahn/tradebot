package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
)

type DefaultTradingStrategyConfig struct {
	Tax                    float64
	TradeSize              float64
	ProfitMarginMin        float64
	ProfitMarginMinPercent float64
	StopLoss               float64
	StopLossPercent        float64
	RequiredBuySignals     int
	RequiredSellSignals    int
}

type DefaultTradingStrategy struct {
	name        string
	params      *common.TradingStrategyParams
	config      *DefaultTradingStrategyConfig
	buySignals  int
	sellSignals int
	common.TradingStrategy
}

func main() {}

func CreateDefaultTradingStrategy(params *common.TradingStrategyParams) (common.TradingStrategy, error) {
	expectedConfigCount := 8
	var strategyConfig *DefaultTradingStrategyConfig
	if len(params.Config) <= 0 {
		strategyConfig = &DefaultTradingStrategyConfig{
			Tax:                    .40,
			TradeSize:              1,
			ProfitMarginMin:        0,
			ProfitMarginMinPercent: .10,
			StopLoss:               0,
			StopLossPercent:        .20,
			RequiredBuySignals:     2,
			RequiredSellSignals:    2}
	} else if len(params.Config) == expectedConfigCount {
		tax, _ := strconv.ParseFloat(params.Config[0], 64)
		tradeSize, _ := strconv.ParseFloat(params.Config[1], 64)
		profitMarginMin, _ := strconv.ParseFloat(params.Config[2], 64)
		profitMarginMinPercent, _ := strconv.ParseFloat(params.Config[3], 64)
		stopLoss, _ := strconv.ParseFloat(params.Config[4], 64)
		stopLossPercent, _ := strconv.ParseFloat(params.Config[5], 64)
		requiredBuySignals, _ := strconv.ParseInt(params.Config[6], 10, 64)
		requiredSellSignals, _ := strconv.ParseInt(params.Config[7], 10, 64)
		strategyConfig = &DefaultTradingStrategyConfig{
			Tax:                    tax,
			TradeSize:              tradeSize,
			ProfitMarginMin:        profitMarginMin,
			ProfitMarginMinPercent: profitMarginMinPercent,
			StopLoss:               stopLoss,
			StopLossPercent:        stopLossPercent,
			RequiredBuySignals:     int(requiredBuySignals),
			RequiredSellSignals:    int(requiredSellSignals)}
	} else {
		errmsg := fmt.Sprintf("Invalid configuration. Expected %d items, received %d.", len(params.Config), expectedConfigCount)
		return nil, errors.New(errmsg)
	}
	strategy := &DefaultTradingStrategy{
		name:   "DefaultTradingStrategy",
		params: params,
		config: strategyConfig}
	for _, name := range strategy.GetRequiredIndicators() {
		if _, ok := params.Indicators[name]; !ok {
			return nil, errors.New(fmt.Sprintf("Strategy requires missing indicator: %s", name))
		}
	}
	return strategy, nil
}

func (strategy *DefaultTradingStrategy) GetRequiredIndicators() []string {
	return []string{"RelativeStrengthIndex", "BollingerBands", "MovingAverageConvergenceDivergence"}
}

func (strategy *DefaultTradingStrategy) GetParameters() *common.TradingStrategyParams {
	return strategy.params
}

func (strategy *DefaultTradingStrategy) GetBuySellSignals() (bool, bool, error) {
	var buy, sell bool
	strategy.countSignals()
	if strategy.buySignals == strategy.config.RequiredBuySignals {
		buy = true
		err := strategy.buy()
		if err != nil {
			return buy, sell, err
		}
	}
	if strategy.sellSignals == strategy.config.RequiredSellSignals {
		sell = true
		err := strategy.sell()
		if err != nil {
			return buy, sell, err
		}
	}
	return buy, sell, nil
}

func (strategy *DefaultTradingStrategy) CalculateFeeAndTax(price float64) (float64, float64) {
	var tax float64
	diff := price - strategy.params.LastTrade.Price
	if strategy.config.Tax > 0 && diff > 0 {
		tax = diff * strategy.config.Tax
	}
	fee := price * strategy.params.TradeFee
	return fee, tax
}

func (strategy *DefaultTradingStrategy) GetTradeAmounts() (float64, float64) {
	var baseAmount, quoteAmount float64
	if strategy.config.TradeSize > 1 {
		strategy.config.TradeSize = 1
	}
	if strategy.config.TradeSize < 0 {
		strategy.config.TradeSize = 0
	}
	for _, coin := range strategy.params.Balances {
		if coin.Currency == strategy.params.CurrencyPair.Base {
			if strategy.config.TradeSize > 0 {
				baseAmount = coin.Available * strategy.config.TradeSize
			}
		}
		if coin.Currency == strategy.params.CurrencyPair.Quote {
			if strategy.config.TradeSize > 0 {
				quoteAmount = coin.Available * strategy.config.TradeSize
			}
		}
		if baseAmount > 0 && quoteAmount > 0 {
			break
		}
	}
	return baseAmount, quoteAmount
}

func (strategy *DefaultTradingStrategy) minSellPrice() float64 {
	var profitMargin float64
	if strategy.config.ProfitMarginMinPercent > 0 {
		profitMargin = strategy.params.LastTrade.Price * strategy.config.ProfitMarginMinPercent
	} else {
		profitMargin = strategy.config.ProfitMarginMin
	}
	price := strategy.params.LastTrade.Price + profitMargin
	fee, tax := strategy.CalculateFeeAndTax(price)
	return price + fee + tax
}

func (strategy *DefaultTradingStrategy) countSignals() error {
	rsi := strategy.params.Indicators["RelativeStrengthIndex"].(indicators.RelativeStrengthIndex)
	if rsi == nil {
		return errors.New("RelativeStrengthIndex indicator required")
	}
	rsiValue := rsi.Calculate(strategy.params.NewPrice)
	if rsi.IsOverBought(rsiValue) {
		strategy.sellSignals++
	} else if rsi.IsOverSold(rsiValue) {
		strategy.buySignals++
	}
	bollinger := strategy.params.Indicators["BollingerBands"].(indicators.BollingerBands)
	if rsi == nil {
		return errors.New("BollingerBands indicator required")
	}
	upper, _, lower := bollinger.Calculate(strategy.params.NewPrice)
	if strategy.params.NewPrice > upper {
		strategy.sellSignals++
	} else if strategy.params.NewPrice < lower {
		strategy.buySignals++
	}
	//macd := strategy.params.Indicators["MovingAverageConvergenceDivergence"].(indicators.MovingAverageConvergenceDivergence)
	//value, signal, histogram := macd.Calculate(strategy.params.NewPrice)
	return nil
}

func (strategy *DefaultTradingStrategy) buy() error {
	_, quoteAmount := strategy.GetTradeAmounts()
	if quoteAmount <= 0 {
		return errors.New(fmt.Sprintf("Out of %s funding!", strategy.params.CurrencyPair.Quote))
	}
	return nil
}

func (strategy *DefaultTradingStrategy) sell() error {
	if strategy.params.LastTrade.Type == "sell" {
		return errors.New("Aborting sale. Buy position required")
	}
	minPrice := strategy.minSellPrice()
	if strategy.params.NewPrice <= minPrice {
		return errors.New(fmt.Sprintf("Aborting sale. Doesn't meet minimum trade requirements. price=%f, minRequired=%f",
			strategy.params.NewPrice, minPrice))
	}
	return nil
}

func (config *DefaultTradingStrategyConfig) ToSlice() []string {
	return []string{
		fmt.Sprintf("%f", config.Tax),
		fmt.Sprintf("%f", config.TradeSize),
		fmt.Sprintf("%f", config.ProfitMarginMin),
		fmt.Sprintf("%f", config.ProfitMarginMinPercent),
		fmt.Sprintf("%f", config.StopLoss),
		fmt.Sprintf("%f", config.StopLossPercent),
		fmt.Sprintf("%d", config.RequiredBuySignals),
		fmt.Sprintf("%d", config.RequiredSellSignals)}
}
