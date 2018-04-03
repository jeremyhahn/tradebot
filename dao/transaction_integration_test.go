// +build integration

package dao

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransactionDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()

	orderDAO := NewTransactionDAO(ctx)
	order1 := &entity.Transaction{
		Id:       "1",
		UserId:   1,
		Date:     time.Now(),
		Network:  "Test",
		Type:     "buy",
		Currency: "TST-USD",
		Quantity: decimal.NewFromFloat(25.67).StringFixed(2),
		Price:    decimal.NewFromFloat(123.45).StringFixed(2),
		Fee:      decimal.NewFromFloat(1.23).StringFixed(2),
		Total:    decimal.NewFromFloat(124.67).StringFixed(2)}
	order2 := &entity.Transaction{
		Id:       "2",
		UserId:   1,
		Date:     time.Now(),
		Network:  "Test 2",
		Type:     "buy",
		Currency: "TST-USD",
		Quantity: decimal.NewFromFloat(25.67).StringFixed(2),
		Price:    decimal.NewFromFloat(123.45).StringFixed(2),
		Fee:      decimal.NewFromFloat(1.23).StringFixed(2),
		Total:    decimal.NewFromFloat(124.67).StringFixed(2)}

	err := orderDAO.Create(order1)
	assert.Nil(t, err)

	err = orderDAO.Create(order2)
	assert.Nil(t, err)

	orders, err := orderDAO.Find()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(orders))
	assert.Equal(t, "1", orders[0].GetId())
	assert.Equal(t, uint(1), orders[0].GetUserId())
	//	assert.Equal(t, order1.GetDate(), orders[0].GetDate())
	assert.Equal(t, order1.GetNetwork(), orders[0].GetNetwork())
	assert.Equal(t, order1.GetType(), orders[0].GetType())
	assert.Equal(t, order1.GetQuantity(), orders[0].GetQuantity())
	assert.Equal(t, order1.GetPrice(), orders[0].GetPrice())
	assert.Equal(t, order1.GetFee(), orders[0].GetFee())
	assert.Equal(t, order1.GetTotal(), orders[0].GetTotal())

	assert.Equal(t, "2", orders[1].GetId())
	assert.Equal(t, uint(1), orders[1].GetUserId())
	//	assert.Equal(t, order2.GetDate(), orders[1].GetDate())
	assert.Equal(t, order2.GetNetwork(), orders[1].GetNetwork())
	assert.Equal(t, order2.GetType(), orders[1].GetType())
	assert.Equal(t, order2.GetQuantity(), orders[1].GetQuantity())
	assert.Equal(t, order2.GetPrice(), orders[1].GetPrice())
	assert.Equal(t, order2.GetFee(), orders[1].GetFee())
	assert.Equal(t, order2.GetTotal(), orders[1].GetTotal())

	CleanupIntegrationTest()
}
