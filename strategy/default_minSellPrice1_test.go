package strategy

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChartDAO_MinSellPrice1 struct {
	dao.ChartDAO
	mock.Mock
}

func TestDefaultTradingStrategy_minSellPrice_Default(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chartService := new(MockChartService)
	chart := new(MockChart)
	chartDAO := new(MockChartDAO_MinSellPrice1)
	chartDAO.On("GetLastTrade", chart).Return(chartDAO.GetLastTrade(chart)).Once()
	strategy := NewDefaultTradingStrategy(ctx, chart, chartDAO, new(MockProfitDAO))
	strategy.OnPriceChange(chartService)
	minPrice := strategy.minSellPrice(chartService.GetExchange().GetTradingFee())
	assert.Equal(t, 10000.0, strategy.lastTrade.Price)
	assert.Equal(t, 10000.0, chartService.GetData().Price)
	assert.Equal(t, 0.025, chartService.GetExchange().GetTradingFee()) // 10000 * 0.025 = 250
	assert.Equal(t, .10, strategy.config.profitMarginMinPercent)       // 1000
	assert.Equal(t, 0.0, strategy.config.profitMarginMin)
	assert.Equal(t, 0.4, strategy.config.tax)
	assert.Equal(t, 11675.0, minPrice)
	chartDAO.AssertExpectations(t)
}

func TestDefaultTradingStrategy_minSellPrice_NoTax(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chartService := new(MockChartService)
	chart := new(MockChart)
	chartDAO := new(MockChartDAO_MinSellPrice1)
	chartDAO.On("GetLastTrade", chart).Return(chartDAO.GetLastTrade(chart)).Once()
	strategy := CreateDefaultTradingStrategy(ctx, chart, chartDAO, new(MockProfitDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    0,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        0,
		profitMarginMinPercent: .10,
		tradeSize:              0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})
	strategy.OnPriceChange(chartService)
	minPrice := strategy.minSellPrice(chartService.GetExchange().GetTradingFee())
	assert.Equal(t, 10000.0, strategy.lastTrade.Price)
	assert.Equal(t, 10000.0, chartService.GetData().Price)
	assert.Equal(t, 0.025, chartService.GetExchange().GetTradingFee()) // 10000 * 0.025 = 250
	assert.Equal(t, .10, strategy.config.profitMarginMinPercent)       // 1000
	assert.Equal(t, 0.0, strategy.config.profitMarginMin)
	assert.Equal(t, 0.0, strategy.config.tax)
	assert.Equal(t, 11275.0, minPrice)
	chartDAO.AssertExpectations(t)
}

func (mdao *MockChartDAO_MinSellPrice1) GetLastTrade(chart dao.IChart) *dao.Trade {
	mdao.Called(chart)
	return &dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    10000,
		UserID:   1}
}

func (mdao *MockChartDAO_MinSellPrice1) Save(dao dao.IChart) {}
