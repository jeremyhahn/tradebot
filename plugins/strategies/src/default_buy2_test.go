package main

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSI_StrategyBuy2 struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategyBuy2 struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategyBuy2 struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_DefaultConfig_Buy2(t *testing.T) {
	strategyIndicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyBuy2),
		"BollingerBands":                     new(MockBBands_StrategyBuy2),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyBuy2)}
	balances := []common.Coin{
		&dto.CoinDTO{
			Currency:  "BTC",
			Available: decimal.NewFromFloat(1.5)},
		&dto.CoinDTO{
			Currency:  "USD",
			Available: decimal.NewFromFloat(0.0)}}
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
		Indicators:   strategyIndicators,
		NewPrice:     decimal.NewFromFloat(11000),
		TradeFee:     decimal.NewFromFloat(.025),
		Balances:     balances,
		LastTrade:    lastTrade}

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
	assert.Equal(t, "Out of USD funding!", err.Error())
	assert.Equal(t, map[string]string{
		"RelativeStrengthIndex":              "29",
		"BollingerBands":                     "14000, 13000, 12000",
		"MovingAverageConvergenceDivergence": "25, 20, 3.25",
	}, data)
}

func (mrsi *MockRSI_StrategyBuy2) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategyBuy2) Calculate(price decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(29.0)
}

func (mrsi *MockRSI_StrategyBuy2) IsOverBought(rsiValue decimal.Decimal) bool {
	return false
}

func (mrsi *MockRSI_StrategyBuy2) IsOverSold(rsiValue decimal.Decimal) bool {
	return true
}

func (mrsi *MockBBands_StrategyBuy2) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategyBuy2) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(14000.0), decimal.NewFromFloat(13000.0), decimal.NewFromFloat(12000.0)
}

func (mrsi *MockMACD_StrategyBuy2) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategyBuy2) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(25), decimal.NewFromFloat(20), decimal.NewFromFloat(3.25)
}
