package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategyBuy struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategyBuy struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategyBuy struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_Buy(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyBuy),
		"BollingerBands":                     new(MockBBands_StrategyBuy),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyBuy)}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     11000,
		TradeFee:     .025,
		LastTrade: &dto.TradeDTO{
			Id:       1,
			ChartId:  1,
			Base:     "BTC",
			Quote:    "USD",
			Exchange: "gdax",
			Type:     "buy",
			Amount:   1,
			Price:    8000}}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	requiredIndicators := strategy.GetRequiredIndicators()
	assert.Equal(t, "RelativeStrengthIndex", requiredIndicators[0])
	assert.Equal(t, "BollingerBands", requiredIndicators[1])
	assert.Equal(t, "MovingAverageConvergenceDivergence", requiredIndicators[2])

	buy, sell, data, err := strategy.Analyze()
	assert.Equal(t, true, buy)
	assert.Equal(t, false, sell)
	assert.Equal(t, nil, err)
	assert.Equal(t, map[string]string{
		"MovingAverageConvergenceDivergence": "25.00, 20.00, 3.25",
		"RelativeStrengthIndex":              "29.00",
		"BollingerBands":                     "14000.00, 13000.00, 12000.00"}, data)
}

func (mrsi *MockRSI_StrategyBuy) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategyBuy) Calculate(price float64) float64 {
	return 29.0
}

func (mrsi *MockRSI_StrategyBuy) IsOverBought(rsiValue float64) bool {
	return false
}

func (mrsi *MockRSI_StrategyBuy) IsOverSold(rsiValue float64) bool {
	return true
}

func (mrsi *MockBBands_StrategyBuy) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategyBuy) Calculate(price float64) (float64, float64, float64) {
	return 14000.0, 13000.0, 12000.0
}

func (mrsi *MockMACD_StrategyBuy) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategyBuy) Calculate(price float64) (float64, float64, float64) {
	return 25, 20, 3.25
}
