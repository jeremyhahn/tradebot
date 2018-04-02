// +build integration

package service

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSlickChartsService(t *testing.T) {
	ctx := test.NewUnitTestContext()
	service := NewSlickChartsService(ctx)
	candlestick, err := service.GetPriceAt("ADA", time.Now())
	assert.Nil(t, err)
	assert.Equal(t, true, candlestick.Close.GreaterThan(decimal.NewFromFloat(0)))
}
