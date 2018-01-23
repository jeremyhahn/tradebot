package strategy

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChartService_Trade4 struct {
	common.ChartService
	mock.Mock
}

type MockChartService_Trade3 struct {
	common.ChartService
	mock.Mock
}

type MockChartService_Trade2 struct {
	common.ChartService
	mock.Mock
}

type MockAutoTradeDAO_Trade struct {
	dao.IAutoTradeDAO
	mock.Mock
}

func TestDefaultTradingStrategy_Trade(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	autoTradeCoin := createAutoTradeCoin()
	autoTradeDAO := dao.NewAutoTradeDAO(ctx)

	strategy := CreateDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, new(MockSignalLogDAO), &DefaultTradingStrategyConfig{
		rsiOverSold:            30,
		rsiOverBought:          70,
		tax:                    0,
		stopLoss:               0,
		stopLossPercent:        .20,
		profitMarginMin:        1000,
		profitMarginMinPercent: 0,
		tradeSizePercent:       0,
		requiredBuySignals:     2,
		requiredSellSignals:    2})

	chart := new(MockChartService)
	strategy.OnPriceChange(chart)
	coins, _ := chart.GetExchange().GetBalances()
	trade := autoTradeDAO.GetLastTrade(autoTradeCoin)
	assert.Equal(t, "buy", trade.Type)
	assert.Equal(t, coins[1].Available, trade.Amount)

	chartTrade2 := new(MockChartService_Trade2)
	strategy.OnPriceChange(chartTrade2)
	coins, _ = chartTrade2.GetExchange().GetBalances()
	trade = autoTradeDAO.GetLastTrade(autoTradeCoin)
	assert.Equal(t, "sell", trade.Type)
	assert.Equal(t, coins[1].Available, trade.Amount)

	chartTrade3 := new(MockChartService_Trade3)
	strategy.OnPriceChange(chartTrade3)
	coins, _ = chartTrade3.GetExchange().GetBalances()
	trade = autoTradeDAO.GetLastTrade(autoTradeCoin)
	assert.Equal(t, "sell", trade.Type)
	assert.Equal(t, 2, len(autoTradeDAO.GetTrades(ctx.User)))

	chartTrade4 := new(MockChartService_Trade4)
	strategy.OnPriceChange(chartTrade4)
	coins, _ = chartTrade4.GetExchange().GetBalances()
	trade = autoTradeDAO.GetLastTrade(autoTradeCoin)
	assert.Equal(t, "buy", trade.Type)
	assert.Equal(t, 3, len(autoTradeDAO.GetTrades(ctx.User)))

	test.CleanupMockContext()
}

func (mcs *MockChartService_Trade2) GetData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               18500,
		RSILive:             85,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func (mcs *MockChartService_Trade2) GetCurrencyPair() common.CurrencyPair {
	return createCurrencyPair()
}

func (cs *MockChartService_Trade2) GetExchange() common.Exchange {
	return new(MockExchange)
}

func createAutoTradeCoin() *dao.AutoTradeCoin {
	sampleTrades := make([]dao.Trade, 0, 0)
	return &dao.AutoTradeCoin{
		UserID:   1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Period:   900,
		Trades:   sampleTrades}
}

func (mcs *MockChartService_Trade3) GetData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               10000,
		RSILive:             85,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func (mcs *MockChartService_Trade3) GetCurrencyPair() common.CurrencyPair {
	return createCurrencyPair()
}

func (cs *MockChartService_Trade3) GetExchange() common.Exchange {
	return new(MockExchange)
}

func (mcs *MockChartService_Trade4) GetData() *common.ChartData {
	return &common.ChartData{
		CurrencyPair:        common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
		Exchange:            "Test",
		Price:               8000,
		RSILive:             21,
		BollingerUpperLive:  15000,
		BollingerMiddleLive: 13000,
		BollingerLowerLive:  11000}
}

func (mcs *MockChartService_Trade4) GetCurrencyPair() common.CurrencyPair {
	return createCurrencyPair()
}

func (cs *MockChartService_Trade4) GetExchange() common.Exchange {
	return new(MockExchange)
}
