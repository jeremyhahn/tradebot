package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/shopspring/decimal"
)

type ChartMapper interface {
	MapChartDtoToEntity(dto common.Chart) entity.ChartEntity
	MapChartEntityToDto(entity entity.ChartEntity) common.Chart
	MapIndicatorEntityToDto(entity entity.ChartIndicator) common.ChartIndicator
	MapIndicatorDtoToEntity(dto common.ChartIndicator) entity.ChartIndicator
	MapStrategyEntityToDto(entity entity.ChartStrategy) common.ChartStrategy
	MapStrategyDtoToEntity(dto common.ChartStrategy) entity.ChartStrategy
	MapTradeEntityToDto(entity entity.TradeEntity) common.Trade
	MapTradeDtoToEntity(trade common.Trade) entity.Trade
}

type DefaultChartMapper struct {
	ctx common.Context
}

func NewChartMapper(ctx common.Context) ChartMapper {
	return &DefaultChartMapper{ctx: ctx}
}

func (mapper *DefaultChartMapper) MapTradeEntityToDto(entity entity.TradeEntity) common.Trade {
	amount, err := decimal.NewFromString(entity.GetAmount())
	if err != nil {
		mapper.ctx.GetLogger().Errorf("[ChartMapper.MapTradeEntityToDto] Error parsing amount decimal: %s", err.Error())
	}
	price, err := decimal.NewFromString(entity.GetPrice())
	if err != nil {
		mapper.ctx.GetLogger().Errorf("[ChartMapper.MapTradeEntityToDto] Error parsing price decimal: %s", err.Error())
	}
	return &dto.TradeDTO{
		Id:        entity.GetId(),
		UserId:    entity.GetUserId(),
		ChartId:   entity.GetChartId(),
		Base:      entity.GetBase(),
		Quote:     entity.GetQuote(),
		Exchange:  entity.GetExchangeName(),
		Date:      entity.GetDate(),
		Type:      entity.GetType(),
		Amount:    amount,
		Price:     price,
		ChartData: entity.GetChartData()}
}

func (mapper *DefaultChartMapper) MapTradeDtoToEntity(trade common.Trade) entity.Trade {
	return entity.Trade{
		Id:        trade.GetId(),
		UserId:    mapper.ctx.GetUser().GetId(),
		ChartId:   trade.GetChartId(),
		Date:      trade.GetDate(),
		Exchange:  trade.GetExchange(),
		Type:      trade.GetType(),
		Base:      trade.GetBase(),
		Quote:     trade.GetQuote(),
		Amount:    trade.GetAmount().String(),
		Price:     trade.GetPrice().String(),
		ChartData: trade.GetChartData()}
}

func (mapper *DefaultChartMapper) MapIndicatorEntityToDto(entity entity.ChartIndicator) common.ChartIndicator {
	return &dto.ChartIndicatorDTO{
		Id:         entity.Id,
		ChartId:    entity.ChartId,
		Name:       entity.Name,
		Parameters: entity.Parameters}
}

func (mapper *DefaultChartMapper) MapIndicatorDtoToEntity(dto common.ChartIndicator) entity.ChartIndicator {
	return entity.ChartIndicator{
		Id:         dto.GetId(),
		ChartId:    dto.GetChartId(),
		Name:       dto.GetName(),
		Parameters: dto.GetParameters()}
}

func (mapper *DefaultChartMapper) MapStrategyEntityToDto(entity entity.ChartStrategy) common.ChartStrategy {
	return &dto.ChartStrategyDTO{
		Id:         entity.GetId(),
		ChartId:    entity.GetChartId(),
		Name:       entity.GetName(),
		Parameters: entity.GetParameters()}
}

func (mapper *DefaultChartMapper) MapStrategyDtoToEntity(dto common.ChartStrategy) entity.ChartStrategy {
	return entity.ChartStrategy{
		Id:         dto.GetId(),
		ChartId:    dto.GetChartId(),
		Name:       dto.GetName(),
		Parameters: dto.GetParameters()}
}

func (mapper *DefaultChartMapper) MapChartDtoToEntity(dto common.Chart) entity.ChartEntity {
	var daoChartIndicators []entity.ChartIndicator
	for _, indicator := range dto.GetIndicators() {
		daoChartIndicators = append(daoChartIndicators, mapper.MapIndicatorDtoToEntity(indicator))
	}
	var daoChartStrategies []entity.ChartStrategy
	for _, strategy := range dto.GetStrategies() {
		daoChartStrategies = append(daoChartStrategies, mapper.MapStrategyDtoToEntity(strategy))
	}
	var daoTrades []entity.Trade
	for _, trade := range dto.GetTrades() {
		daoTrades = append(daoTrades, mapper.MapTradeDtoToEntity(trade))
	}
	return &entity.Chart{
		Id:         dto.GetId(),
		UserId:     mapper.ctx.GetUser().GetId(),
		Base:       dto.GetBase(),
		Quote:      dto.GetQuote(),
		Exchange:   dto.GetExchange(),
		Period:     dto.GetPeriod(),
		AutoTrade:  dto.GetAutoTrade(),
		Indicators: daoChartIndicators,
		Strategies: daoChartStrategies,
		Trades:     daoTrades}
}

func (mapper *DefaultChartMapper) MapChartEntityToDto(entity entity.ChartEntity) common.Chart {
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
