// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestPriceHistoryService(t *testing.T) {
	ctx := test.NewUnitTestContext()
	service := NewPriceHistoryService(ctx)
	data := service.GetPriceHistory("cardano")
	assert.NotNil(t, data)
	assert.Equal(t, true, data[0].GetClose() > 0)
}
