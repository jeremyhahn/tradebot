package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
)

type DefaultTradingStrategyConfig struct {
	Tax                    decimal.Decimal
	TradeSize              decimal.Decimal
	ProfitMarginMin        decimal.Decimal
	ProfitMarginMinPercent decimal.Decimal
	StopLoss               decimal.Decimal
	StopLossPercent        decimal.Decimal
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

func main() {
}

func CreateDefaultTradingStrategy(params *common.TradingStrategyParams) (common.TradingStrategy, error) {
	expectedConfigCount := 8
	var strategyConfig *DefaultTradingStrategyConfig
	if params.Config == nil {
		strategyConfig = &DefaultTradingStrategyConfig{
			Tax:                    decimal.NewFromFloat(.40),
			TradeSize:              decimal.NewFromFloat(1),
			ProfitMarginMin:        decimal.NewFromFloat(0),
			ProfitMarginMinPercent: decimal.NewFromFloat(.10),
			StopLoss:               decimal.NewFromFloat(0),
			StopLossPercent:        decimal.NewFromFloat(.20),
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
			Tax:                    decimal.NewFromFloat(tax),
			TradeSize:              decimal.NewFromFloat(tradeSize),
			ProfitMarginMin:        decimal.NewFromFloat(profitMarginMin),
			ProfitMarginMinPercent: decimal.NewFromFloat(profitMarginMinPercent),
			StopLoss:               decimal.NewFromFloat(stopLoss),
			StopLossPercent:        decimal.NewFromFloat(stopLossPercent),
			RequiredBuySignals:     int(requiredBuySignals),
			RequiredSellSignals:    int(requiredSellSignals)}
	} else {
		errmsg := fmt.Sprintf("Invalid configuration. Expected %d items, received %d (%s)",
			expectedConfigCount, len(params.Config), strings.Join(params.Config, ","))
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

func (strategy *DefaultTradingStrategy) Analyze() (bool, bool, map[string]string, error) {
	var buy, sell bool
	signalData, err := strategy.countSignals()
	if err != nil {
		return buy, sell, signalData, err
	}
	if strategy.buySignals == strategy.config.RequiredBuySignals {
		buy = true
		err := strategy.buy()
		if err != nil {
			return buy, sell, signalData, err
		}
	}
	if strategy.sellSignals == strategy.config.RequiredSellSignals {
		sell = true
		err := strategy.sell()
		if err != nil {
			return buy, sell, signalData, err
		}
	}
	return buy, sell, signalData, nil
}

func (strategy *DefaultTradingStrategy) CalculateFeeAndTax(price decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	var tax decimal.Decimal
	diff := price.Sub(strategy.params.LastTrade.GetPrice())
	zero := decimal.NewFromFloat(0)
	if strategy.config.Tax.GreaterThan(zero) && diff.GreaterThan(zero) {
		tax = diff.Mul(strategy.config.Tax)
	}
	fee := price.Mul(strategy.params.TradeFee)
	return fee, tax
}

func (strategy *DefaultTradingStrategy) GetTradeAmounts() (decimal.Decimal, decimal.Decimal) {
	var baseAmount, quoteAmount decimal.Decimal
	zero := decimal.NewFromFloat(0)
	one := decimal.NewFromFloat(1)
	if strategy.config.TradeSize.GreaterThan(one) {
		strategy.config.TradeSize = one
	}
	if strategy.config.TradeSize.LessThan(zero) {
		strategy.config.TradeSize = zero
	}
	for _, coin := range strategy.params.Balances {
		if coin.GetCurrency() == strategy.params.CurrencyPair.Base {
			if strategy.config.TradeSize.GreaterThan(zero) {
				baseAmount = coin.GetAvailable().Mul(strategy.config.TradeSize)
			}
		}
		if coin.GetCurrency() == strategy.params.CurrencyPair.Quote {
			if strategy.config.TradeSize.GreaterThan(zero) {
				quoteAmount = coin.GetAvailable().Mul(strategy.config.TradeSize)
			}
		}
		if baseAmount.GreaterThan(zero) && quoteAmount.GreaterThan(zero) {
			break
		}
	}
	return baseAmount, quoteAmount
}

func (strategy *DefaultTradingStrategy) minSellPrice() decimal.Decimal {
	var profitMargin decimal.Decimal
	if strategy.config.ProfitMarginMinPercent.GreaterThan(decimal.NewFromFloat(0)) {
		profitMargin = strategy.params.LastTrade.GetPrice().Mul(strategy.config.ProfitMarginMinPercent)
	} else {
		profitMargin = strategy.config.ProfitMarginMin
	}
	price := strategy.params.LastTrade.GetPrice().Add(profitMargin)
	fee, tax := strategy.CalculateFeeAndTax(price)
	return price.Add(fee).Add(tax)
}

func (strategy *DefaultTradingStrategy) countSignals() (map[string]string, error) {
	signalData := make(map[string]string, len(strategy.params.Indicators))

	rsi := strategy.params.Indicators["RelativeStrengthIndex"].(indicators.RelativeStrengthIndex)
	if rsi == nil {
		return nil, errors.New("RelativeStrengthIndex indicator required")
	}
	rsiValue := rsi.Calculate(strategy.params.NewPrice)
	if rsi.IsOverBought(rsiValue) {
		strategy.sellSignals++
	} else if rsi.IsOverSold(rsiValue) {
		strategy.buySignals++
	}
	signalData[rsi.GetName()] = fmt.Sprintf("%s", rsiValue)
	bollinger := strategy.params.Indicators["BollingerBands"].(indicators.BollingerBands)
	if rsi == nil {
		return nil, errors.New("BollingerBands indicator required")
	}
	upper, middle, lower := bollinger.Calculate(strategy.params.NewPrice)
	if strategy.params.NewPrice.GreaterThan(upper) {
		strategy.sellSignals++
	} else if strategy.params.NewPrice.LessThan(lower) {
		strategy.buySignals++
	}
	signalData[bollinger.GetName()] = fmt.Sprintf("%s, %s, %s", upper, middle, lower)

	macd := strategy.params.Indicators["MovingAverageConvergenceDivergence"].(indicators.MovingAverageConvergenceDivergence)
	value, signal, histogram := macd.Calculate(strategy.params.NewPrice)
	signalData[macd.GetName()] = fmt.Sprintf("%s, %s, %s", value, signal, histogram)

	return signalData, nil
}

func (strategy *DefaultTradingStrategy) buy() error {
	_, quoteAmount := strategy.GetTradeAmounts()
	if quoteAmount.LessThanOrEqual(decimal.NewFromFloat(0)) {
		return errors.New(fmt.Sprintf("Out of %s funding!", strategy.params.CurrencyPair.Quote))
	}
	return nil
}

func (strategy *DefaultTradingStrategy) sell() error {
	if strategy.params.LastTrade.GetType() == "sell" {
		return errors.New("Aborting sale. Buy position required")
	}
	minPrice := strategy.minSellPrice()
	if strategy.params.NewPrice.LessThanOrEqual(minPrice) {
		return errors.New(fmt.Sprintf("Aborting sale. Doesn't meet minimum trade requirements. price=%s, minRequired=%s",
			strategy.params.NewPrice, minPrice))
	}
	return nil
}

func (config *DefaultTradingStrategyConfig) ToSlice() []string {
	return []string{
		fmt.Sprintf("%s", config.Tax),
		fmt.Sprintf("%s", config.TradeSize),
		fmt.Sprintf("%s", config.ProfitMarginMin),
		fmt.Sprintf("%s", config.ProfitMarginMinPercent),
		fmt.Sprintf("%s", config.StopLoss),
		fmt.Sprintf("%s", config.StopLossPercent),
		fmt.Sprintf("%d", config.RequiredBuySignals),
		fmt.Sprintf("%d", config.RequiredSellSignals)}
}
