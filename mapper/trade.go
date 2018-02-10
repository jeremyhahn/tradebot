package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type TradeMapper interface {
	MapTradeEntityToDto(entity entity.TradeEntity) common.Trade
	MapTradeDtoToEntity(dto common.Trade) entity.TradeEntity
}

type DefaultTradeMapper struct {
}

func NewTradeMapper() TradeMapper {
	return &DefaultTradeMapper{}
}

func (mapper *DefaultTradeMapper) MapTradeEntityToDto(entity entity.TradeEntity) common.Trade {
	return &dto.TradeDTO{
		Id:        entity.GetId(),
		ChartId:   entity.GetChartId(),
		UserId:    entity.GetUserId(),
		Base:      entity.GetBase(),
		Quote:     entity.GetQuote(),
		Exchange:  entity.GetExchangeName(),
		Date:      entity.GetDate(),
		Type:      entity.GetType(),
		Price:     entity.GetPrice(),
		Amount:    entity.GetAmount(),
		ChartData: entity.GetChartData()}
}

func (mapper *DefaultTradeMapper) MapTradeDtoToEntity(dto common.Trade) entity.TradeEntity {
	return &entity.Trade{
		Id:        dto.GetId(),
		ChartId:   dto.GetChartId(),
		UserId:    dto.GetUserId(),
		Base:      dto.GetBase(),
		Quote:     dto.GetQuote(),
		Exchange:  dto.GetExchange(),
		Date:      dto.GetDate(),
		Type:      dto.GetType(),
		Price:     dto.GetPrice(),
		Amount:    dto.GetAmount(),
		ChartData: dto.GetChartData()}
}
