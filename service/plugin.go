package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
)

const (
	INDICATOR_PLUGIN = "indicators"
	STRATEGY_PLUGIN  = "strategies"
	INDICATOR_TYPE   = "indicator"
	STRATEGY_TYPE    = "strategy"
)

var PLUGINS = make(map[string]*plugin.Plugin)

type PluginService interface {
	GetMapper() mapper.PluginMapper
	GetPlugin(pluginName, pluginType string) (common.Plugin, error)
	CreateIndicator(indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error)
	CreateStrategy(strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error)
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

func (service *DefaultPluginService) CreateIndicator(indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error) {
	indicatorEntity, err := service.dao.Get(indicatorName, INDICATOR_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] Error: %s", err.Error())
		return nil, err
	}
	filename := indicatorEntity.GetFilename()
	lib, err := service.openPlugin(INDICATOR_PLUGIN, filename)
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
		errmsg := fmt.Sprintf("Invalid type - expected (%s) %s(candles []common.Candlestick, params []string) (common.FinancialIndicator, error)", filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateIndicator] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (service *DefaultPluginService) CreateStrategy(strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error) {
	strategyEntity, err := service.dao.Get(strategyName, STRATEGY_TYPE)
	if err != nil {
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] Error: %s", err.Error())
		return nil, err
	}
	filename := strategyEntity.GetFilename()
	lib, err := service.openPlugin(STRATEGY_PLUGIN, filename)
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
		errmsg := fmt.Sprintf("Invalid type - expected (%s) %s(params *common.TradingStrategyParams) (common.TradingStrategy, error))", filename, symbol)
		service.ctx.GetLogger().Errorf("[PluginService.CreateStrategy] %s", errmsg)
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
