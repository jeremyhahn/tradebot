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

func (service *UserService) GetExchange(user *common.User, name string) common.Exchange {
	exchanges := service.dao.GetExchanges(user)
	for _, ex := range exchanges {
		if ex.Name == name {
			currencyPair := exchange.CurrencyPairMap[name]
			return exchange.SupportedExchangeMap[name](&ex, service.ctx.Logger, currencyPair)
		}
	}
	return nil
}

func (service *UserService) GetExchanges(user *common.User) []common.CoinExchange {
	var exchangeList []common.CoinExchange
	var chans []chan common.CoinExchange
	exchanges := service.dao.GetExchanges(user)
	for _, ex := range exchanges {
		c := make(chan common.CoinExchange, 1)
		chans = append(chans, c)
		currencyPair := exchange.CurrencyPairMap[ex.Name]
		exchange := exchange.SupportedExchangeMap[ex.Name](&ex, service.ctx.Logger, currencyPair)
		go func() { c <- exchange.GetExchange() }()
	}
	for i := 0; i < len(exchanges); i++ {
		exchangeList = append(exchangeList, <-chans[i])
	}
	return exchangeList
}

func (service *UserService) GetWallets(user *common.User) []common.CryptoWallet {
	var walletList []common.CryptoWallet
	wallets := service.dao.GetWallets(user)
	var chans []chan common.CryptoWallet
	for _, _wallet := range wallets {
		wallet := _wallet
		c := make(chan common.CryptoWallet, 1)
		chans = append(chans, c)
		balance := service.getBalance(wallet.Currency, wallet.Address)
		netWorth := service.getPrice(wallet.Currency, balance)
		go func() {
			c <- common.CryptoWallet{
				Address:  wallet.Address,
				Currency: wallet.Currency,
				Balance:  service.getBalance(wallet.Currency, wallet.Address),
				NetWorth: netWorth}
		}()
	}
	for i := 0; i < len(wallets); i++ {
		walletList = append(walletList, <-chans[i])
	}
	return walletList
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
