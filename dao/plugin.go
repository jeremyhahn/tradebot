package dao

import (
	"errors"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type PluginDAO interface {
	Create(plugin entity.PluginEntity) error
	Save(plugin entity.PluginEntity) error
	Update(plugin entity.PluginEntity) error
	Find(pluginType string) ([]entity.Plugin, error)
	Get(pluginName, pluginType string) (entity.PluginEntity, error)
}

type PluginDAOImpl struct {
	ctx common.Context
	PluginDAO
}

func NewPluginDAO(ctx common.Context) PluginDAO {
	return &PluginDAOImpl{ctx: ctx}
}

func (dao *PluginDAOImpl) Create(indicator entity.PluginEntity) error {
	return dao.ctx.GetCoreDB().Create(indicator).Error
}

func (dao *PluginDAOImpl) Save(indicator entity.PluginEntity) error {
	return dao.ctx.GetCoreDB().Save(indicator).Error
}

func (dao *PluginDAOImpl) Update(indicator entity.PluginEntity) error {
	return dao.ctx.GetCoreDB().Update(indicator).Error
}

func (dao *PluginDAOImpl) Get(pluginName, pluginType string) (entity.PluginEntity, error) {
	var plugins []entity.Plugin
	if err := dao.ctx.GetCoreDB().Where("name = ? AND type = ?", pluginName, pluginType).Find(&plugins).Error; err != nil {
		return nil, err
	}
	if len(plugins) == 0 {
		return nil, errors.New(fmt.Sprintf("%s (%s) plugin not found in database", pluginName, pluginType))
	}
	return &plugins[0], nil
}

func (dao *PluginDAOImpl) Find(pluginType string) ([]entity.Plugin, error) {
	var plugins []entity.Plugin
	if err := dao.ctx.GetCoreDB().Where("type = ?", pluginType).Order("name asc").Find(&plugins).Error; err != nil {
		return nil, err
	}
	if len(plugins) == 0 {
		return nil, errors.New("no plugins found in database")
	}
	return plugins, nil
}
