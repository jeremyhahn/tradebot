package strategy

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChartService_MinSellPrice2 struct {
	common.ChartService
	mock.Mock
}

type MockAutoTradeDAO_MinSellPrice2 struct {
	dao.IAutoTradeDAO
	mock.Mock
}

type MockAutoTradeDAO_MinSellPrice2_2 struct {
	dao.IAutoTradeDAO
	mock.Mock
}

func TestDefaultTradingStrategy_minSellPrice_LastPriceGreater(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChartService)
	autoTradeCoin := new(MockAutoTradeCoin)
	autoTradeDAO := new(MockAutoTradeDAO_MinSellPrice2)
	autoTradeDAO.On("GetLastTrade", autoTradeCoin).Return(autoTradeDAO.GetLastTrade(autoTradeCoin)).Once()
	strategy := CreateDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    0,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        0,
		profitMarginMinPercent: .10,
		tradeSizePercent:       0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	strategy.OnPriceChange(chart)
	minPrice := strategy.minSellPrice(chart.GetData().Price, chart.GetExchange().GetTradingFee())
	assert.Equal(t, float64(18000), strategy.lastTrade.Price)
	assert.Equal(t, float64(10000), chart.GetData().Price)
	assert.Equal(t, 0.025, chart.GetExchange().GetTradingFee())  // 18000 + 1800 = 19800 * .025 = 495
	assert.Equal(t, .10, strategy.config.profitMarginMinPercent) // 18000 * .10 = 1800
	assert.Equal(t, float64(0), strategy.config.profitMarginMin)
	assert.Equal(t, float64(0), strategy.config.tax)
	assert.Equal(t, float64(20295), minPrice)
	autoTradeDAO.AssertExpectations(t)
}

func TestDefaultTradingStrategy_minSellPrice_NoProfitPercent(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChartService)
	autoTradeCoin := new(MockAutoTradeCoin)
	autoTradeDAO := new(MockAutoTradeDAO_MinSellPrice2)
	autoTradeDAO.On("GetLastTrade", autoTradeCoin).Return(autoTradeDAO.GetLastTrade(autoTradeCoin)).Once()
	strategy := CreateDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    0,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        500,
		profitMarginMinPercent: 0,
		tradeSizePercent:       0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	strategy.OnPriceChange(chart)
	minPrice := strategy.minSellPrice(chart.GetData().Price, chart.GetExchange().GetTradingFee())
	assert.Equal(t, float64(18000), strategy.lastTrade.Price)
	assert.Equal(t, float64(10000), chart.GetData().Price)
	assert.Equal(t, 0.025, chart.GetExchange().GetTradingFee()) // 18000 + 500 = 18500 * 0.025 = 462.5
	assert.Equal(t, float64(0), strategy.config.profitMarginMinPercent)
	assert.Equal(t, float64(500), strategy.config.profitMarginMin)
	assert.Equal(t, float64(0), strategy.config.tax)
	assert.Equal(t, float64(18962.5), minPrice)
	autoTradeDAO.AssertExpectations(t)
}

func TestDefaultTradingStrategy_minSellPrice_DoesntIncludeTax(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChartService)
	autoTradeCoin := new(MockAutoTradeCoin)
	autoTradeDAO := new(MockAutoTradeDAO_MinSellPrice2)
	autoTradeDAO.On("GetLastTrade", autoTradeCoin).Return(autoTradeDAO.GetLastTrade(autoTradeCoin)).Once()
	strategy := CreateDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    .20,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        500,
		profitMarginMinPercent: 0,
		tradeSizePercent:       0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	strategy.OnPriceChange(chart)
	minPrice := strategy.minSellPrice(chart.GetData().Price, chart.GetExchange().GetTradingFee())
	assert.Equal(t, float64(18000), strategy.lastTrade.Price)
	assert.Equal(t, float64(10000), chart.GetData().Price)
	assert.Equal(t, 0.025, chart.GetExchange().GetTradingFee()) // 18000 + 500 = 18500 * .025 = 462.5
	assert.Equal(t, float64(0), strategy.config.profitMarginMinPercent)
	assert.Equal(t, float64(500), strategy.config.profitMarginMin)
	assert.Equal(t, .20, strategy.config.tax) // 500 * .20 = 100
	assert.Equal(t, 19062.5, minPrice)
	autoTradeDAO.AssertExpectations(t)
}

func TestDefaultTradingStrategy_minSellPrice_IncludesTax(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChartService)
	autoTradeCoin := new(MockAutoTradeCoin)
	autoTradeDAO := new(MockAutoTradeDAO_MinSellPrice2_2)
	autoTradeDAO.On("GetLastTrade", autoTradeCoin).Return(autoTradeDAO.GetLastTrade(autoTradeCoin)).Once()
	strategy := CreateDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    .20,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        500,
		profitMarginMinPercent: 0,
		tradeSizePercent:       0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	strategy.OnPriceChange(chart)
	minPrice := strategy.minSellPrice(chart.GetData().Price, chart.GetExchange().GetTradingFee())
	assert.Equal(t, float64(9000), strategy.lastTrade.Price)
	assert.Equal(t, float64(10000), chart.GetData().Price)
	assert.Equal(t, 0.025, chart.GetExchange().GetTradingFee()) // 10000 * .025 = 250
	assert.Equal(t, float64(0), strategy.config.profitMarginMinPercent)
	assert.Equal(t, float64(500), strategy.config.profitMarginMin)
	assert.Equal(t, .20, strategy.config.tax) // 1000 * .20 = 200
	assert.Equal(t, float64(10450), minPrice)
	autoTradeDAO.AssertExpectations(t)
}

func (mdao *MockAutoTradeDAO_MinSellPrice2) GetLastTrade(autoTradeCoin dao.IAutoTradeCoin) *dao.Trade {
	mdao.Called(autoTradeCoin)
	return &dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    18000,
		UserID:   1}
}

func (mdao *MockAutoTradeDAO_MinSellPrice2) Save(dao dao.IAutoTradeCoin) {}

func (mcs *MockChartService_MinSellPrice2) GetData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               18000,
		RSILive:             29,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func (mcs *MockChartService_MinSellPrice2) GetCurrencyPair() common.CurrencyPair {
	return createCurrencyPair()
}

func (cs *MockChartService_MinSellPrice2) GetExchange() common.Exchange {
	return new(MockExchange)
}

func (mdao *MockAutoTradeDAO_MinSellPrice2_2) GetLastTrade(autoTradeCoin dao.IAutoTradeCoin) *dao.Trade {
	mdao.Called(autoTradeCoin)
	return &dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    9000,
		UserID:   1}
}

func (mdao *MockAutoTradeDAO_MinSellPrice2_2) Save(dao dao.IAutoTradeCoin) {}
