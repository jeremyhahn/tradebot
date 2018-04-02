package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

var PLUGINS = make(map[string]*plugin.Plugin)
var PLUGINTYPE = map[string]string{
	common.INDICATOR_PLUGIN_TYPE: "indicators",
	common.STRATEGY_PLUGIN_TYPE:  "strategies",
	common.EXCHANGE_PLUGIN_TYPE:  "exchanges",
	common.WALLET_PLUGIN_TYPE:    "wallets"}

type PluginService interface {
	GetMapper() mapper.PluginMapper
	GetPlugin(pluginName, pluginType string) (common.Plugin, error)
	GetPlugins(pluginType string) ([]string, error)
	ListPlugins(pluginType string) ([]string, error)
	CreateIndicator(indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error)
	CreateStrategy(strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error)
	CreateExchange(exchangeName string) (func(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange, error)
	CreateWallet(currency string) (func(params *common.WalletParams) common.Wallet, error)
}

type DefaultPluginService struct {
	ctx        common.Context
	pluginRoot string
	loaded     map[string]*plugin.Plugin
	dao        dao.PluginDAO
	mapper     mapper.PluginMapper
	PluginService
}

func NewPluginService(ctx common.Context, pluginDAO dao.PluginDAO, pluginMapper mapper.PluginMapper) PluginService {
	return &DefaultPluginService{
		ctx:        ctx,
		pluginRoot: "./plugins",
		dao:        pluginDAO,
		mapper:     pluginMapper}
}

func CreatePluginService(ctx common.Context, pluginRoot string, pluginDAO dao.PluginDAO,
	pluginMapper mapper.PluginMapper) PluginService {

	return &DefaultPluginService{
		ctx:        ctx,
		pluginRoot: pluginRoot,
		loaded:     make(map[string]*plugin.Plugin),
		dao:        pluginDAO,
		mapper:     pluginMapper}
}

func (service *DefaultPluginService) GetMapper() mapper.PluginMapper {
	return service.mapper
}

func (service *DefaultPluginService) GetPlugin(pluginName, pluginType string) (common.Plugin, error) {
	entity, err := service.dao.Get(pluginName, pluginType)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.GetPlugin] Error: %s", err.Error())
		return nil, err
	}
	return service.mapper.MapPluginEntityToDto(entity), nil
}

func (service *DefaultPluginService) GetPlugins(pluginType string) ([]string, error) {
	entities, err := service.dao.Find(pluginType)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.GetPlugin] Error: %s", err.Error())
		return nil, err
	}
	plugins := make([]string, len(entities))
	for i, entity := range entities {
		plugins[i] = entity.GetName()
	}
	return plugins, nil
}

func (service *DefaultPluginService) CreateIndicator(indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error) {
	indicatorEntity, err := service.dao.Get(indicatorName, common.INDICATOR_PLUGIN_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] Error: %s", err.Error())
		return nil, err
	}
	filename := indicatorEntity.GetFilename()
	lib, err := service.openPlugin(PLUGINTYPE[common.INDICATOR_PLUGIN_TYPE], filename)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] Error loading %s. %s", filename, err.Error())
		return nil, err
	}
	symbolName := strings.Split(indicatorName, ".")
	symbol := fmt.Sprintf("Create%s", symbolName[0])
	service.ctx.GetLogger().Debugf("[PluginService.CreateIndicator] Looking up indicator symbol %s", symbol)
	indicator, err := lib.Lookup(symbol)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] Error loading symbol %s.%s. %s", filename, symbol, err.Error())
		return nil, err
	}
	impl, ok := indicator.(func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error))
	if !ok {
		errmsg := fmt.Sprintf("Invalid plugin, expected factory method: (%s) %s(candles []common.Candlestick, params []string) (common.FinancialIndicator, error)",
			filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (service *DefaultPluginService) CreateStrategy(strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error) {
	strategyEntity, err := service.dao.Get(strategyName, common.STRATEGY_PLUGIN_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] Error: %s", err.Error())
		return nil, err
	}
	filename := strategyEntity.GetFilename()
	lib, err := service.openPlugin(PLUGINTYPE[common.STRATEGY_PLUGIN_TYPE], filename)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] Error loading %s. %s", filename, err.Error())
		return nil, err
	}
	symbolName := strings.Split(strategyName, ".")
	symbol := fmt.Sprintf("Create%s", symbolName[0])
	service.ctx.GetLogger().Debugf("[PluginService.CreateStrategy] Looking up strategy symbol %s", symbol)
	strategy, err := lib.Lookup(symbol)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] %s", err.Error())
		return nil, err
	}
	impl, ok := strategy.(func(params *common.TradingStrategyParams) (common.TradingStrategy, error))
	if !ok {
		errmsg := fmt.Sprintf("Invalid plugin, expected factory method: (%s) %s(params *common.TradingStrategyParams) (common.TradingStrategy, error))",
			filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (service *DefaultPluginService) CreateExchange(exchangeName string) (func(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange, error) {
	exchangeEntity, err := service.dao.Get(exchangeName, common.EXCHANGE_PLUGIN_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateExchange] Error loading exchange from database: %s", err.Error())
		return nil, err
	}
	filename := exchangeEntity.GetFilename()
	lib, err := service.openPlugin(PLUGINTYPE[common.EXCHANGE_PLUGIN_TYPE], filename)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateExchange] Error loading %s. %s", filename, err.Error())
		return nil, err
	}
	symbol := fmt.Sprintf("Create%s", exchangeName)
	service.ctx.GetLogger().Debugf("[PluginService.CreateExchange] Looking up exchange symbol %s", symbol)
	exchange, err := lib.Lookup(symbol)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateExchange] Error looking up exchange symbol: %s", err.Error())
		return nil, err
	}
	impl, ok := exchange.(func(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange)
	if !ok {
		errmsg := fmt.Sprintf("Invalid plugin, expected factory method: (%s) %s(params *common.TradingStrategyParams) (common.TradingStrategy, error))",
			filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateExchange] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (service *DefaultPluginService) CreateWallet(currency string) (func(params *common.WalletParams) common.Wallet, error) {
	walletEntity, err := service.dao.Get(currency, common.WALLET_PLUGIN_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateWallet] Error loading wallet from database: %s", err.Error())
		return nil, err
	}
	filename := walletEntity.GetFilename()
	lib, err := service.openPlugin(PLUGINTYPE[common.WALLET_PLUGIN_TYPE], filename)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateWallet] Error loading %s. %s", filename, err.Error())
		return nil, err
	}
	walletName := strings.Title(strings.ToLower(currency))
	symbol := fmt.Sprintf("Create%sWallet", walletName)
	service.ctx.GetLogger().Debugf("[PluginService.CreateWallet] Looking up wallet symbol %s", symbol)
	wallet, err := lib.Lookup(symbol)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateWallet] Error looking up wallet symbol: %s", err.Error())
		return nil, err
	}
	impl, ok := wallet.(func(params *common.WalletParams) common.Wallet)
	if !ok {
		errmsg := fmt.Sprintf("Invalid plugin, expected factory method: (%s) %s(params *common.WalletParams) common.Wallet)",
			filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateWallet] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (service *DefaultPluginService) openPlugin(which, name string) (*plugin.Plugin, error) {
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s", service.pluginRoot, which, name))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	if indicator, ok := PLUGINS[name]; ok {
		return indicator, nil
	} else {
		service.ctx.GetLogger().Debugf("[PluginService.openPlugin] Loading plugin %s", path)
		_lib, err := plugin.Open(path)
		if err != nil {
			service.ctx.GetLogger().Errorf("[PluginService.openPlugin] Error loading plugin %s. %s", name, err.Error())
			return nil, err
		}
		PLUGINS[name] = _lib
		return PLUGINS[name], nil
	}
}

func (service *DefaultPluginService) ListPlugins(pluginType string) ([]string, error) {
	var plugins []string
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s", service.pluginRoot, PLUGINTYPE[pluginType]))
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		pieces := strings.Split(f.Name(), ".")
		plugins = append(plugins, pieces[0])
	}
	return plugins, nil
}
