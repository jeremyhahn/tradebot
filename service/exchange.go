package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type DefaultExchangeService struct {
	ctx                common.Context
	userDAO            dao.UserDAO
	userMapper         mapper.UserMapper
	userExchangeMapper mapper.UserExchangeMapper
	pluginService      PluginService
	ExchangeService
}

func NewExchangeService(ctx common.Context, userDAO dao.UserDAO, userMapper mapper.UserMapper,
	userExchangeMapper mapper.UserExchangeMapper, pluginService PluginService) ExchangeService {
	return &DefaultExchangeService{
		ctx:                ctx,
		userDAO:            userDAO,
		userMapper:         userMapper,
		userExchangeMapper: userExchangeMapper,
		pluginService:      pluginService}
}

func (service *DefaultExchangeService) CreateExchange(exchangeName string) (common.Exchange, error) {
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	userCryptoExchange, err := service.userDAO.GetExchange(userEntity, exchangeName)
	if err != nil {
		service.ctx.GetLogger().Errorf("[ExchangeService.CreateExchange] Error: %s", err.Error())
		return nil, err
	}
	exchange, err := service.pluginService.CreateExchange(exchangeName)
	return exchange(service.ctx, userCryptoExchange), nil
}

func (service *DefaultExchangeService) GetDisplayNames() ([]string, error) {
	names, err := service.pluginService.GetPlugins(common.EXCHANGE_PLUGIN_TYPE)
	if err != nil {
		return names, err
	}
	return names, nil
}

func (service *DefaultExchangeService) GetExchanges() ([]common.Exchange, error) {
	var exchanges []common.Exchange
	userEntity := &entity.User{Id: service.ctx.GetUser().GetId()}
	userCryptoExchanges := service.userDAO.GetExchanges(userEntity)
	for _, userCryptExchange := range userCryptoExchanges {
		exchange, err := service.pluginService.CreateExchange(userCryptExchange.GetName())
		if err != nil {
			service.ctx.GetLogger().Errorf("[ExchangeService.GetExchanges] Error: %s", err.Error())
			return nil, err
		}
		exchanges = append(exchanges, exchange(service.ctx, &userCryptExchange))
	}
	return exchanges, nil
}

func (service *DefaultExchangeService) GetExchange(exchangeName string) (common.Exchange, error) {
	exchanges, err := service.GetExchanges()
	if err != nil {
		service.ctx.GetLogger().Errorf("[ExchangeService.GetExchange] Error: %s", err.Error())
		return nil, err
	}
	for _, ex := range exchanges {
		if ex.GetName() == exchangeName {
			return ex, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to locate exchange: %s", exchangeName))
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
