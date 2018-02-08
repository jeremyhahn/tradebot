package mapper

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dto"
	"github.com/stretchr/testify/assert"
)

func TestStrategyMapper(t *testing.T) {
	mapper := NewStrategyMapper()
	dto := &dto.StrategyDTO{
		Name:     "TestStrategy",
		Filename: "test.so",
		Version:  "0.0.1a"}

	entity := mapper.MapStrategyDtoToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, "TestStrategy", entity.GetName())
	assert.Equal(t, "test.so", entity.GetFilename())
	assert.Equal(t, "0.0.1a", entity.GetVersion())

	mappedDTO := mapper.MapStrategyEntityToDto(entity)
	assert.NotNil(t, entity)
	assert.Equal(t, mappedDTO.GetName(), entity.GetName())
	assert.Equal(t, mappedDTO.GetFilename(), entity.GetFilename())
	assert.Equal(t, mappedDTO.GetVersion(), entity.GetVersion())
}
