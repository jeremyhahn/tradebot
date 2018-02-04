package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/indicators"
	"github.com/jeremyhahn/tradebot/test"
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
	lastTrade := &common.Trade{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    10000}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		TradeSize:              1,
		ProfitMarginMin:        10000,
		ProfitMarginMinPercent: 0,
		StopLoss:               0,
		StopLossPercent:        .20,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   indicators,
		LastTrade:    lastTrade,
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)
	strategy := s.(*DefaultTradingStrategy)

	_, sell, _ := strategy.GetBuySellSignals()
	assert.Equal(t, sell, true)

	minSellPrice := strategy.minSellPrice()
	assert.Equal(t, minSellPrice, 24500.0)

	fees, tax := strategy.CalculateFeeAndTax(minSellPrice)
	assert.Equal(t, fees, 612.5)
	assert.Equal(t, tax, 5800.0)
	assert.Equal(t, "Aborting sale. Doesn't meet minimum trade requirements. price=11000.000000, minRequired=24500.000000", strategy.sell().Error())
}

func TestDefaultTradingStrategy_CustomTradeSize_ProfitMarginFixed2(t *testing.T) {
	helper := &test.StrategyTestHelper{}
	indicators := map[string]common.FinancialIndicator{
		"RelativeStrengthIndex":              new(MockRSI_StrategyProfitMargin),
		"BollingerBands":                     new(MockBBands_StrategyProfitMargin),
		"MovingAverageConvergenceDivergence": new(MockMACD_StrategyProfitMargin)}
	lastTrade := &common.Trade{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    10000}
	config := &DefaultTradingStrategyConfig{
		Tax:                    .40,
		StopLoss:               0,
		StopLossPercent:        .20,
		ProfitMarginMin:        200,
		ProfitMarginMinPercent: 0,
		TradeSize:              1,
		RequiredBuySignals:     2,
		RequiredSellSignals:    2}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		NewPrice:     11000,
		Indicators:   indicators,
		LastTrade:    lastTrade,
		TradeFee:     .025,
		Config:       config.ToSlice()}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	_, sell, _ := strategy.GetBuySellSignals()
	assert.Equal(t, sell, true)

	fees, tax := strategy.CalculateFeeAndTax(params.NewPrice)
	assert.Equal(t, strategy.minSellPrice(), 10535.0)
	assert.Equal(t, fees, 275.0)
	assert.Equal(t, tax, 400.0)
	assert.Equal(t, nil, strategy.sell())
}

func (mrsi *MockRSI_StrategyProfitMargin) Calculate(price float64) float64 {
	return 79.0
}

func (mrsi *MockRSI_StrategyProfitMargin) IsOverBought(rsiValue float64) bool {
	return true
}

func (mrsi *MockRSI_StrategyProfitMargin) IsOverSold(rsiValue float64) bool {
	return false
}

func (mrsi *MockBBands_StrategyProfitMargin) Calculate(price float64) (float64, float64, float64) {
	return 10000.0, 9000.0, 8000.0
}
