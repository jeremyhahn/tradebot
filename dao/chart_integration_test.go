//// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/stretchr/testify/assert"
)

func TestChartDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)

	daoUser := &entity.User{
		Id:            ctx.User.GetId(),
		Username:      ctx.User.GetUsername(),
		LocalCurrency: ctx.User.GetLocalCurrency()}

	chart := createIntegrationTestChart(ctx)
	indicators := chart.GetIndicators()
	trades := chart.GetTrades()

	err := chartDAO.Create(chart)
	assert.Equal(t, nil, err)

	persistedChart, chartErr := chartDAO.Get(chart.GetId())
	assert.Equal(t, nil, chartErr)
	assert.Equal(t, chart.GetUserId(), persistedChart.GetUserId())
	assert.Equal(t, chart.GetBase(), persistedChart.GetBase())
	assert.Equal(t, chart.GetQuote(), persistedChart.GetQuote())
	assert.Equal(t, chart.GetExchangeName(), persistedChart.GetExchangeName())
	assert.Equal(t, chart.GetPeriod(), persistedChart.GetPeriod())
	assert.Equal(t, chart.IsAutoTrade(), persistedChart.IsAutoTrade())
	assert.Equal(t, 3, len(chart.GetIndicators()))
	assert.Equal(t, 2, len(chart.GetTrades()))

	persistedIndicators, terr := chartDAO.GetIndicators(chart)
	assert.Equal(t, nil, terr)
	assert.Equal(t, 3, len(persistedIndicators))

	assert.Equal(t, uint(1), persistedIndicators[0].GetId())
	assert.Equal(t, chart.GetId(), persistedIndicators[0].GetChartId())
	assert.Equal(t, indicators[0].GetName(), persistedIndicators[0].GetName())
	assert.Equal(t, indicators[0].GetParameters(), persistedIndicators[0].GetParameters())

	assert.Equal(t, uint(2), persistedIndicators[1].Id)
	assert.Equal(t, chart.GetId(), persistedIndicators[1].ChartId)
	assert.Equal(t, indicators[1].Name, persistedIndicators[1].Name)
	assert.Equal(t, indicators[1].Parameters, persistedIndicators[1].Parameters)

	persistedTrades, terr := chartDAO.GetTrades(ctx.User)
	assert.Equal(t, nil, terr)
	assert.Equal(t, 2, len(persistedTrades))

	assert.Equal(t, uint(1), persistedTrades[0].GetId())
	assert.Equal(t, chart.GetId(), persistedTrades[0].GetChartId())
	assert.Equal(t, daoUser.Id, persistedTrades[0].GetUserId())
	assert.Equal(t, trades[0].GetBase(), persistedTrades[0].GetBase())
	assert.Equal(t, trades[0].GetQuote(), persistedTrades[0].GetQuote())
	assert.Equal(t, trades[0].GetExchangeName(), persistedTrades[0].GetExchangeName())
	assert.Equal(t, trades[0].GetDate().UTC().String(), persistedTrades[0].GetDate().UTC().String())
	assert.Equal(t, trades[0].GetType(), persistedTrades[0].GetType())
	assert.Equal(t, trades[0].GetAmount(), persistedTrades[0].GetAmount())
	assert.Equal(t, trades[0].GetPrice(), persistedTrades[0].GetPrice())
	assert.Equal(t, trades[0].GetChartData(), persistedTrades[0].GetChartData())

	assert.Equal(t, uint(2), persistedTrades[1].GetId())
	assert.Equal(t, chart.GetId(), persistedTrades[1].GetChartId())
	assert.Equal(t, daoUser.Id, persistedTrades[1].GetUserId())
	assert.Equal(t, trades[1].GetBase(), persistedTrades[1].GetBase())
	assert.Equal(t, trades[1].GetQuote(), persistedTrades[1].GetQuote())
	assert.Equal(t, trades[1].GetExchangeName(), persistedTrades[1].GetExchangeName())
	assert.Equal(t, trades[1].GetDate().UTC().String(), persistedTrades[1].GetDate().UTC().String())
	assert.Equal(t, trades[1].GetType(), persistedTrades[1].GetType())
	assert.Equal(t, trades[1].GetAmount(), persistedTrades[1].GetAmount())
	assert.Equal(t, trades[1].GetPrice(), persistedTrades[1].GetPrice())
	assert.Equal(t, trades[1].GetChartData(), persistedTrades[1].GetChartData())

	CleanupIntegrationTest()
}

func TestChartDAO_Find(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	err := chartDAO.Create(chart)
	assert.Nil(t, err)

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(charts))

	CleanupIntegrationTest()
}

func TestChartDAO_GetLastTrade(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	chartDAO.Create(chart)

	trade, err := chartDAO.GetLastTrade(chart)
	assert.Equal(t, nil, err)
	assert.Equal(t, uint(2), trade.GetId())
	assert.Equal(t, "sell", trade.GetType())

	CleanupIntegrationTest()
}
