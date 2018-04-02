package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategySell struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategySell struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategySell struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_SellSuccess(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategySell),
		"BollingerBands":                     new(MockBBands_StrategySell),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategySell)}
	lastTrade := &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   decimal.NewFromFloat(1),
		Price:    decimal.NewFromFloat(8000)}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     decimal.NewFromFloat(16000),
		LastTrade:    lastTrade,
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
	assert.Equal(t, sell, true)
	assert.Equal(t, map[string]string{
		"RelativeStrengthIndex":              "71",
		"BollingerBands":                     "15000, 12500, 10000",
		"MovingAverageConvergenceDivergence": "25, 20, 3.25"}, data)
	assert.Equal(t, err, nil)
}

func (mrsi *MockRSI_StrategySell) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategySell) Calculate(price decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(71.0)
}

func (mrsi *MockRSI_StrategySell) IsOverBought(price decimal.Decimal) bool {
	return true
}

func (mrsi *MockRSI_StrategySell) IsOverSold(price decimal.Decimal) bool {
	return false
}

func (mrsi *MockBBands_StrategySell) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategySell) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(15000.0), decimal.NewFromFloat(12500.0), decimal.NewFromFloat(10000.0)
}

func (mrsi *MockMACD_StrategySell) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategySell) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(25), decimal.NewFromFloat(20), decimal.NewFromFloat(3.25)
}
