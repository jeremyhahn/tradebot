//// +build integration

package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExchange_Chart struct {
	common.Exchange
	mock.Mock
}

type MockExchangeService_Chart struct {
	common.Exchange
	mock.Mock
}

func TestChartDAO(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)

	chart := createChart(ctx)
	chartDAO.Create(chart)

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
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
	assert.Equal(t, uint(1), trades[0].UserId)

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

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
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

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
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
	assert.Equal(t, uint(1), trades[0].UserId)

	lastTrade, err := chartDAO.GetLastTrade(chart)
	assert.Equal(t, nil, err)
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

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())
	assert.Equal(t, 2, len(charts[0].GetTrades()))

	indicators := charts[0].GetIndicators()
	assert.Equal(t, 3, len(indicators))
	assert.Equal(t, "BollingerBands", indicators[0].Name)
	assert.Equal(t, "20,2", indicators[0].Parameters)
	assert.Equal(t, "MovingAverageConvergenceDivergence", indicators[1].Name)
	assert.Equal(t, "12,26,9", indicators[1].Parameters)
	assert.Equal(t, "RelativeStrengthIndex", indicators[2].Name)
	assert.Equal(t, "14,70,30", indicators[2].Parameters)

	mapper := mapper.NewChartMapper(ctx)

	service := NewChartService(ctx, chartDAO, new(MockExchangeService_Chart))

	commonChart := mapper.MapChartEntityToDto(&charts[0])
	Indicators, err := service.GetIndicators(&commonChart)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(Indicators))

	test.CleanupMockContext()
}

func TestChartService_Stream(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	chartDAO := dao.NewChartDAO(ctx)
	chart := createChart(ctx)
	chartDAO.Create(chart)

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "BTC", charts[0].GetBase())
	assert.Equal(t, "USD", charts[0].GetQuote())
	assert.Equal(t, "gdax", charts[0].GetExchangeName())
	assert.Equal(t, 900, charts[0].GetPeriod())
	assert.Equal(t, true, charts[0].IsAutoTrade())

	var indicators []common.Indicator
	service := NewChartService(ctx, chartDAO, new(MockExchangeService_Chart))
	mapper := mapper.NewChartMapper(ctx)
	commonChart := mapper.MapChartEntityToDto(&charts[0])
	service.Stream(&commonChart, func(newPrice float64) error {
		indicators = commonChart.Indicators
		service.StopStream(&commonChart)
		return nil
	})
	assert.Equal(t, 3, len(indicators))

	test.CleanupMockContext()
}

func createChartIndicators() []dao.Indicator {
	var indicators []dao.Indicator
	indicators = append(indicators, dao.Indicator{
		Name:       "RelativeStrengthIndex",
		Parameters: "14,70,30"})
	indicators = append(indicators, dao.Indicator{
		Name:       "BollingerBands",
		Parameters: "20,2"})
	indicators = append(indicators, dao.Indicator{
		Name:       "MovingAverageConvergenceDivergence",
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
		UserId:   1})
	trades = append(trades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -10),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    12000,
		UserId:   1})
	return trades
}

func createChart(ctx *common.Context) *dao.Chart {
	return &dao.Chart{
		UserId:     ctx.User.Id,
		Base:       "BTC",
		Quote:      "USD",
		Exchange:   "gdax",
		Period:     900,
		AutoTrade:  1,
		Indicators: createChartIndicators(),
		Trades:     createChartTrades()}
}

func (mcs *MockExchange_Chart) GetName() string {
	return "Test"
}

func (mcs *MockExchange_Chart) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	return "BTC-USD"
}

func (mcs *MockExchange_Chart) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) []common.Candlestick {

	return createIntegrationTestCandles()
}

func (mcs *MockExchange_Chart) SubscribeToLiveFeed(currencyPair *common.CurrencyPair, priceChange chan common.PriceChange) {
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

func (mes *MockExchangeService_Chart) CreateExchange(user *common.User, exchangeName string) common.Exchange {
	return new(MockExchange_Chart)
}

func (mes *MockExchangeService_Chart) GetExchange(user *common.User, exchangeName string) common.Exchange {
	return new(MockExchange_Chart)
}

func (mes *MockExchangeService_Chart) GetExchanges(user *common.User) []common.Exchange {
	ctx := &common.Context{
		User: &common.User{
			Id:            1,
			Username:      test.TEST_USERNAME,
			LocalCurrency: "USD"}}
	testExchange := &dao.UserCryptoExchange{
		Name:   "Test Exchange",
		URL:    "https://www.example.com",
		Key:    "ABC123",
		Secret: "$ecret!",
		Extra:  "Exchange specific data here"}
	return []common.Exchange{exchange.NewGDAX(ctx, testExchange)}
}