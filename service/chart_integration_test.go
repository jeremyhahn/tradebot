// +build integration

package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExchange struct {
	common.Exchange
	mock.Mock
}

func TestChartDAO(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)

	chart := createChart(ctx)
	chartDAO.Create(chart)

	charts := chartDAO.Find(ctx.User)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())

	trades := charts[0].GetTrades()
	assert.Equal(t, 2, len(trades))
	assert.Equal(t, true, trades[0].Date.Second() > 0)
	assert.Equal(t, "buy", trades[0].Type)
	assert.Equal(t, "BTC", trades[0].Base)
	assert.Equal(t, "USD", trades[0].Quote)
	assert.Equal(t, "gdax", trades[0].Exchange)
	assert.Equal(t, 1.0, trades[0].Amount)
	assert.Equal(t, 10000.0, trades[0].Price)
	assert.Equal(t, uint(1), trades[0].UserID)

	indicators := charts[0].GetIndicators()
	assert.Equal(t, 3, len(indicators))

	test.CleanupMockContext()
}

func TestChartDAO_GetIndicators(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)

	chart := createChart(ctx)
	chart.Trades = nil
	chart.Indicators = nil
	chartDAO.Create(chart)

	charts := chartDAO.Find(ctx.User)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())
	assert.Equal(t, 0, len(charts[0].GetTrades()))
	assert.Equal(t, 0, len(charts[0].GetIndicators()))

	charts[0].SetIndicators(createChartIndicators())
	chartDAO.Save(&charts[0])

	indicators := charts[0].GetIndicators()
	assert.Equal(t, 3, len(indicators))

	test.CleanupMockContext()
}

func TestChartDAO_GetTrades(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)

	chart := createChart(ctx)
	chart.Trades = nil
	chart.Indicators = nil
	chartDAO.Create(chart)

	charts := chartDAO.Find(ctx.User)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())
	assert.Equal(t, 0, len(charts[0].GetTrades()))
	assert.Equal(t, 0, len(charts[0].GetIndicators()))

	charts[0].SetTrades(createChartTrades())
	chartDAO.Save(&charts[0])

	trades := charts[0].GetTrades()
	assert.Equal(t, 2, len(trades))
	assert.Equal(t, true, trades[0].Date.Second() > 0)
	assert.Equal(t, "buy", trades[0].Type)
	assert.Equal(t, "BTC", trades[0].Base)
	assert.Equal(t, "USD", trades[0].Quote)
	assert.Equal(t, "gdax", trades[0].Exchange)
	assert.Equal(t, 1.0, trades[0].Amount)
	assert.Equal(t, 10000.0, trades[0].Price)
	assert.Equal(t, uint(1), trades[0].UserID)

	lastTrade := chartDAO.GetLastTrade(chart)
	assert.Equal(t, "sell", lastTrade.Type)
	assert.Equal(t, "BTC", lastTrade.Base)
	assert.Equal(t, "USD", lastTrade.Quote)
	assert.Equal(t, "gdax", lastTrade.Exchange)
	assert.Equal(t, 1.0, lastTrade.Amount)
	assert.Equal(t, 12000.0, lastTrade.Price)

	test.CleanupMockContext()
}

func TestChartService_GetIndicators(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)
	chart := createChart(ctx)
	chartDAO.Create(chart)

	charts := chartDAO.Find(ctx.User)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())
	assert.Equal(t, 2, len(charts[0].GetTrades()))

	indicators := charts[0].GetIndicators()
	assert.Equal(t, 3, len(indicators))
	assert.Equal(t, "BBands", indicators[0].Name)
	assert.Equal(t, "20,2", indicators[0].Parameters)
	assert.Equal(t, "MACD", indicators[1].Name)
	assert.Equal(t, "12,26,9", indicators[1].Parameters)
	assert.Equal(t, "RSI", indicators[2].Name)
	assert.Equal(t, "14,70,30", indicators[2].Parameters)

	service := NewChartService(ctx, chartDAO, &charts[0], new(MockExchange))
	Indicators := service.GetIndicators()
	assert.Equal(t, 3, len(Indicators))

	test.CleanupMockContext()
}

func TestChartService_Stream(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)
	chart := createChart(ctx)
	chartDAO.Create(chart)

	charts := chartDAO.Find(ctx.User)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())

	var receivedService common.ChartService
	service := NewChartService(ctx, chartDAO, &charts[0], new(MockExchange))
	service.Stream(func(chart common.ChartService) {
		receivedService = chart
		service.StopStream()
	})
	indicators := receivedService.GetIndicators()
	assert.Equal(t, true, len(indicators) == 3)

	test.CleanupMockContext()
}

func createChartIndicators() []dao.Indicator {
	var indicators []dao.Indicator
	indicators = append(indicators, dao.Indicator{
		Name:       "RSI",
		Parameters: "14,70,30"})
	indicators = append(indicators, dao.Indicator{
		Name:       "BBands",
		Parameters: "20,2"})
	indicators = append(indicators, dao.Indicator{
		Name:       "MACD",
		Parameters: "12,26,9"})
	return indicators
}

func createChartTrades() []dao.Trade {
	var trades []dao.Trade
	trades = append(trades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    10000,
		UserID:   1})
	trades = append(trades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -10),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    12000,
		UserID:   1})
	return trades
}

func createChart(ctx *common.Context) *dao.Chart {
	return &dao.Chart{
		UserID:     ctx.User.Id,
		Base:       "BTC",
		Quote:      "USD",
		Exchange:   "gdax",
		Period:     900,
		AutoTrade:  1,
		Indicators: createChartIndicators(),
		Trades:     createChartTrades()}
}

func (mcs *MockExchange) GetName() string {
	return "Test"
}

func (mcs *MockExchange) FormattedCurrencyPair() string {
	return "BTC-USD"
}

func (mcs *MockExchange) GetPriceHistory(start, end time.Time, granularity int) []common.Candlestick {
	return createChartCandles()
}

func (mcs *MockExchange) SubscribeToLiveFeed(priceChange chan common.PriceChange) {
	fmt.Println("Subscribing to feed")
	priceChange <- common.PriceChange{
		CurrencyPair: &common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"},
		Exchange: "gdax",
		Price:    12345.0,
		Satoshis: 0.12345678}
	fmt.Println("Price change sent")
}

func (mcs *MockExchange) GetCurrencyPair() common.CurrencyPair {
	return common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
}
