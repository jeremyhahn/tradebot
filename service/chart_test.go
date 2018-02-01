package service

import (
	"github.com/jeremyhahn/tradebot/common"
)

/*
type MockChartDAO struct {
	dao.ChartDAO
	mock.Mock
}

type MockChartEntity struct {
	dao.IChart
	mock.Mock
}

type MockExchangeService struct {
	common.Exchange
	mock.Mock
}

func TestChartService(t *testing.T) {
	ctx := test.NewUnitTestContext()
	chartDAO := new(MockChartDAO)
	entity := new(MockChartEntity)
	exchangeService := new(MockExchangeService)
	chartService := NewChartService(ctx, chartDAO, entity, exchangeService)
	chartService.Stream(func(chart common.ChartService) {

	})
}

func (mcdao *MockChartDAO) GetIndicators(chart dao.IChart) map[string]dao.Indicator {
	return map[string]dao.Indicator{
		"TestIndicator": dao.Indicator{
			Id:         1,
			ChartID:    1,
			Name:       "TestIndicator",
			Parameters: "1,2,3"}}
}

func (mce *MockChartEntity) GetPeriod() int {
	return 14
}

func (eces *MockExchangeService) GetName() string {
	return "TestExchange"
}

func (eces *MockExchangeService) FormattedCurrencyPair() string {
	return "BTC-USD"
}

func (eces *MockExchangeService) GetPriceHistory(time.Time, time.Time, int) []common.Candlestick {
	return createChartCandles()
}
*/
func createChartCandles() []common.Candlestick {
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
	candles = append(candles, common.Candlestick{Close: 3200.00})
	candles = append(candles, common.Candlestick{Close: 3300.00})
	candles = append(candles, common.Candlestick{Close: 3400.00})
	candles = append(candles, common.Candlestick{Close: 3500.00})
	candles = append(candles, common.Candlestick{Close: 3600.00})
	candles = append(candles, common.Candlestick{Close: 3700.00})
	candles = append(candles, common.Candlestick{Close: 3800.00})
	candles = append(candles, common.Candlestick{Close: 3900.00})
	candles = append(candles, common.Candlestick{Close: 4000.00})
	return candles
}
