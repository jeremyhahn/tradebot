package dto

import "github.com/jeremyhahn/tradebot/common"

type OrderPairDTO struct {
	BuyOrder  common.Order
	SellOrder common.Order
	common.OrderPair
}

func NewOrderPair() common.OrderPair {
	return &OrderPairDTO{}
}

func (op *OrderPairDTO) GetBuyOrder() common.Order {
	return op.BuyOrder
}

func (op *OrderPairDTO) GetSellOrder() common.Order {
	return op.BuyOrder
}
