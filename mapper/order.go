package mapper

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type OrderMapper interface {
	MapOrderEntityToDto(entity entity.OrderEntity) common.Order
	MapOrderDtoToEntity(dto common.Order) entity.OrderEntity
	MapOrderDtoToViewModel(dto common.Order) viewmodel.Order
}

type DefaultOrderMapper struct {
	ctx *common.Context
	OrderMapper
}

func NewOrderMapper(ctx *common.Context) OrderMapper {
	return &DefaultOrderMapper{
		ctx: ctx}
}

func (mapper *DefaultOrderMapper) MapOrderEntityToDto(entity entity.OrderEntity) common.Order {
	return &dto.OrderDTO{
		Id:           fmt.Sprintf("%d", entity.GetId()),
		Exchange:     entity.GetExchange(),
		Date:         entity.GetDate(),
		Type:         entity.GetType(),
		CurrencyPair: common.NewCurrencyPair(entity.GetCurrency(), mapper.ctx.GetUser().GetLocalCurrency()),
		Quantity:     entity.GetQuantity(),
		Price:        entity.GetPrice(),
		Fee:          entity.GetFee(),
		Total:        entity.GetTotal()}
}

func (mapper *DefaultOrderMapper) MapOrderDtoToEntity(dto common.Order) entity.OrderEntity {
	return &entity.Order{
		UserId:   mapper.ctx.GetUser().GetId(),
		Date:     dto.GetDate(),
		Exchange: dto.GetExchange(),
		Type:     dto.GetType(),
		Currency: dto.GetCurrencyPair().String(),
		Quantity: dto.GetQuantity(),
		Price:    dto.GetPrice(),
		Fee:      dto.GetFee(),
		Total:    dto.GetTotal()}
}

func (mapper *DefaultOrderMapper) MapOrderDtoToViewModel(dto common.Order) viewmodel.Order {
	return viewmodel.Order{
		Id:            dto.GetId(),
		Exchange:      dto.GetExchange(),
		Date:          dto.GetDate().Format(common.TIME_DISPLAY_FORMAT),
		Type:          dto.GetType(),
		CurrencyPair:  dto.GetCurrencyPair(),
		Quantity:      dto.GetQuantity(),
		Price:         dto.GetPrice(),
		Fee:           dto.GetFee(),
		Total:         dto.GetTotal(),
		PriceCurrency: dto.GetPriceCurrency(),
		FeeCurrency:   dto.GetFeeCurrency(),
		TotalCurrency: dto.GetTotalCurrency()}
}
