// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserDAO_GetById(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)
	user, err := userDAO.GetById(1)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.GetId())
	assert.Equal(t, "test", user.GetUsername())
	assert.Equal(t, "USD", user.GetLocalCurrency())

	CleanupIntegrationTest()
}

func TestUserDAO_Create(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)
	err := userDAO.Create(&entity.User{
		Username:      "integrationtest",
		LocalCurrency: "USD"})
	assert.Nil(t, err)

	persistedUser, err := userDAO.GetById(2)
	assert.Nil(t, err)
	assert.NotNil(t, persistedUser)
	assert.Equal(t, "USD", persistedUser.GetLocalCurrency())
	assert.Empty(t, persistedUser.GetEtherbase())

	CleanupIntegrationTest()
}

func TestUserDAO_Update(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)
	err := userDAO.Create(&entity.User{
		Username:      "integrationtest",
		LocalCurrency: "GDP"})
	assert.Nil(t, err)

	persistedUser, err := userDAO.GetById(2)
	assert.Nil(t, err)
	assert.NotNil(t, persistedUser)
	assert.Equal(t, "integrationtest", persistedUser.GetUsername())
	assert.Equal(t, "GDP", persistedUser.GetLocalCurrency())
	assert.Empty(t, persistedUser.GetEtherbase())

	user := &entity.User{
		Id:            persistedUser.GetId(),
		Username:      persistedUser.GetUsername(),
		LocalCurrency: "YEN",
		Etherbase:     "0xabc123",
		Keystore:      "/some/path/to/keystore"}

	err = userDAO.Save(user)
	assert.Nil(t, err)

	persistedUser2, err := userDAO.GetById(2)
	assert.Nil(t, err)
	assert.Equal(t, user.GetLocalCurrency(), persistedUser2.GetLocalCurrency())
	assert.Equal(t, user.GetEtherbase(), persistedUser2.GetEtherbase())
	assert.Equal(t, user.GetKeystore(), persistedUser2.GetKeystore())

	CleanupIntegrationTest()
}
