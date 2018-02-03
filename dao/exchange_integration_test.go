package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeDAO(t *testing.T) {
	ctx := NewIntegrationTestContext()
	exchangeDAO := NewExchangeDAO(ctx)

	exchange := &CryptoExchange{
		UserID: ctx.User.Id,
		Name:   "Test Exchange",
		Key:    "ABC123",
		Secret: "$ecret!",
		URL:    "https://www.example.com",
		Extra:  "Anything specific to this exchange can be stored here"}
	err := exchangeDAO.Create(exchange)
	assert.Equal(t, nil, err)

	persisted, exErr := exchangeDAO.Get("Test Exchange")
	assert.Equal(t, nil, exErr)
	assert.Equal(t, exchange.UserID, persisted.UserID)
	assert.Equal(t, exchange.Name, persisted.Name)
	assert.Equal(t, exchange.URL, persisted.URL)
	assert.Equal(t, exchange.Key, persisted.Key)
	assert.Equal(t, exchange.Secret, persisted.Secret)
	assert.Equal(t, exchange.Extra, persisted.Extra)

	CleanupIntegrationTest()
}
