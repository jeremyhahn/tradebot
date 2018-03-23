package test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestShopSpring_gasUsed(t *testing.T) {
	gasUsed, err := decimal.NewFromString("25153")
	assert.Nil(t, err)
	answer := gasUsed.Div(decimal.NewFromFloat(100000000))

	expected := float64(25153) / 100000000

	assert.Equal(t, gasUsed.StringFixed(0), "25153")
	assert.Equal(t, fmt.Sprintf("%.8f", expected), answer.StringFixed(8))
}

func TestShopSpring_gasPrice(t *testing.T) {
	gasPrice, err := decimal.NewFromString("57631341878")
	assert.Nil(t, err)
	answer := gasPrice.Div(decimal.NewFromFloat(1000000000000000000))

	expected := float64(57631341878) / 1000000000000000000

	assert.Equal(t, gasPrice.StringFixed(0), "57631341878")
	assert.Equal(t, fmt.Sprintf("%.8f", expected), answer.StringFixed(8))
}
