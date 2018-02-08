package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type ChartMapper interface {
	MapChartDtoToEntity(dto common.Chart) dao.ChartEntity
	MapChartEntityToDto(entity dao.ChartEntity) common.Chart
	MapIndicatorEntityToDto(entity dao.ChartIndicator) common.ChartIndicator
	MapIndicatorDtoToEntity(dto common.ChartIndicator) dao.ChartIndicator
	MapStrategyEntityToDto(entity dao.ChartStrategy) common.ChartStrategy
	MapStrategyDtoToEntity(dto common.ChartStrategy) dao.ChartStrategy
	MapTradeEntityToDto(entity *dao.Trade) common.Trade
	MapTradeDtoToEntity(trade common.Trade) dao.Trade
}

type DefaultChartMapper struct {
	ctx *common.Context
}

func NewChartMapper(ctx *common.Context) ChartMapper {
	return &DefaultChartMapper{ctx: ctx}
}

func (mapper *DefaultChartMapper) MapTradeEntityToDto(entity *dao.Trade) common.Trade {
	return &dto.TradeDTO{
		Id:        entity.GetId(),
		UserId:    entity.GetUserId(),
		ChartId:   entity.GetChartId(),
		Base:      entity.GetBase(),
		Quote:     entity.GetQuote(),
		Exchange:  entity.GetExchangeName(),
		Date:      entity.GetDate(),
		Type:      entity.GetType(),
		Amount:    entity.GetAmount(),
		Price:     entity.GetPrice(),
		ChartData: entity.GetChartData()}
}

func (mapper *DefaultChartMapper) MapTradeDtoToEntity(trade common.Trade) dao.Trade {
	return dao.Trade{
		Id:        trade.GetId(),
		UserId:    mapper.ctx.User.Id,
		ChartId:   trade.GetChartId(),
		Date:      trade.GetDate(),
		Exchange:  trade.GetExchange(),
		Type:      trade.GetType(),
		Base:      trade.GetBase(),
		Quote:     trade.GetQuote(),
		Amount:    trade.GetAmount(),
		Price:     trade.GetPrice(),
		ChartData: trade.GetChartData()}
}

func (mapper *DefaultChartMapper) MapIndicatorEntityToDto(entity dao.ChartIndicator) common.ChartIndicator {
	return &dto.ChartIndicatorDTO{
		Id:         entity.Id,
		ChartId:    entity.ChartId,
		Name:       entity.Name,
		Parameters: entity.Parameters}
}

func (mapper *DefaultChartMapper) MapIndicatorDtoToEntity(dto common.ChartIndicator) dao.ChartIndicator {
	return dao.ChartIndicator{
		Id:         dto.GetId(),
		ChartId:    dto.GetChartId(),
		Name:       dto.GetName(),
		Parameters: dto.GetParameters()}
}

func (mapper *DefaultChartMapper) MapStrategyEntityToDto(entity dao.ChartStrategy) common.ChartStrategy {
	return &dto.ChartStrategyDTO{
		Id:         entity.GetId(),
		ChartId:    entity.GetChartId(),
		Name:       entity.GetName(),
		Parameters: entity.GetParameters()}
}

func (mapper *DefaultChartMapper) MapStrategyDtoToEntity(dto common.ChartStrategy) dao.ChartStrategy {
	return dao.ChartStrategy{
		Id:         dto.GetId(),
		ChartId:    dto.GetChartId(),
		Name:       dto.GetName(),
		Parameters: dto.GetParameters()}
}

func (mapper *DefaultChartMapper) MapChartDtoToEntity(dto common.Chart) dao.ChartEntity {
	var daoChartIndicators []dao.ChartIndicator
	for _, indicator := range dto.GetIndicators() {
		daoChartIndicators = append(daoChartIndicators, mapper.MapIndicatorDtoToEntity(indicator))
	}
	var daoChartStrategies []dao.ChartStrategy
	for _, strategy := range dto.GetStrategies() {
		daoChartStrategies = append(daoChartStrategies, mapper.MapStrategyDtoToEntity(strategy))
	}
	var daoTrades []dao.Trade
	for _, trade := range dto.GetTrades() {
		daoTrades = append(daoTrades, mapper.MapTradeDtoToEntity(trade))
	}
	return &dao.Chart{
		Id:         dto.GetId(),
		UserId:     mapper.ctx.User.Id,
		Base:       dto.GetBase(),
		Quote:      dto.GetQuote(),
		Exchange:   dto.GetExchange(),
		Period:     dto.GetPeriod(),
		AutoTrade:  dto.GetAutoTrade(),
		Indicators: daoChartIndicators,
		Strategies: daoChartStrategies,
		Trades:     daoTrades}
}

func (mapper *DefaultChartMapper) MapChartEntityToDto(entity dao.ChartEntity) common.Chart {
	var indicators []common.ChartIndicator
	for _, indicator := range entity.GetIndicators() {
		indicators = append(indicators, mapper.MapIndicatorEntityToDto(indicator))
	}
	var strategies []common.ChartStrategy
	for _, strategy := range entity.GetStrategies() {
		strategies = append(strategies, mapper.MapStrategyEntityToDto(strategy))
	}
	var trades []common.Trade
	for _, trade := range entity.GetTrades() {
		trades = append(trades, mapper.MapTradeEntityToDto(&trade))
	}
	return dto.ChartDTO{
		Id:         entity.GetId(),
		Base:       entity.GetBase(),
		Quote:      entity.GetQuote(),
		Exchange:   entity.GetExchangeName(),
		Period:     entity.GetPeriod(),
		AutoTrade:  entity.GetAutoTrade(),
		Indicators: indicators,
		Strategies: strategies,
		Trades:     trades}
}
