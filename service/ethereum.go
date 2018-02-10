package service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
)

type EthereumService interface {
	Authenticate(address, passphrase string) error
	GetBalance(address string) (*big.Int, error)
	CreateAccount(passphrase string) (accounts.Account, error)
	DeleteAccount(passphrase string) error
	AuthService
}

type EthService struct {
	ctx      *common.Context
	client   *ethclient.Client
	keystore *ethkeystore.KeyStore
	userDAO  dao.UserDAO
	EthereumService
}

func NewEthereumService(ctx *common.Context, ipc, keystore string, userDAO dao.UserDAO) (EthereumService, error) {
	client, err := ethclient.Dial(ipc)
	if err != nil {
		return nil, err
	}
	return &EthService{
		ctx:     ctx,
		client:  client,
		userDAO: userDAO,
		keystore: ethkeystore.NewKeyStore(
			keystore,
			ethkeystore.LightScryptN,
			ethkeystore.LightScryptP)}, nil
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

func (eth *EthService) Login(password string) error {
	return eth.Authenticate("", password)
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
	return eth.userDAO.Save(user)
}
