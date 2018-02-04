package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/jeremyhahn/tradebot/common"
)

const (
	INDICATOR_PLUGIN = "indicators"
	STRATEGY_PLUGIN  = "strategies"
)

type PluginService interface {
	GetIndicator(pluginName, indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error)
	GetStrategy(pluginName, strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error)
}

type PluginServiceImpl struct {
	ctx        *common.Context
	pluginRoot string
	PluginService
}

func NewPluginService(ctx *common.Context) PluginService {
	return &PluginServiceImpl{
		ctx:        ctx,
		pluginRoot: "./plugins"}
}

func CreatePluginService(ctx *common.Context, pluginRoot string) PluginService {
	return &PluginServiceImpl{
		ctx:        ctx,
		pluginRoot: pluginRoot}
}

func (p *PluginServiceImpl) GetIndicator(pluginName, indicatorName string) (func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error), error) {
	lib, err := p.openPlugin(INDICATOR_PLUGIN, pluginName)
	if err != nil {
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetIndicator] Error loading %s. %s", pluginName, err.Error())
		return nil, err
	}
	symbolName := strings.Split(indicatorName, ".")
	symbol := fmt.Sprintf("Create%s", symbolName[0])
	p.ctx.Logger.Debugf("[PluginServiceImpl.GetIndicator] Looking up symbol %s", symbol)
	indicator, err := lib.Lookup(symbol)
	impl, ok := indicator.(func(candles []common.Candlestick, params []string) (common.FinancialIndicator, error))
	if !ok {
		errmsg := fmt.Sprintf("Wrong type - expected (%s) %s(candles []common.Candlestick, params []string) (common.FinancialIndicator, error)", pluginName, symbol)
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetIndicator] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (p *PluginServiceImpl) GetStrategy(pluginName, strategyName string) (func(params *common.TradingStrategyParams) (common.TradingStrategy, error), error) {
	lib, err := p.openPlugin(STRATEGY_PLUGIN, pluginName)
	if err != nil {
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetStrategy] Error loading %s. %s", pluginName, err.Error())
		return nil, err
	}
	symbolName := strings.Split(strategyName, ".")
	symbol := fmt.Sprintf("Create%s", symbolName[0])
	p.ctx.Logger.Debugf("[PluginServiceImpl.GetStrategy] Looking up symbol %s", symbol)
	strategy, err := lib.Lookup(symbol)
	if err != nil {
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetStrategy] %s", err.Error())
		return nil, err
	}
	impl, ok := strategy.(func(params *common.TradingStrategyParams) (common.TradingStrategy, error))
	if !ok {
		errmsg := fmt.Sprintf("Wrong type - expected (%s) %s(params *common.TradingStrategyParams) (common.TradingStrategy, error))", pluginName, symbol)
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetStrategy] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}

func (p *PluginServiceImpl) openPlugin(which, name string) (*plugin.Plugin, error) {
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s", p.pluginRoot, which, name))
	p.ctx.Logger.Debugf("[PluginServiceImpl.openPlugin] Loading plugin %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	return plugin.Open(path)
}
