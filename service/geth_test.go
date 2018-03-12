// +build integration

package service

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/stretchr/testify/assert"
)

var GETH_PASSPHRASE = "test"
var GETH_ETHERBASE = "0x411e50dde8844a77323849f5031be52c1f592383"
var GETH_BALANCE = "42000000000000000000"

func TestGethService_CreateDeleteAccount(t *testing.T) {
	ctx := NewIntegrationTestContext()

	userMapper := mapper.NewUserMapper()
	userDAO := dao.NewUserDAO(ctx)
	service, err := NewGethService(ctx, userDAO, userMapper)
	assert.Nil(t, err)

	err = service.Register("testuser", GETH_PASSPHRASE)
	assert.Equal(t, nil, err)
	assert.NotNil(t, ctx.GetUser().GetEtherbase())

	err = service.DeleteAccount(GETH_PASSPHRASE)
	assert.Equal(t, nil, err)

	CleanupIntegrationTest()
}

func TestGethService_Authenticate(t *testing.T) {
	ctx := NewIntegrationTestContext()

	userMapper := mapper.NewUserMapper()
	userDAO := dao.NewUserDAO(ctx)
	service, err := NewGethService(ctx, userDAO, userMapper)
	assert.Equal(t, nil, err)

	testuser := "testuser"

	err = service.Register(testuser, GETH_PASSPHRASE)
	assert.Equal(t, nil, err)

	err = service.Authenticate(ctx.GetUser().GetEtherbase(), GETH_PASSPHRASE)
	assert.Nil(t, err)

	err = service.Authenticate(ctx.GetUser().GetEtherbase(), "nogood")
	assert.NotNil(t, err)
	assert.Equal(t, "could not decrypt key with given passphrase", err.Error())

	user, err := service.Login(testuser, GETH_PASSPHRASE)
	assert.Nil(t, err)
	assert.Equal(t, uint(2), user.GetId())
	assert.Equal(t, "testuser", user.GetUsername())

	user, err = service.Login(testuser, "badpass")
	assert.NotNil(t, err)
	assert.Equal(t, "could not decrypt key with given passphrase", err.Error())

	err = service.DeleteAccount(GETH_PASSPHRASE)
	assert.Equal(t, nil, err)

	CleanupIntegrationTest()
}

/*
func TestGethService_GetBalance(t *testing.T) {
	ctx := NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	service, err := NewGethService(ctx, userDAO)
	assert.Equal(t, nil, err)

	err = service.Register("testuser", GETH_PASSPHRASE)
	assert.Equal(t, nil, err)

	ctx.Logger.Debugf("Make sure you've got an account and mined some coins, and GETH_ETHERBASE holds the valid address...")

	balance, err := service.GetBalance(GETH_ETHERBASE)
	//expected := new(big.Int)
	//expected.SetString(GETH_BALANCE, 10)
	//assert.Equal(t, expected, balance)
	assert.Equal(t, true, balance.Int64() > 0)

	CleanupIntegrationTest()
}*/

func TestGethService_Register(t *testing.T) {
	ctx := NewIntegrationTestContext()

	userMapper := mapper.NewUserMapper()
	userDAO := dao.NewUserDAO(ctx)
	service, err := NewGethService(ctx, userDAO, userMapper)
	assert.Nil(t, err)

	err = service.Register("ethtest", GETH_PASSPHRASE)
	assert.Nil(t, err)

	assert.NotNil(t, ctx.GetUser().GetKeystore())
	assert.NotNil(t, ctx.GetUser().GetEtherbase())

	err = service.DeleteAccount(GETH_PASSPHRASE)
	assert.Nil(t, err)

	newUser, err := userDAO.GetByName("ethtest")
	assert.Nil(t, err)
	assert.Equal(t, true, newUser.GetId() > 0)
	assert.Equal(t, "ethtest", newUser.GetUsername())
	assert.NotNil(t, newUser.GetEtherbase())
	assert.NotNil(t, newUser.GetKeystore())

	CleanupIntegrationTest()
}

func getGethBalance(t *testing.T) *big.Int {
	bal, err := new(big.Int).SetString(GETH_BALANCE, 10)
	assert.Nil(t, err)
	return bal
}

func createGethSimulatedBackend() *backends.SimulatedBackend {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(150000000)}
	return backends.NewSimulatedBackend(alloc)
}
