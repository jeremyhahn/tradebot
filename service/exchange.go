package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type DefaultExchangeService struct {
	ctx            common.Context
	pluginDAO      dao.PluginDAO
	userDAO        dao.UserDAO
	userMapper     mapper.UserMapper
	exchangeMapper mapper.UserExchangeMapper
	exchangeMap    map[string]func(common.Context, entity.UserExchangeEntity) common.Exchange
	ExchangeService
}

func NewExchangeService(ctx common.Context, pluginDAO dao.PluginDAO, userDAO dao.UserDAO,
	userMapper mapper.UserMapper, exchangeMapper mapper.UserExchangeMapper) ExchangeService {
	return &DefaultExchangeService{
		ctx:            ctx,
		pluginDAO:      pluginDAO,
		userDAO:        userDAO,
		userMapper:     userMapper,
		exchangeMapper: exchangeMapper,
		exchangeMap: map[string]func(ctx common.Context, exchange entity.UserExchangeEntity) common.Exchange{
			"gdax":    exchange.NewGDAX,
			"bittrex": exchange.NewBittrex,
			"binance": exchange.NewBinance}}
}

func (service *DefaultExchangeService) CreateExchange(user common.UserContext, exchangeName string) (common.Exchange, error) {
	userEntity := &entity.User{Id: user.GetId()}
	exchange, err := service.userDAO.GetExchange(userEntity, exchangeName)
	if err != nil {
		service.ctx.GetLogger().Errorf("[ExchangeService.CreateExchange] Error: %s", err.Error())
		return nil, err
	}
	return service.exchangeMap[exchangeName](service.ctx, exchange), nil
}

func (service *DefaultExchangeService) GetDisplayNames(user common.UserContext) []string {
	var exchanges []string
	userEntity := &entity.User{Id: user.GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, ex.Name)
	}
	return exchanges
}

func (service *DefaultExchangeService) GetUserExchanges(user common.UserContext) []viewmodel.UserCryptoExchange {
	var exchanges []viewmodel.UserCryptoExchange
	userEntity := &entity.User{Id: user.GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		viewmodel := service.exchangeMapper.MapEntityToViewModel(&ex)
		exchanges = append(exchanges, *viewmodel)
	}
	return exchanges
}

func (service *DefaultExchangeService) GetExchanges(user common.UserContext) []common.Exchange {
	var exchanges []common.Exchange
	userEntity := &entity.User{Id: user.GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, service.exchangeMap[ex.Name](service.ctx, &ex))
	}
	return exchanges
}

func (service *DefaultExchangeService) GetExchange(user common.UserContext, exchangeName string) common.Exchange {
	for _, ex := range service.GetExchanges(user) {
		if ex.GetName() == exchangeName {
			return ex
		}
	}
	return nil
}

func (service *DefaultExchangeService) GetCurrencyPairs(user common.UserContext, exchangeName string) ([]common.CurrencyPair, error) {
	userEntity := service.userMapper.MapUserDtoToEntity(user)
	userCryptoExchange, err := service.userDAO.GetExchange(userEntity.(*entity.User), exchangeName)
	if err != nil {
		return nil, err
	}
	return service.parseCurrencyPairs(userCryptoExchange.GetExtra(), exchangeName), nil
}

func (service *DefaultExchangeService) parseCurrencyPairs(configuredPairs, exchangeName string) []common.CurrencyPair {
	var currencyPairs []common.CurrencyPair
	pairs := strings.Split(configuredPairs, ",")
	for _, pair := range pairs {
		pieces := strings.Split(pair, "-")
		if len(pieces) != 2 {
			service.ctx.GetLogger().Errorf("[DefaultExchangeService.parseCurrencyPairs] Invalid currency pair configured for %s: %+v", exchangeName, pair)
			continue
		}
		currencyPairs = append(currencyPairs, common.CurrencyPair{
			Base:          pieces[0],
			Quote:         pieces[1],
			LocalCurrency: service.ctx.GetUser().GetLocalCurrency()})
	}
	return currencyPairs
}
