// +build integration

package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserDAO_GetById(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)
	user := userDAO.GetById(1)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.Id)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "USD", user.LocalCurrency)

	CleanupIntegrationTest()
}
