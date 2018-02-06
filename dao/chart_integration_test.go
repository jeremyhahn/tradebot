// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChartDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)

	daoUser := User{
		Id:            ctx.User.Id,
		Username:      ctx.User.Username,
		LocalCurrency: ctx.User.LocalCurrency}

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

	assert.Equal(t, uint(1), persistedIndicators[0].Id)
	assert.Equal(t, chart.GetId(), persistedIndicators[0].ChartId)
	assert.Equal(t, indicators[0].Name, persistedIndicators[0].Name)
	assert.Equal(t, indicators[0].Parameters, persistedIndicators[0].Parameters)

	assert.Equal(t, uint(2), persistedIndicators[1].Id)
	assert.Equal(t, chart.GetId(), persistedIndicators[1].ChartId)
	assert.Equal(t, indicators[1].Name, persistedIndicators[1].Name)
	assert.Equal(t, indicators[1].Parameters, persistedIndicators[1].Parameters)

	persistedTrades, terr := chartDAO.GetTrades(ctx.User)
	assert.Equal(t, nil, terr)
	assert.Equal(t, 2, len(persistedTrades))

	assert.Equal(t, uint(1), persistedTrades[0].GetId())
	assert.Equal(t, chart.GetId(), persistedTrades[0].ChartId)
	assert.Equal(t, daoUser.Id, persistedTrades[0].UserId)
	assert.Equal(t, trades[0].Base, persistedTrades[0].Base)
	assert.Equal(t, trades[0].Quote, persistedTrades[0].Quote)
	assert.Equal(t, trades[0].Exchange, persistedTrades[0].Exchange)
	assert.Equal(t, trades[0].Date.UTC().String(), persistedTrades[0].Date.UTC().String())
	assert.Equal(t, trades[0].Type, persistedTrades[0].Type)
	assert.Equal(t, trades[0].Amount, persistedTrades[0].Amount)
	assert.Equal(t, trades[0].Price, persistedTrades[0].Price)
	assert.Equal(t, trades[0].ChartData, persistedTrades[0].ChartData)

	assert.Equal(t, uint(2), persistedTrades[1].GetId())
	assert.Equal(t, chart.GetId(), persistedTrades[1].ChartId)
	assert.Equal(t, daoUser.Id, persistedTrades[1].UserId)
	assert.Equal(t, trades[1].Base, persistedTrades[1].Base)
	assert.Equal(t, trades[1].Quote, persistedTrades[1].Quote)
	assert.Equal(t, trades[1].Exchange, persistedTrades[1].Exchange)
	assert.Equal(t, trades[1].Date.UTC().String(), persistedTrades[1].Date.UTC().String())
	assert.Equal(t, trades[1].Type, persistedTrades[1].Type)
	assert.Equal(t, trades[1].Amount, persistedTrades[1].Amount)
	assert.Equal(t, trades[1].Price, persistedTrades[1].Price)
	assert.Equal(t, trades[1].ChartData, persistedTrades[1].ChartData)

	CleanupIntegrationTest()
}

func TestChartDAO_Find(t *testing.T) {
	ctx := NewIntegrationTestContext()
	chartDAO := NewChartDAO(ctx)

	chart := createIntegrationTestChart(ctx)
	chartDAO.Create(chart)

	chart2 := createIntegrationTestChart(ctx)
	chart2.Base = "ETH"
	chartDAO.Create(chart2)

	charts, err := chartDAO.Find(ctx.User)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(charts))

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
