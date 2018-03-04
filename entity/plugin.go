package entity

type Plugin struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
	Type     string `gorm:"not null"`
	PluginEntity
}

func (entity *Plugin) GetName() string {
	return entity.Name
}

func (entity *Plugin) GetFilename() string {
	return entity.Filename
}

func (entity *Plugin) GetVersion() string {
	return entity.Version
}

func (entity *Plugin) GetType() string {
	return entity.Type
}
