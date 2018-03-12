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
	"github.com/shopspring/decimal"
)

type GethServiceImpl struct {
	ctx        common.Context
	client     *ethclient.Client
	keystore   *ethkeystore.KeyStore
	userDAO    dao.UserDAO
	userMapper mapper.UserMapper
	EthereumService
}

func NewGethService(ctx common.Context, userDAO dao.UserDAO, userMapper mapper.UserMapper) (GethService, error) {
	client, err := ethclient.Dial(ctx.GetIPC())
	if err != nil {
		return nil, err
	}
	return &GethServiceImpl{
		ctx:        ctx,
		client:     client,
		userDAO:    userDAO,
		userMapper: userMapper,
		keystore: ethkeystore.NewKeyStore(
			ctx.GetKeystore(),
			ethkeystore.StandardScryptN,
			ethkeystore.StandardScryptP)}, nil
}

func (geth *GethServiceImpl) CreateAccount(passphrase string) (common.UserContext, error) {
	acct, err := geth.keystore.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	return &dto.UserContextDTO{
		Etherbase: acct.Address.String(),
		Keystore:  acct.URL.String()}, nil
}

func (geth *GethServiceImpl) DeleteAccount(passphrase string) error {
	user := geth.ctx.GetUser()
	if user == nil {
		geth.ctx.GetLogger().Error("[GethService.DeleteAccount] No user context")
		return errors.New("No user context")
	}
	acct := accounts.Account{
		Address: ethcommon.HexToAddress(geth.ctx.GetUser().GetEtherbase()),
		URL:     accounts.URL{Path: geth.ctx.GetUser().GetKeystore()}}
	return geth.keystore.Delete(acct, passphrase)
}

func (geth *GethServiceImpl) GetAccounts() ([]common.UserContext, error) {
	geth.ctx.GetLogger().Error("[GethService.GetAccounts]")
	var accounts []common.UserContext
	user := geth.ctx.GetUser()
	if user == nil {
		geth.ctx.GetLogger().Error("[GethService.GetAccounts] No user context")
		return nil, errors.New("No user context")
	}
	for _, acct := range geth.keystore.Accounts() {
		accounts = append(accounts, &dto.UserContextDTO{
			Etherbase: acct.Address.String(),
			Keystore:  acct.URL.String()})
	}
	return accounts, nil
}

func (geth *GethServiceImpl) GetWallet(address string) (common.UserCryptoWallet, error) {
	geth.ctx.GetLogger().Errorf("[GethService.GetWallet] address: %s", address)
	ctx := context.Background()
	balance, err := geth.client.BalanceAt(ctx, ethcommon.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	fBalance, _ := new(big.Float).SetInt(balance).Float64()
	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  fBalance,
		Currency: "ETH",
		Value:    0}, nil
}

func (geth *GethServiceImpl) Authenticate(address, passphrase string) error {
	acct := accounts.Account{Address: ethcommon.HexToAddress(address)}
	return geth.keystore.Unlock(acct, passphrase)
}

func (geth *GethServiceImpl) GetTransactions(contractAddress string) ([]common.Transaction, error) {
	var transactions []common.Transaction
	return transactions, nil
}

func (geth *GethServiceImpl) GetToken(walletAddress string, contractAddress string) (common.EthereumToken, error) {
	var err error

	token, err := common.NewTokenCaller(ethcommon.HexToAddress(contractAddress), geth.client)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to instantiate a Token contract: %v", err)
		return nil, err
	}

	address := ethcommon.HexToAddress(walletAddress)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed hex address: %s, error: %s", walletAddress, err.Error())
		return nil, err
	}

	ethAmount, err := geth.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to get ethereum balance from contract %s. Error: %s", address, err.Error())
		return nil, err
	}

	balance, err := token.BalanceOf(nil, address)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to get balance from contract %s. Error: %s", contractAddress, err.Error())
		return nil, err
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to get symbol from contract %s. Error: %s", contractAddress, err.Error())
		return nil, err
	}
	tokenDecimals, err := token.Decimals(nil)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to get decimals from contract %s. Error: %s", contractAddress, err.Error())
		return nil, err
	}
	name, err := token.Name(nil)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[GethService.GetToken] Failed to retrieve token name from contract %s. Error: %s", contractAddress, err.Error())
		return nil, err
	}

	ethBalance, _ := decimal.NewFromString(ethAmount.String())
	ethFac, _ := decimal.NewFromString("0.000000000000000001")
	ethCorrected := ethBalance.Mul(ethFac)

	tokenBalance, err := decimal.NewFromString(balance.String())
	if err != nil {
		return nil, err
	}
	tokenMul := decimal.NewFromFloat(float64(0.1)).Pow(decimal.NewFromFloat(float64(tokenDecimals)))
	tokenCorrected := tokenBalance.Mul(tokenMul)
	fBalance, exact := tokenCorrected.Float64()
	if !exact {
		geth.ctx.GetLogger().Warningf("[GethService.GetToken] Error: tokenBalance float conversion not exact")
	}

	return &dto.EthereumTokenDTO{
		Name:            name,
		Symbol:          symbol,
		Balance:         fBalance,
		Decimals:        tokenDecimals,
		EthBalance:      ethCorrected.String(),
		ContractAddress: contractAddress,
		WalletAddress:   walletAddress}, nil
}

func (geth *GethServiceImpl) GetTokenTransactions(contractAddress string) ([]common.Transaction, error) {
	var transactions []common.Transaction
	return transactions, nil
}

func (geth *GethServiceImpl) Login(username, password string) (common.UserContext, error) {
	userEntity, err := geth.userDAO.GetByName(username)
	if err != nil {
		return nil, err
	}
	geth.ctx.GetLogger().Debugf("[GethService.Login] userEntity: %+v", userEntity)
	err = geth.Authenticate(userEntity.GetEtherbase(), password)
	if err != nil {
		geth.ctx.GetLogger().Errorf("[EhtereumService.Login] %s", err.Error())
		return nil, err
	}
	return geth.userMapper.MapUserEntityToDto(userEntity), err
}

func (geth *GethServiceImpl) Register(username, password string) error {
	_, err := geth.userDAO.GetByName(username)
	if err != nil && err.Error() != "record not found" {
		geth.ctx.GetLogger().Errorf("[GethService.Register] %s", err.Error())
		return errors.New("Unexpected error")
	}
	acct, err := geth.CreateAccount(password)
	if err != nil {
		return err
	}
	newUserEntity := &entity.User{
		Username:      username,
		LocalCurrency: "USD",
		Etherbase:     acct.GetEtherbase(),
		Keystore:      acct.GetKeystore()}
	err = geth.userDAO.Save(newUserEntity)
	userDTO := geth.userMapper.MapUserEntityToDto(newUserEntity)
	geth.ctx.SetUser(userDTO)
	if err != nil {
		geth.DeleteAccount(password)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("User already exists")
		}
		geth.ctx.GetLogger().Errorf("[GethService.Register] %s", err.Error())
		return errors.New("Unexpected error")
	}
	return err
}
