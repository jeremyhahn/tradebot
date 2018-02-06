package service

import (
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/stretchr/testify/mock"
)

type MockChartDAO struct {
	dao.ChartDAO
	mock.Mock
}

type MockExchangeService struct {
	ExchangeService
	mock.Mock
}

type MockTradeService struct {
	TradeService
	mock.Mock
}

type MockProfitService struct {
	ProfitService
	mock.Mock
}

/*
func TestAutoTradeService(t *testing.T) {
	ctx := NewUnitTestContext()
	chartDAO := new(MockChartDAO)
	exchangeService := new(MockExchangeService)
	tradeService := new(MockTradeService)
	profitService := new(MockProfitService)
	autoTradeService := NewAutoTradeService(ctx, exchangeService, chartDAO, tradeService, profitService)
	autoTradeService.Trade()
}

func (mts *MockTradeService) GetLastTrade(chart *common.Chart) *common.Trade {
	return &common.Trade{}
}

func (mts *MockTradeService) Save(trade *common.Trade) {

}

func (mps *MockProfitService) Find() {

}

func (mcdao *MockChartDAO) Find(user *common.User) []dao.Chart {
	return []dao.Chart{
		dao.Chart{
			ID:         1,
			UserID:     1,
			Base:       "BTC",
			Quote:      "USD",
			Exchange:   "gdax",
			Period:     900,
			AutoTrade:  1,
			Indicators: nil,
			Trades:     nil}}
}

func (mcdao *MockChartDAO) GetIndicators(entity dao.IChart) map[string]dao.Indicator {
	return map[string]dao.Indicator{
		"Test": dao.Indicator{
			Id:         1,
			ChartID:    1,
			Name:       "TestIndicator",
			Parameters: "1,2,3"}}
}

func (mes *MockExchangeService) NewExchange(user *common.User, exchangeName string,
	currencyPair *common.CurrencyPair) common.Exchange {

	return nil
}
*/
