package entity

type Strategy struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
}

func (entity *Strategy) GetName() string {
	return entity.Name
}

func (entity *Strategy) GetFilename() string {
	return entity.Filename
}

func (entity *Strategy) GetVersion() string {
	return entity.Version
}
