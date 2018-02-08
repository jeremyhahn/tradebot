package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
)

type TradeMapper interface {
	MapTradeEntityToDto(entity dao.TradeEntity) common.Trade
	MapTradeDtoToEntity(dto common.Trade) dao.TradeEntity
}

type DefaultTradeMapper struct {
}

func NewTradeMapper() TradeMapper {
	return &DefaultTradeMapper{}
}

func (mapper *DefaultTradeMapper) MapTradeEntityToDto(entity dao.TradeEntity) common.Trade {
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

func (mapper *DefaultTradeMapper) MapTradeDtoToEntity(dto common.Trade) dao.TradeEntity {
	return &dao.Trade{
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
