package service

import (
	"errors"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultUserService struct {
	ctx                common.Context
	userDAO            dao.UserDAO
	marketcapService   common.MarketCapService
	ethereumService    EthereumService
	userMapper         mapper.UserMapper
	userExchangeMapper mapper.UserExchangeMapper
	exchangeService    ExchangeService
	walletService      WalletService
	UserService
}

func NewUserService(ctx common.Context, userDAO dao.UserDAO,
	userMapper mapper.UserMapper, userExchangeMapper mapper.UserExchangeMapper,
	marketcapService common.MarketCapService, ethereumService EthereumService,
	exchangeService ExchangeService, walletService WalletService) UserService {
	return &DefaultUserService{
		ctx:                ctx,
		userDAO:            userDAO,
		marketcapService:   marketcapService,
		ethereumService:    ethereumService,
		userMapper:         userMapper,
		userExchangeMapper: userExchangeMapper,
		exchangeService:    exchangeService,
		walletService:      walletService}
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
			exchange, err := service.exchangeService.CreateExchange(ex.Name)
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

func (service *DefaultUserService) CreateExchange(userCryptoExchange common.UserCryptoExchange) (common.UserCryptoExchange, error) {
	entity := service.userExchangeMapper.MapDtoToEntity(userCryptoExchange)
	err := service.userDAO.CreateExchange(entity)
	return userCryptoExchange, err
}

func (service *DefaultUserService) DeleteExchange(exchangeName string) error {
	user := service.ctx.GetUser()
	userEntity := service.userMapper.MapUserDtoToEntity(user)
	userExchange, _ := service.userDAO.GetExchange(userEntity, exchangeName)
	if userExchange != nil {
		return service.userDAO.DeleteExchange(userExchange)
	}
	return errors.New("Exchange not found")
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
		exchange, err := service.exchangeService.CreateExchange(ex.Name)
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

func (service *DefaultUserService) GetTokensFor(wallet string) ([]common.EthereumToken, error) {
	var walletTokens []common.EthereumToken
	tokens, err := service.GetTokens()
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

func (service *DefaultUserService) GetTokens() ([]common.EthereumToken, error) {
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

func (service *DefaultUserService) GetWallet(currency string) (common.UserCryptoWallet, error) {
	service.ctx.GetLogger().Debugf("[UserService.GetWallet] user: %s, currency: %s",
		service.ctx.GetUser().GetUsername(), currency)
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	walletEntity := service.userDAO.GetWallet(daoUser, currency)
	wallet, err := service.walletService.CreateWallet(walletEntity.GetCurrency(), walletEntity.GetAddress())
	if err != nil {
		service.ctx.GetLogger().Errorf("[UserService.GetWallet] Error: %s", err.Error())
		return nil, err
	}
	return wallet.GetWallet()
}

func (service *DefaultUserService) GetWallets() []common.UserCryptoWallet {
	var walletList []common.UserCryptoWallet
	daoUser := &entity.User{Id: service.ctx.GetUser().GetId()}
	walletEntities := service.userDAO.GetWallets(daoUser)
	var chans []chan common.UserCryptoWallet
	for _, walletEntity := range walletEntities {
		entity := walletEntity
		c := make(chan common.UserCryptoWallet, 1)
		chans = append(chans, c)
		go func() {
			wallet, err := service.walletService.CreateWallet(entity.Currency, entity.Address)
			if err != nil {
				service.ctx.GetLogger().Errorf("[UserService.GetWallets] Unable to create %s wallet instance: %s",
					entity.Currency, err.Error())
			}
			userWallet, err := wallet.GetWallet()
			if err != nil {
				service.ctx.GetLogger().Errorf("[UserService.GetWallets] Unable to retrieve %s's %s wallet: %s",
					service.ctx.GetUser().GetUsername(), entity.Currency, err.Error())
			}
			c <- userWallet
		}()
	}
	for i := 0; i < len(walletEntities); i++ {
		walletList = append(walletList, <-chans[i])
	}
	return walletList
}
