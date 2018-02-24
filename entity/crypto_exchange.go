package entity

type CryptoExchange struct {
	Name     string `gorm:"primary_key"`
	Filename string `gorm:"not null"`
	Version  string `gorm:"not null"`
}

func (entity *CryptoExchange) GetName() string {
	return entity.Name
}

func (entity *CryptoExchange) GetFilename() string {
	return entity.Filename
}

func (entity *CryptoExchange) GetVersion() string {
	return entity.Version
}
