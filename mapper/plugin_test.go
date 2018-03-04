package mapper

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dto"
	"github.com/stretchr/testify/assert"
)

func TestPluginMapper_Indicators(t *testing.T) {
	mapper := NewPluginMapper()
	dto := &dto.PluginDTO{
		Name:     "TestIndicator",
		Filename: "test.so",
		Version:  "0.0.1a",
		Type:     "indicator"}

	entity := mapper.MapPluginDtoToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, "TestIndicator", entity.GetName())
	assert.Equal(t, "test.so", entity.GetFilename())
	assert.Equal(t, "0.0.1a", entity.GetVersion())
	assert.Equal(t, "indicator", entity.GetType())

	mappedDTO := mapper.MapPluginEntityToDto(entity)
	assert.NotNil(t, entity)
	assert.Equal(t, mappedDTO.GetName(), entity.GetName())
	assert.Equal(t, mappedDTO.GetFilename(), entity.GetFilename())
	assert.Equal(t, mappedDTO.GetVersion(), entity.GetVersion())
	assert.Equal(t, mappedDTO.GetType(), entity.GetType())
}

func TestStrategyMapper_Strategy(t *testing.T) {
	mapper := NewPluginMapper()
	dto := &dto.PluginDTO{
		Name:     "TestStrategy",
		Filename: "test.so",
		Version:  "0.0.1a",
		Type:     "strategy"}

	entity := mapper.MapPluginDtoToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, "TestStrategy", entity.GetName())
	assert.Equal(t, "test.so", entity.GetFilename())
	assert.Equal(t, "0.0.1a", entity.GetVersion())
	assert.Equal(t, "strategy", entity.GetType())

	mappedDTO := mapper.MapPluginEntityToDto(entity)
	assert.NotNil(t, entity)
	assert.Equal(t, mappedDTO.GetName(), entity.GetName())
	assert.Equal(t, mappedDTO.GetFilename(), entity.GetFilename())
	assert.Equal(t, mappedDTO.GetVersion(), entity.GetVersion())
	assert.Equal(t, mappedDTO.GetType(), entity.GetType())
}
