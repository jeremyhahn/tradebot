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
	lastTrade := &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   decimal.NewFromFloat(1),
		Price:    decimal.NewFromFloat(10000)}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     nil,
		Indicators:   indicators,
		NewPrice:     decimal.NewFromFloat(13000),
		LastTrade:    lastTrade,
		TradeFee:     decimal.NewFromFloat(.025)}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	buy, sell, data, err := strategy.Analyze()
	assert.Equal(t, false, buy)
	assert.Equal(t, true, sell)
	assert.Equal(t, map[string]string{
		"RelativeStrengthIndex":              "79",
		"BollingerBands":                     "10000, 9000, 8000",
		"MovingAverageConvergenceDivergence": "25, 20, 3.25"}, data)
	assert.Equal(t, nil, err)

	minPrice := strategy.minSellPrice()
	assert.Equal(t, decimal.NewFromFloat(11675.0).String(), minPrice.String())

	fees, tax := strategy.CalculateFeeAndTax(params.NewPrice)
	assert.Equal(t, decimal.NewFromFloat(325.0).String(), fees.String())
	assert.Equal(t, decimal.NewFromFloat(1200.0).String(), tax.String())
}

func (mrsi *MockRSI_StrategyCalculate) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategyCalculate) Calculate(price decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(79.0)
}

func (mrsi *MockRSI_StrategyCalculate) IsOverBought(rsiValue decimal.Decimal) bool {
	return true
}

func (mrsi *MockRSI_StrategyCalculate) IsOverSold(rsiValue decimal.Decimal) bool {
	return false
}

func (mrsi *MockBBands_StrategyCalculate) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategyCalculate) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(10000.0), decimal.NewFromFloat(9000.0), decimal.NewFromFloat(8000.0)
}

func (mrsi *MockMACD_StrategyCalculate) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategyCalculate) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(25), decimal.NewFromFloat(20), decimal.NewFromFloat(3.25)
}
