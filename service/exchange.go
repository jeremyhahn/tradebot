package service

import (
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultExchangeService struct {
	ctx                 common.Context
	pluginDAO           dao.PluginDAO
	userDAO             dao.UserDAO
	userMapper          mapper.UserMapper
	userExchangeMapper  mapper.UserExchangeMapper
	priceHistoryService common.PriceHistoryService
	exchangeMap         map[string]func(common.Context, entity.UserExchangeEntity, common.PriceHistoryService) common.Exchange
	ExchangeService
}

func NewExchangeService(ctx common.Context, pluginDAO dao.PluginDAO, userDAO dao.UserDAO,
	userMapper mapper.UserMapper, userExchangeMapper mapper.UserExchangeMapper,
	priceHistoryService common.PriceHistoryService) ExchangeService {
	return &DefaultExchangeService{
		ctx:                 ctx,
		pluginDAO:           pluginDAO,
		userDAO:             userDAO,
		userMapper:          userMapper,
		userExchangeMapper:  userExchangeMapper,
		priceHistoryService: priceHistoryService,
		exchangeMap: map[string]func(ctx common.Context, exchange entity.UserExchangeEntity,
			priceHistoryService common.PriceHistoryService) common.Exchange{
			"gdax":    exchange.NewGDAX,
			"bittrex": exchange.NewBittrex,
			"binance": exchange.NewBinance}}
}

func (service *DefaultExchangeService) CreateExchange(exchangeName string) (common.Exchange, error) {
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	exchange, err := service.userDAO.GetExchange(userEntity, exchangeName)
	if err != nil {
		service.ctx.GetLogger().Errorf("[ExchangeService.CreateExchange] Error: %s", err.Error())
		return nil, err
	}
	return service.exchangeMap[exchangeName](service.ctx, exchange, service.priceHistoryService), nil
}

func (service *DefaultExchangeService) GetDisplayNames() []string {
	var exchanges []string
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, ex.Name)
	}
	return exchanges
}

func (service *DefaultExchangeService) GetExchanges() []common.Exchange {
	var exchanges []common.Exchange
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	userExchanges := service.userDAO.GetExchanges(userEntity)
	for _, ex := range userExchanges {
		exchanges = append(exchanges, service.exchangeMap[ex.Name](service.ctx, &ex, service.priceHistoryService))
	}
	return exchanges
}

func (service *DefaultExchangeService) GetExchange(exchangeName string) common.Exchange {
	for _, ex := range service.GetExchanges() {
		if ex.GetName() == exchangeName {
			return ex
		}
	}
	return nil
}

func (service *DefaultExchangeService) GetCurrencyPairs(exchangeName string) ([]common.CurrencyPair, error) {
	userEntity := service.userMapper.MapUserDtoToEntity(service.ctx.GetUser())
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
