package mapper

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestChartMapper(t *testing.T) {
	ctx := test.NewUnitTestContext()
	mapper := NewChartMapper(ctx)

	chartIndicatorDTOs := []common.ChartIndicator{
		&dto.ChartIndicatorDTO{
			ChartId:    1,
			Name:       "TestIndicator1",
			Filename:   "test_indicator1.so",
			Parameters: "1,2,3"},
		&dto.ChartIndicatorDTO{
			ChartId:    1,
			Name:       "TestIndicator2",
			Filename:   "test_indicator2.so",
			Parameters: "4,5,6"}}

	chartStrategyDTOs := []common.ChartStrategy{
		&dto.ChartStrategyDTO{
			ChartId:    1,
			Name:       "TestStrategy",
			Filename:   "test_strategy.so",
			Parameters: "1,2,3"}}

	chartTradeDTOs := []common.Trade{
		&dto.TradeDTO{
			Id:       1,
			ChartId:  1,
			UserId:   1,
			Base:     "BTC",
			Quote:    "USD",
			Exchange: "gdax",
			Date:     time.Now().AddDate(0, 0, -1),
			Type:     "buy",
			Price:    5000,
			Amount:   1.5},
		&dto.TradeDTO{
			Id:       2,
			ChartId:  1,
			UserId:   1,
			Base:     "BTC",
			Quote:    "USD",
			Exchange: "gdax",
			Date:     time.Now(),
			Type:     "sell",
			Price:    10000,
			Amount:   1.0}}

	chartDTO := &dto.ChartDTO{
		Base:       "BTC",
		Quote:      "USD",
		Exchange:   "TestExchange",
		Period:     900,
		Price:      12000,
		AutoTrade:  1,
		Indicators: chartIndicatorDTOs,
		Strategies: chartStrategyDTOs,
		Trades:     chartTradeDTOs}

	chartEntity := mapper.MapChartDtoToEntity(chartDTO)
	assert.NotNil(t, chartEntity)
	assert.Equal(t, 2, len(chartEntity.GetIndicators()))
	assert.Equal(t, 1, len(chartEntity.GetStrategies()))
	assert.Equal(t, 2, len(chartEntity.GetTrades()))

	assert.Equal(t, chartEntity.GetId(), chartDTO.GetId())
	assert.Equal(t, chartEntity.GetBase(), chartDTO.GetBase())
	assert.Equal(t, chartEntity.GetQuote(), chartDTO.GetQuote())
	assert.Equal(t, chartEntity.GetExchangeName(), chartDTO.GetExchange())
	assert.Equal(t, chartEntity.GetPeriod(), chartDTO.GetPeriod())
	assert.Equal(t, chartEntity.GetAutoTrade(), chartDTO.GetAutoTrade())

	chartIndicatorEntities := chartEntity.GetIndicators()
	assert.Equal(t, chartIndicatorEntities[0].GetId(), chartDTO.GetIndicators()[0].GetId())
	assert.Equal(t, chartIndicatorEntities[0].GetName(), chartDTO.GetIndicators()[0].GetName())
	assert.Equal(t, chartIndicatorEntities[0].GetChartId(), chartDTO.GetIndicators()[0].GetChartId())
	assert.Equal(t, chartIndicatorEntities[0].GetName(), chartDTO.GetIndicators()[0].GetName())
	assert.Equal(t, chartIndicatorEntities[0].GetParameters(), chartDTO.GetIndicators()[0].GetParameters())

	chartStrategyEntities := chartEntity.GetStrategies()
	assert.Equal(t, chartStrategyEntities[0].GetId(), chartDTO.GetStrategies()[0].GetId())
	assert.Equal(t, chartStrategyEntities[0].GetName(), chartDTO.GetStrategies()[0].GetName())
	assert.Equal(t, chartStrategyEntities[0].GetChartId(), chartDTO.GetStrategies()[0].GetChartId())
	assert.Equal(t, chartStrategyEntities[0].GetName(), chartDTO.GetStrategies()[0].GetName())
	assert.Equal(t, chartStrategyEntities[0].GetParameters(), chartDTO.GetStrategies()[0].GetParameters())

	mappedDTO := mapper.MapChartEntityToDto(chartEntity)
	assert.Equal(t, chartEntity.GetId(), mappedDTO.GetId())
	assert.Equal(t, chartEntity.GetBase(), mappedDTO.GetBase())
	assert.Equal(t, chartEntity.GetQuote(), mappedDTO.GetQuote())
	assert.Equal(t, chartEntity.GetExchangeName(), mappedDTO.GetExchange())
	assert.Equal(t, chartEntity.GetPeriod(), mappedDTO.GetPeriod())
	assert.Equal(t, chartEntity.GetAutoTrade(), mappedDTO.GetAutoTrade())
	assert.Equal(t, chartEntity.IsAutoTrade(), mappedDTO.IsAutoTrade())
	assert.Equal(t, chartEntity.GetIndicators()[0].GetId(), mappedDTO.GetIndicators()[0].GetId())
	assert.Equal(t, chartEntity.GetIndicators()[0].GetChartId(), mappedDTO.GetIndicators()[0].GetChartId())
	assert.Equal(t, chartEntity.GetIndicators()[0].GetName(), mappedDTO.GetIndicators()[0].GetName())
	assert.Equal(t, chartEntity.GetIndicators()[0].GetParameters(), mappedDTO.GetIndicators()[0].GetParameters())

	assert.NotNil(t, mappedDTO.GetStrategies()[0])

	assert.Equal(t, chartEntity.GetStrategies()[0].GetId(), mappedDTO.GetStrategies()[0].GetId())
	assert.Equal(t, chartEntity.GetStrategies()[0].GetName(), mappedDTO.GetStrategies()[0].GetName())
	assert.Equal(t, chartEntity.GetStrategies()[0].GetChartId(), mappedDTO.GetStrategies()[0].GetChartId())
	assert.Equal(t, chartEntity.GetStrategies()[0].GetName(), mappedDTO.GetStrategies()[0].GetName())
	assert.Equal(t, chartEntity.GetStrategies()[0].GetParameters(), mappedDTO.GetStrategies()[0].GetParameters())
}
