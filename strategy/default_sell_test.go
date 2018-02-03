// +build unit

package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
	"github.com/jeremyhahn/tradebot/test"
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
	lastTrade := &common.Trade{
		ID:       1,
		ChartID:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "buy",
		Amount:   1,
		Price:    8000}
	params := &TradingStrategyParams{
		CurrencyPair: &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Balances:     helper.CreateBalances(),
		Indicators:   strategyIndicators,
		NewPrice:     16000,
		LastTrade:    lastTrade,
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
	assert.Equal(t, sell, true)
	assert.Equal(t, err, nil)
}

func (mrsi *MockRSI_StrategySell) Calculate(price float64) float64 {
	return 71.0
}

func (mrsi *MockRSI_StrategySell) IsOverBought(price float64) bool {
	return true
}

func (mrsi *MockRSI_StrategySell) IsOverSold(price float64) bool {
	return false
}

func (mrsi *MockBBands_StrategySell) Calculate(price float64) (float64, float64, float64) {
	return 15000.0, 12500.0, 10000.0
}
