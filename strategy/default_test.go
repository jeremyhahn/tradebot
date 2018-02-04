package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
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
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     11000,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	requiredIndicators := strategy.GetRequiredIndicators()
	assert.Equal(t, "RelativeStrengthIndex", requiredIndicators[0])
	assert.Equal(t, "BollingerBands", requiredIndicators[1])
	assert.Equal(t, "MovingAverageConvergenceDivergence", requiredIndicators[2])

	buy, sell, err := strategy.GetBuySellSignals()
	assert.Equal(t, buy, false)
	assert.Equal(t, sell, false)
	assert.Equal(t, err, nil)
}

func TestDefaultTradingStrategy_CustomTradeSize_Percentage(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        0,
		ProfitMarginMinPercent: .10,
		TradeSize:              .1,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)
	strategy := s.(*DefaultTradingStrategy)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base, 0.2)
	assert.Equal(t, quote, 2000.0)
}

func TestDefaultTradingStrategy_CustomTradeSize_AvailableBalance(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        0,
		ProfitMarginMinPercent: .10,
		TradeSize:              1,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base, 2.0)
	assert.Equal(t, quote, 20000.0)
}

func TestDefaultTradingStrategy_CustomTradeSize_Zero(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        0,
		ProfitMarginMinPercent: .10,
		TradeSize:              0,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base, 0.0)
	assert.Equal(t, quote, 0.0)
}

func TestDefaultTradingStrategy_CustomTradeSize_GreaterThanOne(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        0,
		ProfitMarginMinPercent: .10,
		TradeSize:              2,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base, 2.0)
	assert.Equal(t, quote, 20000.0)
}

func TestDefaultTradingStrategy_CustomTradeSize_LessThanZero(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRelativeStrengthIndex),
		"BollingerBands":                     new(MockBollingerBands),
		"MovingAverageConvergenceDivergence": new(MockMovingAverageConvergenceDivergence)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        0,
		ProfitMarginMinPercent: .10,
		TradeSize:              -2,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   strategyIndicators,
		LastTrade:    helper.CreateLastTrade(),
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	strategy := s.(*DefaultTradingStrategy)
	assert.Equal(t, nil, err)

	base, quote := strategy.GetTradeAmounts()
	assert.Equal(t, base, 0.0)
	assert.Equal(t, quote, 0.0)
}

func (mrsi *MockRelativeStrengthIndex) Calculate(price float64) float64 {
	return 31.0
}

func (mrsi *MockRelativeStrengthIndex) IsOverBought(price float64) bool {
	return false
}

func (mrsi *MockRelativeStrengthIndex) IsOverSold(price float64) bool {
	return false
}

func (mrsi *MockBollingerBands) Calculate(price float64) (float64, float64, float64) {
	return 15000.0, 12500.0, 10000.0
}
