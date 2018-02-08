package mapper

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dto"
	"github.com/stretchr/testify/assert"
)

func TestIndicatorMapper(t *testing.T) {
	mapper := NewIndicatorMapper()
	dto := &dto.IndicatorDTO{
		Name:     "TestIndicator",
		Filename: "test.so",
		Version:  "0.0.1a"}

	entity := mapper.MapIndicatorDtoToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, "TestIndicator", entity.GetName())
	assert.Equal(t, "test.so", entity.GetFilename())
	assert.Equal(t, "0.0.1a", entity.GetVersion())

	mappedDTO := mapper.MapIndicatorEntityToDto(entity)
	assert.NotNil(t, entity)
	assert.Equal(t, mappedDTO.GetName(), entity.GetName())
	assert.Equal(t, mappedDTO.GetFilename(), entity.GetFilename())
	assert.Equal(t, mappedDTO.GetVersion(), entity.GetVersion())
}
