package dto

type PlatformIndicator interface {
	GetName() string
	GetFilename() string
	GetVersion() string
}

type PlatformIndicatorDTO struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	PlatformIndicator
}

func NewPlatformIndicatorDTO() PlatformIndicator {
	return &PlatformIndicatorDTO{}
}

func CreatePlatformIndicatorDTO(name, filename, version string) PlatformIndicator {
	return &PlatformIndicatorDTO{
		Name:     name,
		Filename: filename,
		Version:  version}
}

func (dto *PlatformIndicatorDTO) GetName() string {
	return dto.Name
}

func (dto *PlatformIndicatorDTO) GetFilename() string {
	return dto.Filename
}

func (dto *PlatformIndicatorDTO) GetVersion() string {
	return dto.Version
}
