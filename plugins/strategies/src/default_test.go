package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRelativeStrengthIndex struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBollingerBands struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMovingAverageConvergenceDivergence struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     decimal.NewFromFloat(11000),
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025)}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	requiredIndicators := strategy.GetRequiredIndicators()
	assert.Equal(t, "RelativeStrengthIndex", requiredIndicators[0])
	assert.Equal(t, "BollingerBands", requiredIndicators[1])
	assert.Equal(t, "MovingAverageConvergenceDivergence", requiredIndicators[2])

	buy, sell, data, err := strategy.Analyze()
	assert.Equal(t, buy, false)
	assert.Equal(t, sell, false)
	assert.Equal(t, map[string]string{
		"MovingAverageConvergenceDivergence": "25, 20, 3.25",
		"RelativeStrengthIndex":              "15000, 12500, 10000"}, data)
	assert.Equal(t, err, nil)
}

func TestDefaultTradingStrategy_CustomTradeSize_Percentage(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(0),
		ProfitMarginMinPercent: decimal.NewFromFloat(.10),
		TradeSize:              decimal.NewFromFloat(.1),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)
	strategy := s.(*DefaultTradingStrategy)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base.String(), decimal.NewFromFloat(0.2).String())
	assert.Equal(t, quote.String(), decimal.NewFromFloat(2000.0).String())
}

func TestDefaultTradingStrategy_CustomTradeSize_AvailableBalance(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(0),
		ProfitMarginMinPercent: decimal.NewFromFloat(.10),
		TradeSize:              decimal.NewFromFloat(1),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base.String(), decimal.NewFromFloat(2.0).String())
	assert.Equal(t, quote.String(), decimal.NewFromFloat(20000.0).String())
}

func TestDefaultTradingStrategy_CustomTradeSize_Zero(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(0),
		ProfitMarginMinPercent: decimal.NewFromFloat(.10),
		TradeSize:              decimal.NewFromFloat(0),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base.String(), decimal.NewFromFloat(0.0).String())
	assert.Equal(t, quote.String(), decimal.NewFromFloat(0.0).String())
}

func TestDefaultTradingStrategy_CustomTradeSize_GreaterThanOne(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(0),
		ProfitMarginMinPercent: decimal.NewFromFloat(.10),
		TradeSize:              decimal.NewFromFloat(2),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base.String(), decimal.NewFromFloat(2.0).String())
	assert.Equal(t, quote.String(), decimal.NewFromFloat(20000.0).String())
}

func TestDefaultTradingStrategy_CustomTradeSize_LessThanZero(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(0),
		ProfitMarginMinPercent: decimal.NewFromFloat(.10),
		TradeSize:              decimal.NewFromFloat(-2),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base.String(), decimal.NewFromFloat(0.0).String())
	assert.Equal(t, quote.String(), decimal.NewFromFloat(0.0).String())
}

func (mrsi *MockRelativeStrengthIndex) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRelativeStrengthIndex) Calculate(price decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(31.0)
}

func (mrsi *MockRelativeStrengthIndex) IsOverBought(price decimal.Decimal) bool {
	return false
}

func (mrsi *MockRelativeStrengthIndex) IsOverSold(price decimal.Decimal) bool {
	return false
}

func (mrsi *MockBollingerBands) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockBollingerBands) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(15000.0), decimal.NewFromFloat(12500.0), decimal.NewFromFloat(10000.0)
}

func (mrsi *MockMovingAverageConvergenceDivergence) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMovingAverageConvergenceDivergence) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(25), decimal.NewFromFloat(20), decimal.NewFromFloat(3.25)
}
