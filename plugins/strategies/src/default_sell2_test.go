package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategySell2 struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategySell2 struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategySell2 struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_SellDoesntMeetMinimumRequired(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategySell2),
		"BollingerBands":                     new(MockBBands_StrategySell2),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategySell2)}
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
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     9000,
		LastTrade:    lastTrade,
		TradeFee:     .025}

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
		"RelativeStrengthIndex":              "71.00",
		"BollingerBands":                     "8000.00, 7000.00, 6000.00",
		"MovingAverageConvergenceDivergence": "25.00, 20.00, 3.25"}, data)
	assert.Equal(t, "Aborting sale. Doesn't meet minimum trade requirements. price=9000.000000, minRequired=11675.000000", err.Error())
}

func (mrsi *MockRSI_StrategySell2) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategySell2) Calculate(price float64) float64 {
	return 71.0
}

func (mrsi *MockRSI_StrategySell2) IsOverBought(price float64) bool {
	return true
}

func (mrsi *MockRSI_StrategySell2) IsOverSold(price float64) bool {
	return false
}

func (mrsi *MockBBands_StrategySell2) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategySell2) Calculate(price float64) (float64, float64, float64) {
	return 8000.0, 7000.0, 6000.0
}

func (mrsi *MockMACD_StrategySell2) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategySell2) Calculate(price float64) (float64, float64, float64) {
	return 25, 20, 3.25
}
