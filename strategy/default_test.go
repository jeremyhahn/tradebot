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

type MockedExchange struct {
	mock.Mock
}

type MockedChartService struct {
	mock.Mock
}

func TestDefaultTradingStrategy_SignalCount(t *testing.T) {

	ctx := test.NewTestContext()

	trades := make([]dao.Trade, 0, 5)
	trades = append(trades, dao.Trade{
		Date:     time.Now().AddDate(0, -1, 0),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    15000,
		UserID:   ctx.User.Id})
	trades = append(trades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    16000,
		UserID:   ctx.User.Id})

	autoTradeCoin := &dao.AutoTradeCoin{
		UserID:   ctx.User.Id,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Period:   900,
		Trades:   trades}

	autoTradeDAO := dao.NewAutoTradeDAO(ctx)
	autoTradeDAO.Create(autoTradeCoin)

	var coins []common.Coin
	coins = append(coins, common.Coin{
		Currency:  "BTC",
		Available: 25.01020304})

	strategy := NewDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, dao.NewSignalLogDAO(ctx))

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

func createCandlesticks() []common.Candlestick {
	var candles []common.Candlestick
	candles = append(candles, common.Candlestick{Close: 100.00})
	candles = append(candles, common.Candlestick{Close: 200.00})
	candles = append(candles, common.Candlestick{Close: 300.00})
	candles = append(candles, common.Candlestick{Close: 400.00})
	candles = append(candles, common.Candlestick{Close: 500.00})
	candles = append(candles, common.Candlestick{Close: 600.00})
	candles = append(candles, common.Candlestick{Close: 700.00})
	candles = append(candles, common.Candlestick{Close: 800.00})
	candles = append(candles, common.Candlestick{Close: 900.00})
	candles = append(candles, common.Candlestick{Close: 1000.00})
	candles = append(candles, common.Candlestick{Close: 1100.00})
	candles = append(candles, common.Candlestick{Close: 1200.00})
	candles = append(candles, common.Candlestick{Close: 1300.00})
	candles = append(candles, common.Candlestick{Close: 1400.00})
	candles = append(candles, common.Candlestick{Close: 1500.00})
	candles = append(candles, common.Candlestick{Close: 1600.00})
	candles = append(candles, common.Candlestick{Close: 1700.00})
	candles = append(candles, common.Candlestick{Close: 1800.00})
	candles = append(candles, common.Candlestick{Close: 1900.00})
	candles = append(candles, common.Candlestick{Close: 2000.00})
	candles = append(candles, common.Candlestick{Close: 2100.00})
	candles = append(candles, common.Candlestick{Close: 2200.00})
	candles = append(candles, common.Candlestick{Close: 2300.00})
	candles = append(candles, common.Candlestick{Close: 2400.00})
	candles = append(candles, common.Candlestick{Close: 2500.00})
	candles = append(candles, common.Candlestick{Close: 2600.00})
	candles = append(candles, common.Candlestick{Close: 2700.00})
	candles = append(candles, common.Candlestick{Close: 2800.00})
	candles = append(candles, common.Candlestick{Close: 2900.00})
	candles = append(candles, common.Candlestick{Close: 3000.00})
	return candles
}

func (mcs *MockedChartService) GetData() *common.ChartData {
	return createChartData()
}

/*
  exchange := new(MockedExchange)
  exchange.On("GetName").Return("TestExchange")
  exchange.On("GetBalances").Return(coins, 25.01020304)
  exchange.On("GetTradeHistory").Return(createCandlesticks())
  exchange.On("FormattedCurrencyPair").Return("BTC-USD")

  chart := new(MockedChartService)
  chart.On("GetExchange").Return(exchange)
  chart.On("GetCurrencyPair").Return(currencyPair)
  chart.On("GetData").Return(createChartData())
*/

/*currencyPair := common.CurrencyPair{
Base:          "BTC",
Quote:         "USD",
LocalCurrency: "USD"}*/
