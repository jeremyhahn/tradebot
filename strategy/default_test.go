package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChartService struct {
	common.ChartService
	mock.Mock
}

type MockExchange struct {
	common.Exchange
	mock.Mock
}

type MockProfitDAO struct {
	dao.ProfitDAO
	mock.Mock
}

type MockChartDAO struct {
	dao.ChartDAO
	mock.Mock
}

type MockChart struct {
	dao.IChart
	mock.Mock
}

func TestDefaultTradingStrategy_SignalCount(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChart)
	chartDAO := new(MockChartDAO)
	strategy := NewDefaultTradingStrategy(ctx, chart, chartDAO, new(MockProfitDAO))
	buySignals, sellSignals := strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               13000,
		RSILive:             50,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 0)
	assert.Equal(t, sellSignals, 0)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               12000,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 1)
	assert.Equal(t, sellSignals, 0)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               10000,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 2)
	assert.Equal(t, sellSignals, 0)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               15001,
		RSILive:             50,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 0)
	assert.Equal(t, sellSignals, 1)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               15001,
		RSILive:             80,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 0)
	assert.Equal(t, sellSignals, 2)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               15001,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 1)
	assert.Equal(t, sellSignals, 1)

	buySignals, sellSignals = strategy.countSignals(&common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               10001,
		RSILive:             80,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000})
	assert.Equal(t, buySignals, 1)
	assert.Equal(t, sellSignals, 1)
}

func TestDefaultTradingStrategy_getTradeAmounts_WithoutTradeSizePercent(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chartService := new(MockChartService)
	chart := new(MockChart)
	chartDAO := new(MockChartDAO)
	strategy := NewDefaultTradingStrategy(ctx, chart, chartDAO, new(MockProfitDAO))
	buyAmount, quoteAmount := strategy.getTradeAmounts(chartService)
	assert.Equal(t, 1.0, buyAmount)
	assert.Equal(t, 50.25, quoteAmount)
}

func TestDefaultTradingStrategy_getTradeAmounts_WithTradeSizePercent(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chartService := new(MockChartService)
	chart := new(MockChart)
	chartDAO := new(MockChartDAO)
	strategy := CreateDefaultTradingStrategy(ctx, chart, chartDAO, new(MockProfitDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    0,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        0,
		profitMarginMinPercent: .10,
		tradeSize:              .10,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	buyAmount, quoteAmount := strategy.getTradeAmounts(chartService)
	assert.Equal(t, 0.10, buyAmount)
	assert.Equal(t, 5.025, quoteAmount)
}

// -------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------

func createChartData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               10000,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func createCurrencyPair() common.CurrencyPair {
	return common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
}

func (mcs *MockChartService) GetData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               10000,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func (mcs *MockChartService) GetCurrencyPair() common.CurrencyPair {
	return createCurrencyPair()
}

func (cs *MockChartService) GetExchange() common.Exchange {
	return new(MockExchange)
}

func (mcs *MockExchange) GetBalances() ([]common.Coin, float64) {
	btc := 1.0
	usd := 50.25
	ltc := 75.50
	var coins []common.Coin
	coins = append(coins, common.Coin{
		Currency:  "USD",
		Available: usd})
	coins = append(coins, common.Coin{
		Currency:  "BTC",
		Available: btc})
	coins = append(coins, common.Coin{
		Currency:  "LTC",
		Available: ltc})
	return coins, btc + usd + ltc
}

func (mcs *MockExchange) GetTradingFee() float64 {
	return .025
}

func (mdao *MockChartDAO) GetLastTrade(chart dao.IChart) *dao.Trade {
	mdao.Called(chart)
	trades := chart.GetTrades()
	return &trades[len(trades)-1]
}

func (mdao *MockChartDAO) Save(dao dao.IChart)    {}
func (mdao *MockProfitDAO) Save(dao *dao.Profit)  {}
func (mdao *MockChart) AddTrade(trade *dao.Trade) {}
