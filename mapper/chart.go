package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type ChartMapper interface {
	MapChartDtoToEntity(dto common.Chart) dao.Chart
	MapChartEntityToDto(entity dao.ChartEntity) common.Chart
	MapIndicatorEntityToDto(entity dao.ChartIndicator) common.ChartIndicator
	MapIndicatorDtoToEntity(dto common.ChartIndicator) dao.ChartIndicator
	MapStrategyEntityToDto(entity dao.ChartStrategy) common.ChartStrategy
	MapStrategyDtoToEntity(dto common.ChartStrategy) dao.ChartStrategy
	MapTradeEntityToDto(entity dao.TradeEntity) common.Trade
	MapTradeDtoToEntity(trade common.Trade) dao.Trade
}

type ChartMapperImpl struct {
	ctx *common.Context
}

func NewChartMapper(ctx *common.Context) ChartMapper {
	return &ChartMapperImpl{ctx: ctx}
}

func (mapper *ChartMapperImpl) MapTradeEntityToDto(entity dao.TradeEntity) common.Trade {
	return common.Trade{
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

func (mapper *ChartMapperImpl) MapTradeDtoToEntity(trade common.Trade) dao.Trade {
	return dao.Trade{
		Id:        trade.Id,
		UserId:    mapper.ctx.User.Id,
		ChartId:   trade.ChartId,
		Date:      trade.Date,
		Exchange:  trade.Exchange,
		Type:      trade.Type,
		Base:      trade.Base,
		Quote:     trade.Quote,
		Amount:    trade.Amount,
		Price:     trade.Price,
		ChartData: trade.ChartData}
}

func (mapper *ChartMapperImpl) MapIndicatorEntityToDto(entity dao.ChartIndicator) common.ChartIndicator {
	return common.ChartIndicator{
		Id:         entity.Id,
		ChartId:    entity.ChartId,
		Name:       entity.Name,
		Parameters: entity.Parameters}
}

func (mapper *ChartMapperImpl) MapIndicatorDtoToEntity(dto common.ChartIndicator) dao.ChartIndicator {
	return dao.ChartIndicator{
		Id:         dto.Id,
		ChartId:    dto.ChartId,
		Name:       dto.Name,
		Parameters: dto.Parameters}
}

func (mapper *ChartMapperImpl) MapStrategyEntityToDto(entity dao.ChartStrategy) common.ChartStrategy {
	return common.ChartStrategy{
		Id:         entity.Id,
		ChartId:    entity.ChartId,
		Name:       entity.Name,
		Parameters: entity.Parameters}
}

func (mapper *ChartMapperImpl) MapStrategyDtoToEntity(dto common.ChartStrategy) dao.ChartStrategy {
	return dao.ChartStrategy{
		Id:         dto.Id,
		ChartId:    dto.ChartId,
		Name:       dto.Name,
		Parameters: dto.Parameters}
}

func (mapper *ChartMapperImpl) MapChartDtoToEntity(dto common.Chart) dao.Chart {
	var daoChartIndicators []dao.ChartIndicator
	for _, indicator := range dto.Indicators {
		daoChartIndicators = append(daoChartIndicators, mapper.MapIndicatorDtoToEntity(indicator))
	}
	var daoTrades []dao.Trade
	for _, trade := range dto.Trades {
		daoTrades = append(daoTrades, mapper.MapTradeDtoToEntity(trade))
	}
	return dao.Chart{
		Id:         dto.Id,
		UserId:     mapper.ctx.User.Id,
		Base:       dto.Base,
		Quote:      dto.Quote,
		Exchange:   dto.Exchange,
		Period:     dto.Period,
		AutoTrade:  dto.AutoTrade,
		Indicators: daoChartIndicators,
		Trades:     daoTrades}
}

func (mapper *ChartMapperImpl) MapChartEntityToDto(entity dao.ChartEntity) common.Chart {
	var indicators []common.ChartIndicator
	for _, indicator := range entity.GetIndicators() {
		indicators = append(indicators, mapper.MapIndicatorEntityToDto(indicator))
	}
	var trades []common.Trade
	for _, trade := range entity.GetTrades() {
		trades = append(trades, mapper.MapTradeEntityToDto(&trade))
	}
	return common.Chart{
		Id:         entity.GetId(),
		Base:       entity.GetBase(),
		Quote:      entity.GetQuote(),
		Exchange:   entity.GetExchangeName(),
		Period:     entity.GetPeriod(),
		AutoTrade:  entity.GetAutoTrade(),
		Indicators: indicators,
		Trades:     trades}
}
