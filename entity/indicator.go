package entity

type Indicator struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
}

func (entity *Indicator) GetName() string {
	return entity.Name
}

func (entity *Indicator) GetFilename() string {
	return entity.Filename
}

func (entity *Indicator) GetVersion() string {
	return entity.Version
}
