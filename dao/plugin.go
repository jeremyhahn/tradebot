package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type PluginDAO interface {
	Create(plugin entity.PluginEntity) error
	Save(plugin entity.PluginEntity) error
	Update(plugin entity.PluginEntity) error
	Find() ([]entity.Plugin, error)
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
	if err := dao.ctx.GetCoreDB().Where("name = ?, type = ?", pluginName, pluginType).Model(&entity.Plugin{}).Error; err != nil {
		return nil, err
	}
	return &plugins[0], nil
}

func (dao *PluginDAOImpl) Find() ([]entity.Plugin, error) {
	var plugins []entity.Plugin
	if err := dao.ctx.GetCoreDB().Order("name asc").Model(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}
