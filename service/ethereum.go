package service

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type EthereumService interface {
	Authenticate(address, passphrase string) error
	GetBalance(address string) (*big.Int, error)
	CreateAccount(passphrase string) (accounts.Account, error)
	DeleteAccount(passphrase string) error
	AuthService
}

type EthService struct {
	ctx        *common.Context
	client     *ethclient.Client
	keystore   *ethkeystore.KeyStore
	userDAO    dao.UserDAO
	userMapper mapper.UserMapper
	EthereumService
}

func NewEthereumService(ctx *common.Context, ipc, keystore string, userDAO dao.UserDAO,
	userMapper mapper.UserMapper) (EthereumService, error) {
	client, err := ethclient.Dial(ipc)
	if err != nil {
		return nil, err
	}
	return &EthService{
		ctx:        ctx,
		client:     client,
		userDAO:    userDAO,
		userMapper: userMapper,
		keystore: ethkeystore.NewKeyStore(
			keystore,
			ethkeystore.StandardScryptN,
			ethkeystore.StandardScryptP)}, nil
}

func (eth *EthService) CreateAccount(passphrase string) (accounts.Account, error) {
	return eth.keystore.NewAccount(passphrase)
}

func (eth *EthService) DeleteAccount(passphrase string) error {
	acct := accounts.Account{
		Address: ethcommon.HexToAddress(eth.ctx.GetUser().GetEtherbase()),
		URL:     accounts.URL{Path: eth.ctx.GetUser().GetKeystore()}}
	return eth.keystore.Delete(acct, passphrase)
}

func (eth *EthService) GetBalance(address string) (*big.Int, error) {
	ctx := context.Background()
	return eth.client.BalanceAt(ctx, ethcommon.HexToAddress(address), nil)
}

func (eth *EthService) Authenticate(address, passphrase string) error {
	acct := accounts.Account{Address: ethcommon.HexToAddress(address)}
	return eth.keystore.Unlock(acct, passphrase)
}

func (eth *EthService) Login(username, password string) (common.User, error) {
	userEntity, err := eth.userDAO.GetByName(username)
	if err != nil {
		return nil, err
	}
	eth.ctx.Logger.Debugf("[EthereumService.Login] userEntity: %+v", userEntity)
	err = eth.Authenticate(userEntity.GetEtherbase(), password)
	if err != nil {
		eth.ctx.Logger.Errorf("[EhtereumService.Login] %s", err.Error())
		return nil, err
	}
	userDTO := eth.userMapper.MapUserEntityToDto(userEntity)
	eth.ctx.SetUser(userDTO)
	return userDTO, err
}

func (eth *EthService) Register(username, password string) error {
	acct, err := eth.CreateAccount(password)
	if err != nil {
		return err
	}
	user := &entity.User{
		Username:      username,
		LocalCurrency: "USD",
		Etherbase:     acct.Address.String(),
		Keystore:      acct.URL.String()}
	eth.ctx.SetUser(&dto.UserDTO{
		Id:            user.GetId(),
		Username:      user.GetUsername(),
		LocalCurrency: user.GetLocalCurrency(),
		Etherbase:     user.GetEtherbase(),
		Keystore:      user.GetKeystore()})
	err = eth.userDAO.Save(user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			err = errors.New("User already exists")
		}
	}
	return err
}
