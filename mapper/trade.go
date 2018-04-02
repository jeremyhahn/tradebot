package mapper

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/shopspring/decimal"
)

type TradeMapper interface {
	MapTradeEntityToDto(entity entity.TradeEntity) common.Trade
	MapTradeDtoToEntity(dto common.Trade) entity.TradeEntity
}

type DefaultTradeMapper struct {
	ctx common.Context
}

func NewTradeMapper(ctx common.Context) TradeMapper {
	return &DefaultTradeMapper{ctx: ctx}
}

func (mapper *DefaultTradeMapper) MapTradeEntityToDto(entity entity.TradeEntity) common.Trade {
	amount, err := decimal.NewFromString(entity.GetAmount())
	if err != nil {
		mapper.ctx.GetLogger().Errorf("[Trademapper.MapTradeEntityToDto] Error parsing amount decimal: %s", err.Error())
	}
	price, err := decimal.NewFromString(entity.GetPrice())
	if err != nil {
		mapper.ctx.GetLogger().Errorf("[Trademapper.MapTradeEntityToDto] Error parsing price decimal: %s", err.Error())
	}
	return &dto.TradeDTO{
		Id:        entity.GetId(),
		ChartId:   entity.GetChartId(),
		UserId:    entity.GetUserId(),
		Base:      entity.GetBase(),
		Quote:     entity.GetQuote(),
		Exchange:  entity.GetExchangeName(),
		Date:      entity.GetDate(),
		Type:      entity.GetType(),
		Price:     price,
		Amount:    amount,
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
		Price:     dto.GetPrice().String(),
		Amount:    dto.GetAmount().String(),
		ChartData: dto.GetChartData()}
}
