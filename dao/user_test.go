// +build integration

package dao

import (
	"testing"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
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

func TestUserDAO_CreateGetUserExchange(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)
	mapper := mapper.NewUserMapper()

	userExchange := &entity.UserCryptoExchange{
		UserId: ctx.GetUser().GetId(),
		Name:   "Test Exchange",
		Key:    "ABC123",
		Secret: "$ecret!",
		URL:    "https://www.example.com",
		Extra:  "Anything specific to this exchange can be stored here"}

	err := userDAO.CreateExchange(userExchange)
	assert.Equal(t, nil, err)

	userContext := mapper.MapUserDtoToEntity(ctx.GetUser())
	persisted, exErr := userDAO.GetExchange(userContext, "Test Exchange")
	assert.Equal(t, nil, exErr)
	assert.Equal(t, userExchange.UserId, persisted.GetUserId())
	assert.Equal(t, userExchange.Name, persisted.GetName())
	assert.Equal(t, userExchange.URL, persisted.GetURL())
	assert.Equal(t, userExchange.Key, persisted.GetKey())
	assert.Equal(t, userExchange.Secret, persisted.GetSecret())
	assert.Equal(t, userExchange.Extra, persisted.GetExtra())

	CleanupIntegrationTest()
}

func TestUserDAO_GetTokens(t *testing.T) {
	ctx := NewIntegrationTestContext()
	userDAO := NewUserDAO(ctx)

	userDAO.CreateToken(&entity.UserToken{
		UserId:          ctx.GetUser().GetId(),
		Symbol:          "TEST",
		ContractAddress: "0xabc123",
		WalletAddress:   "0xdef456"})

	tokens := userDAO.GetTokens(&entity.User{Id: 1})

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "TEST", tokens[0].GetSymbol())
	assert.Equal(t, "0xabc123", tokens[0].GetContractAddress())
	assert.Equal(t, "0xdef456", tokens[0].GetWalletAddress())

	CleanupIntegrationTest()
}
