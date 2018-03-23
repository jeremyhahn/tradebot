package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultUserService struct {
	ctx                common.Context
	userDAO            dao.UserDAO
	marketcapService   MarketCapService
	ethereumService    EthereumService
	userMapper         mapper.UserMapper
	userExchangeMapper mapper.UserExchangeMapper
	pluginService      PluginService
	UserService
}

func NewUserService(ctx common.Context, userDAO dao.UserDAO,
	userMapper mapper.UserMapper, userExchangeMapper mapper.UserExchangeMapper,
	marketcapService MarketCapService, ethereumService EthereumService, pluginService PluginService) UserService {
	return &DefaultUserService{
		ctx:                ctx,
		userDAO:            userDAO,
		marketcapService:   marketcapService,
		ethereumService:    ethereumService,
		userMapper:         userMapper,
		userExchangeMapper: userExchangeMapper,
		pluginService:      pluginService}
}

func (service *DefaultUserService) CreateUser(user common.UserContext) {
	service.userDAO.Create(&entity.User{
		Username: user.GetUsername()})
}

func (service *DefaultUserService) GetCurrentUser() (common.UserContext, error) {
	entity, err := service.userDAO.GetById(service.ctx.GetUser().GetId())
	if err != nil {
		return nil, err
	}
	return service.userMapper.MapUserEntityToDto(entity), nil
}

func (service *DefaultUserService) GetUserById(userId uint) (common.UserContext, error) {
	entity, err := service.userDAO.GetById(userId)
	if err != nil {
		return nil, err
	}
	return service.userMapper.MapUserEntityToDto(entity), nil
}

func (service *DefaultUserService) GetUserByName(username string) (common.UserContext, error) {
	entity, err := service.userDAO.GetByName(username)
	if err != nil {
		return nil, err
	}
	return service.userMapper.MapUserEntityToDto(entity), nil
}

func (service *DefaultUserService) GetExchange(user common.UserContext, name string, currencyPair *common.CurrencyPair) (common.Exchange, error) {
	daoUser := &entity.User{Id: user.GetId()}
	exchanges := service.userDAO.GetExchanges(daoUser)
	for _, ex := range exchanges {
		if ex.Name == name {
			exchange, err := NewExchangeService(service.ctx, service.userDAO, service.userMapper,
				service.userExchangeMapper, service.pluginService).CreateExchange(ex.Name)
			if err != nil {
				service.ctx.GetLogger().Debugf("[UserService.GetExchange] Error: %s", err.Error())
				return nil, err
			}
			return exchange, nil
		}
	}
	return nil, errors.New("Exchange not found")
}

func (service *DefaultUserService) GetConfiguredExchanges() []common.UserCryptoExchange {
	var exchanges []common.UserCryptoExchange
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, service.userExchangeMapper.MapEntityToDto(&ex))
	}
	return exchanges
}

func (service *DefaultUserService) GetExchangeSummary(currencyPair *common.CurrencyPair) ([]common.CryptoExchangeSummary, error) {
	user := service.ctx.GetUser()
	service.ctx.GetLogger().Debugf("[UserService.GetExchangeSummary] %+v, %+v", user, currencyPair)
	var exchangeList []common.CryptoExchangeSummary
	var chans []chan common.CryptoExchangeSummary
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	exchanges := service.userDAO.GetExchanges(daoUser)
	c := make(chan common.CryptoExchangeSummary, len(exchanges))
	for _, ex := range exchanges {
		chans = append(chans, c)
		exchangeService := NewExchangeService(service.ctx, service.userDAO, service.userMapper,
			service.userExchangeMapper, service.pluginService)
		exchange, err := exchangeService.CreateExchange(ex.GetName())
		if err != nil {
			service.ctx.GetLogger().Debugf("[UserService.GetExchangeSummary] Error: %s", err.Error())
			return nil, err
		}
		go func() { c <- exchange.GetSummary() }()
	}
	for i := 0; i < len(exchanges); i++ {
		exchangeList = append(exchangeList, <-chans[i])
	}
	service.ctx.GetLogger().Debugf("[UserService.GetExchanges] %+v", exchangeList)
	return exchangeList, nil
}

func (service *DefaultUserService) GetWallets() []common.UserCryptoWallet {
	var walletList []common.UserCryptoWallet
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	wallets := service.userDAO.GetWallets(daoUser)
	var chans []chan common.UserCryptoWallet
	for _, _wallet := range wallets {
		wallet := _wallet
		c := make(chan common.UserCryptoWallet, 1)
		chans = append(chans, c)
		go func() {
			balance, err := service.getBalance(wallet.Currency, wallet.Address)
			if err != nil {
				service.ctx.GetLogger().Errorf("[UserService.GetWallets] Error: %s", err.Error())
			}
			c <- &dto.UserCryptoWalletDTO{
				Address:  wallet.Address,
				Currency: wallet.Currency,
				Balance:  balance,
				Value:    balance * service.getPrice(wallet.Currency, balance)}
		}()
	}
	for i := 0; i < len(wallets); i++ {
		walletList = append(walletList, <-chans[i])
	}
	return walletList
}

func (service *DefaultUserService) CreateToken(token common.EthereumToken) error {
	return service.userDAO.CreateToken(&entity.UserToken{
		UserId:          service.ctx.GetUser().GetId(),
		Symbol:          token.GetSymbol(),
		ContractAddress: token.GetContractAddress(),
		WalletAddress:   token.GetWalletAddress()})
}

func (service *DefaultUserService) GetToken(symbol string) (common.EthereumToken, error) {
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	token := service.userDAO.GetToken(daoUser, symbol)
	//return service.ethereumService.GetTokenBalance(token.GetContract(), token.GetWallet())
	ethereumToken, err := service.ethereumService.GetToken(token.GetWalletAddress(), token.GetContractAddress())
	if err != nil {
		return nil, err
	}
	return ethereumToken, nil
}

func (service *DefaultUserService) GetTokens(wallet string) ([]common.EthereumToken, error) {
	var walletTokens []common.EthereumToken
	tokens, err := service.GetAllTokens()
	if err != nil {
		return nil, err
	}
	for _, t := range tokens {
		if t.GetWalletAddress() == wallet {
			walletTokens = append(walletTokens, t)
		}
	}
	return walletTokens, nil
}

func (service *DefaultUserService) GetAllTokens() ([]common.EthereumToken, error) {
	var tokenList []common.EthereumToken
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	tokens := service.userDAO.GetTokens(daoUser)
	var chans []chan common.EthereumToken
	for _, token := range tokens {
		_token := token
		c := make(chan common.EthereumToken, 1)
		chans = append(chans, c)
		go func() {
			tokenDTO, err := service.ethereumService.GetToken(_token.GetWalletAddress(), _token.GetContractAddress())
			if err != nil {
				service.ctx.GetLogger().Errorf("[UserService.GetTokens] Error: %s", err.Error())
			}
			c <- tokenDTO
		}()
	}
	for i := 0; i < len(tokens); i++ {
		tokenList = append(tokenList, <-chans[i])
	}
	return tokenList, nil
}

func (service *DefaultUserService) CreateWallet(wallet common.UserCryptoWallet) error {
	return service.userDAO.CreateWallet(&entity.UserWallet{
		UserId:   service.ctx.GetUser().GetId(),
		Address:  wallet.GetAddress(),
		Currency: wallet.GetCurrency()})
}

func (service *DefaultUserService) GetWallet(currency string) common.UserCryptoWallet {
	service.ctx.GetLogger().Debugf("[UserService.GetWallet] user: %s, currency: %s",
		service.ctx.GetUser().GetUsername(), currency)
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	wallet := service.userDAO.GetWallet(daoUser, currency)
	balance, err := service.getBalance(wallet.GetCurrency(), wallet.GetAddress())
	if err != nil {
		service.ctx.GetLogger().Errorf("[UserService.GetWallet] Error: %s", err.Error())
	}
	return &dto.UserCryptoWalletDTO{
		Address:  wallet.GetAddress(),
		Currency: wallet.GetCurrency(),
		Balance:  balance,
		Value:    service.getPrice(wallet.GetCurrency(), balance)}
}

func (service *DefaultUserService) getBalance(currency, address string) (float64, error) {
	// TODO: Replace with plugin service factory method
	service.ctx.GetLogger().Debugf("[UserService.getBalance] currency=%s, address=%s", currency, address)
	if currency == "XRP" {
		rippleService := NewRippleService(service.ctx, service.userDAO, service.marketcapService)
		wallet, err := rippleService.GetWallet(address)
		if err != nil {
			return 0.0, err
		}
		return wallet.GetBalance(), nil
	} else if currency == "BTC" {
		return NewBlockchainInfo(service.ctx).GetBalance(address).GetBalance(), nil
	} else if currency == "ETH" {
		wallet, err := service.ethereumService.GetWallet(address)
		return wallet.GetBalance(), err
	}
	return 0.0, errors.New(fmt.Sprintf("Unknown currency: %s", currency))
}

func (service *DefaultUserService) getPrice(currency string, amt float64) float64 {
	service.ctx.GetLogger().Debugf("[UserService.getPrice] currency=%s, amt=%.8f", currency, amt)
	/*
		if currency == "BTC" {
			return util.TruncateFloat(NewBlockchainInfo(service.ctx).GetPrice()*amt, 8)
		}*/
	f, _ := strconv.ParseFloat(service.marketcapService.GetMarket(currency).PriceUSD, 64)
	return f
}
