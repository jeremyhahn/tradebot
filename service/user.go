package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jeremyhahn/tradebot/util"
)

type UserService struct {
	ctx *common.Context
	dao *dao.UserDAO
}

func NewUserService(ctx *common.Context, dao *dao.UserDAO) *UserService {
	return &UserService{
		ctx: ctx,
		dao: dao}
}

func (service *UserService) CreateUser(user *common.User) {
	service.dao.Create(&dao.User{
		Username: user.Username})
}

func (service *UserService) GetCurrenctUser() *common.User {
	return service.dao.GetById(service.ctx.User.Id)
}

func (service *UserService) GetUserById(userId uint) *common.User {
	return service.dao.GetById(userId)
}

func (service *UserService) GetUserByName(username string) *common.User {
	return service.dao.GetByName(username)
}

func (service *UserService) GetExchanges(user *common.User) []common.CoinExchange {
	var exchangeList []common.CoinExchange
	for _, ex := range service.dao.GetExchanges(user) {
		currencyPair := exchange.CurrencyPairMap[ex.Name]
		exchange := exchange.SupportedExchangeMap[ex.Name](&ex, service.ctx.Logger, currencyPair)
		exchangeList = append(exchangeList, exchange.GetExchange())
	}
	return exchangeList
}

func (service *UserService) GetExchangesAsync(user *common.User) chan []common.CoinExchange {
	exchangeChan := make(chan []common.CoinExchange)
	go func() {
		var exchangeList []common.CoinExchange
		var queued []chan common.CoinExchange
		for _, ex := range service.dao.GetExchanges(user) {
			currencyPair := exchange.CurrencyPairMap[ex.Name]
			exchange := exchange.SupportedExchangeMap[ex.Name](&ex, service.ctx.Logger, currencyPair)
			queued = append(queued, exchange.GetExchangeAsync())
		}
		for _, q := range queued {
			exchangeList = append(exchangeList, <-q)
		}
		exchangeChan <- exchangeList
	}()
	return exchangeChan
}

func (service *UserService) GetWallets(user *common.User) []common.CryptoWallet {
	var cryptoWallets []common.CryptoWallet
	wallets := service.dao.GetWallets(user)
	for _, wallet := range wallets {
		balance := service.getBalance(wallet.Currency, wallet.Address)
		netWorth := service.getPrice(wallet.Currency, balance)
		cryptoWallets = append(cryptoWallets, common.CryptoWallet{
			Address:  wallet.Address,
			Currency: wallet.Currency,
			Balance:  service.getBalance(wallet.Currency, wallet.Address),
			NetWorth: netWorth})
	}
	return cryptoWallets
}

func (service *UserService) GetWalletsAsync(user *common.User) chan []common.CryptoWallet {
	var walletList []common.CryptoWallet
	walletListChan := make(chan []common.CryptoWallet)
	walletChan := make(chan common.CryptoWallet)
	wallets := service.dao.GetWallets(user)
	for _, _wallet := range wallets {
		wallet := _wallet
		go func() {
			balance := service.getBalance(wallet.Currency, wallet.Address)
			netWorth := service.getPrice(wallet.Currency, balance)
			walletChan <- common.CryptoWallet{
				Address:  wallet.Address,
				Currency: wallet.Currency,
				Balance:  service.getBalance(wallet.Currency, wallet.Address),
				NetWorth: netWorth}
		}()
		walletList = append(walletList, <-walletChan)
	}
	go func() { walletListChan <- walletList }()
	return walletListChan
}

func (service *UserService) GetWallet(user *common.User, currency string) *common.CryptoWallet {
	wallet := service.dao.GetWallet(user, currency)
	balance := service.getBalance(wallet.Currency, wallet.Address)
	return &common.CryptoWallet{
		Address:  wallet.Address,
		Currency: wallet.Currency,
		Balance:  balance,
		NetWorth: service.getPrice(wallet.Currency, balance)}
}

func (service *UserService) GetWalletAsync(user *common.User, currency string) chan common.CryptoWallet {
	walletChan := make(chan common.CryptoWallet, 1)
	go func() {
		wallet := service.dao.GetWallet(user, currency)
		balance := service.getBalance(wallet.Currency, wallet.Address)
		walletChan <- common.CryptoWallet{
			Address:  wallet.Address,
			Currency: wallet.Currency,
			Balance:  balance,
			NetWorth: service.getPrice(wallet.Currency, balance)}
	}()
	return walletChan
}

func (service *UserService) getBalance(currency, address string) float64 {
	service.ctx.Logger.Debugf("[UserService.getBalance] currency=%s, address=%s", currency, address)
	if currency == "XRP" {
		return NewRipple(service.ctx).GetBalance(address).Balance
	} else if currency == "BTC" {
		return NewBlockchainInfo(service.ctx).GetBalance(address).Balance
	}
	return 0.0
}

func (service *UserService) getPrice(currency string, amt float64) float64 {
	service.ctx.Logger.Debugf("[UserService.getPrice] currency=%s, amt=%.8f", currency, amt)
	if currency == "BTC" {
		return util.TruncateFloat(NewBlockchainInfo(service.ctx).GetPrice()*amt, 8)
	}
	return 0.0
}
