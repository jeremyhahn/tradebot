package dto

type PlatformStrategy interface {
	GetName() string
	GetFilename() string
	GetVersion() string
}

type PlatformStrategyDTO struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	PlatformStrategy
}

func NewPlatformStrategyDTO() PlatformStrategy {
	return &PlatformStrategyDTO{}
}

func CreatePlatformStrategyDTO(name, filename, version string) PlatformStrategy {
	return &PlatformStrategyDTO{
		Name:     name,
		Filename: filename,
		Version:  version}
}

func (dto *PlatformStrategyDTO) GetName() string {
	return dto.Name
}

func (dto *PlatformStrategyDTO) GetFilename() string {
	return dto.Filename
}

func (dto *PlatformStrategyDTO) GetVersion() string {
	return dto.Version
}
