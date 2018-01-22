package strategy

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAutoTradeDAO_MinSellPrice1 struct {
	dao.IAutoTradeDAO
	mock.Mock
}

func TestDefaultTradingStrategy_minSellPrice_Default(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chart := new(MockChartService)
	autoTradeCoin := new(MockAutoTradeCoin)
	autoTradeDAO := new(MockAutoTradeDAO_MinSellPrice1)
	autoTradeDAO.On("GetLastTrade", autoTradeCoin).Return(autoTradeDAO.GetLastTrade(autoTradeCoin)).Once()
	strategy := NewDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO))
	strategy.OnPriceChange(chart)
	minPrice := strategy.minSellPrice(chart.GetData().Price, chart.GetExchange().GetTradingFee())
	assert.Equal(t, float64(10000), strategy.lastTrade.Price)
	assert.Equal(t, float64(10000), chart.GetData().Price)
	assert.Equal(t, 0.025, chart.GetExchange().GetTradingFee())  // 10000 + 1000 = 11000 * 0.025 = 275
	assert.Equal(t, .10, strategy.config.profitMarginMinPercent) // 1000
	assert.Equal(t, float64(0), strategy.config.profitMarginMin)
	assert.Equal(t, float64(0), strategy.config.tax)
	assert.Equal(t, float64(11275), minPrice)
	autoTradeDAO.AssertExpectations(t)
}

func (mdao *MockAutoTradeDAO_MinSellPrice1) GetLastTrade(autoTradeCoin dao.IAutoTradeCoin) *dao.Trade {
	mdao.Called(autoTradeCoin)
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

func (mdao *MockAutoTradeDAO_MinSellPrice1) Save(dao dao.IAutoTradeCoin) {}
