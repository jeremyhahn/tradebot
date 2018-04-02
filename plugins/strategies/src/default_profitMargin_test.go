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

type MockRSI_StrategyProfitMargin struct {
	indicators.RelativeStrengthIndex
	mock.Mock
}

type MockBBands_StrategyProfitMargin struct {
	indicators.BollingerBands
	mock.Mock
}

type MockMACD_StrategyProfitMargin struct {
	indicators.MovingAverageConvergenceDivergence
	mock.Mock
}

func TestDefaultTradingStrategy_CustomTradeSize_ProfitMarginFixed1(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	indicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyProfitMargin),
		"BollingerBands":                     new(MockBBands_StrategyProfitMargin),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyProfitMargin)}
	lastTrade := &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   decimal.NewFromFloat(1),
		Price:    decimal.NewFromFloat(10000)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		TradeSize:              decimal.NewFromFloat(1),
		ProfitMarginMin:        decimal.NewFromFloat(10000),
		ProfitMarginMinPercent: decimal.NewFromFloat(0),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   indicators,
		LastTrade:    lastTrade,
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)
	strategy := s.(*DefaultTradingStrategy)

	_, sell, data, _ := strategy.Analyze()
	assert.Equal(t, sell, true)
	assert.Equal(t, map[string]string{
		"BollingerBands":                     "10000, 9000, 8000",
		"MovingAverageConvergenceDivergence": "25, 20, 3.25",
		"RelativeStrengthIndex":              "79"}, data)

	minSellPrice := strategy.minSellPrice()
	assert.Equal(t, minSellPrice.String(), decimal.NewFromFloat(24500.0).String())

	fees, tax := strategy.CalculateFeeAndTax(minSellPrice)
	assert.Equal(t, fees.String(), decimal.NewFromFloat(612.5).String())
	assert.Equal(t, tax.String(), decimal.NewFromFloat(5800.0).String())
	assert.Equal(t, "Aborting sale. Doesn't meet minimum trade requirements. price=11000, minRequired=24500", strategy.sell().Error())
}

func TestDefaultTradingStrategy_CustomTradeSize_ProfitMarginFixed2(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	indicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyProfitMargin),
		"BollingerBands":                     new(MockBBands_StrategyProfitMargin),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyProfitMargin)}
	lastTrade := &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   decimal.NewFromFloat(1),
		Price:    decimal.NewFromFloat(10000)}
	config := &DefaultTradingStrategyConfig{
		Tax:                    decimal.NewFromFloat(.40),
		StopLoss:               decimal.NewFromFloat(0),
		StopLossPercent:        decimal.NewFromFloat(.20),
		ProfitMarginMin:        decimal.NewFromFloat(200),
		ProfitMarginMinPercent: decimal.NewFromFloat(0),
		TradeSize:              decimal.NewFromFloat(1),
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &common.TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     decimal.NewFromFloat(11000),
		Indicators:   indicators,
		LastTrade:    lastTrade,
		TradeFee:     decimal.NewFromFloat(.025),
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	_, sell, _, _ := strategy.Analyze()
	assert.Equal(t, sell, true)

	fees, tax := strategy.CalculateFeeAndTax(params.NewPrice)
	assert.Equal(t, strategy.minSellPrice().String(), decimal.NewFromFloat(10535.0).String())
	assert.Equal(t, fees.String(), decimal.NewFromFloat(275.0).String())
	assert.Equal(t, tax.String(), decimal.NewFromFloat(400.0).String())
	assert.Equal(t, nil, strategy.sell())
}

func (mrsi *MockRSI_StrategyProfitMargin) GetName() string {
	return "RelativeStrengthIndex"
}

func (mrsi *MockRSI_StrategyProfitMargin) Calculate(price decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(79.0)
}

func (mrsi *MockRSI_StrategyProfitMargin) IsOverBought(rsiValue decimal.Decimal) bool {
	return true
}

func (mrsi *MockRSI_StrategyProfitMargin) IsOverSold(rsiValue decimal.Decimal) bool {
	return false
}

func (mrsi *MockBBands_StrategyProfitMargin) GetName() string {
	return "BollingerBands"
}

func (mrsi *MockBBands_StrategyProfitMargin) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(10000.0), decimal.NewFromFloat(9000.0), decimal.NewFromFloat(8000.0)
}

func (mrsi *MockMACD_StrategyProfitMargin) GetName() string {
	return "MovingAverageConvergenceDivergence"
}

func (mrsi *MockMACD_StrategyProfitMargin) Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	return decimal.NewFromFloat(25), decimal.NewFromFloat(20), decimal.NewFromFloat(3.25)
}
