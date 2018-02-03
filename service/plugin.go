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
	GetIndicator(pluginName string) (func(candles []common.Candlestick, params []string) common.FinancialIndicator, error)
}

type PluginServiceImpl struct {
	ctx       *common.Context
	directory string
	PluginService
}

func NewPluginService(ctx *common.Context, pluginType string) PluginService {
	return &PluginServiceImpl{
		ctx:       ctx,
		directory: fmt.Sprintf("./%s", pluginType)}
}

func CreatePluginService(ctx *common.Context, pluginRoot, pluginType string) PluginService {
	return &PluginServiceImpl{
		ctx:       ctx,
		directory: fmt.Sprintf("%s/%s", pluginRoot, pluginType)}
}

func (p *PluginServiceImpl) GetIndicator(pluginName string) (func(candles []common.Candlestick,
	params []string) common.FinancialIndicator, error) {

	path, _ := filepath.Abs(fmt.Sprintf("%s/%s", p.directory, pluginName))
	p.ctx.Logger.Debugf("[PluginServiceImpl.GetIndicator] Loading plugin %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New(err.Error())
	}
	lib, err := plugin.Open(path)
	if err != nil {
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetIndicator] Error loading %s. %s", pluginName, err.Error())
		return nil, err
	}
	symbolName := strings.Split(pluginName, ".")
	symbol := fmt.Sprintf("Create%s", symbolName[0])
	p.ctx.Logger.Debugf("[PluginServiceImpl.GetIndicator] Looking up symbol %s", symbol)
	indicator, err := lib.Lookup(symbol)
	impl, ok := indicator.(func(candles []common.Candlestick, params []string) common.FinancialIndicator)
	if !ok {
		errmsg := fmt.Sprintf("Wrong type - expected func - %s", pluginName)
		p.ctx.Logger.Errorf("[PluginServiceImpl.GetIndicator] %s", errmsg)
		return nil, errors.New(errmsg)
	}
	return impl, nil
}
