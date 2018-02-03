// +build unit

package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
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
		ID:       1,
		ChartID:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    10000}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     nil,
		Indicators:   indicators,
		NewPrice:     13000,
		LastTrade:    lastTrade,
		TradeFee:     .025}

	s, err := CreateDefaultTradingStrategy(params)
	assert.Equal(t, nil, err)

	strategy := s.(*DefaultTradingStrategy)

	buy, sell, err := strategy.GetBuySellSignals()
	assert.Equal(t, false, buy)
	assert.Equal(t, true, sell)
	assert.Equal(t, nil, err)

	minPrice := strategy.minSellPrice()
	assert.Equal(t, 11675.0, minPrice)

	fees, tax := strategy.CalculateFeeAndTax(params.NewPrice)
	assert.Equal(t, 325.0, fees)
	assert.Equal(t, 1200.0, tax)
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

func (mrsi *MockBBands_StrategyCalculate) Calculate(price float64) (float64, float64, float64) {
	return 10000.0, 9000.0, 8000.0
}
