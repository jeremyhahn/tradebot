package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategyCalculate struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategyCalculate struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategyCalculate struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_Calculate(t *testing.T) {
	indicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyCalculate),
		"BollingerBands":                     new(MockBBands_StrategyCalculate),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyCalculate)}
	lastTrade := &common.Trade{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    10000}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     nil,
		Indicators:   indicators,
		NewPrice:     13000,
		LastTrade:    lastTrade,
		TradeFee:     .025}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	buy, sell, data, err := strategy.Analyze()
	assert.Equal(t, false, buy)
	assert.Equal(t, true, sell)
	assert.Equal(t, map[string]string{
		"RelativeStrengthIndex":              "79.00",
		"BollingerBands":                     "10000.00, 9000.00, 8000.00",
		"MovingAverageConvergenceDivergence": "25.00, 20.00, 3.25"}, data)
	assert.Equal(t, nil, err)

	minPrice := strategy.minSellPrice()
	assert.Equal(t, 11675.0, minPrice)

	fees, tax := strategy.CalculateFeeAndTax(params.NewPrice)
	assert.Equal(t, 325.0, fees)
	assert.Equal(t, 1200.0, tax)
}

func (mrsi *MockRSI_StrategyCalculate) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategyCalculate) Calculate(price float64) float64 {
	return 79.0
}

func (mrsi *MockRSI_StrategyCalculate) IsOverBought(rsiValue float64) bool {
	return true
}

func (mrsi *MockRSI_StrategyCalculate) IsOverSold(rsiValue float64) bool {
	return false
}

func (mrsi *MockBBands_StrategyCalculate) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategyCalculate) Calculate(price float64) (float64, float64, float64) {
	return 10000.0, 9000.0, 8000.0
}

func (mrsi *MockMACD_StrategyCalculate) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategyCalculate) Calculate(price float64) (float64, float64, float64) {
	return 25, 20, 3.25
}
