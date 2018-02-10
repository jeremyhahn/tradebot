package mapper

import (
	"testing"

	"github.com/jeremyhahn/tradebot/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserMapper(t *testing.T) {
	mapper := NewUserMapper()

	userDTO := &dto.UserDTO{
		Id:            1,
		Username:      "Test",
		LocalCurrency: "USD",
		Etherbase:     "0xabc123",
		Keystore:      "/tmp/keystore"}

	entity := mapper.MapUserDtoToEntity(userDTO)
	assert.NotNil(t, entity)
	assert.Equal(t, userDTO.GetId(), entity.GetId())
	assert.Equal(t, userDTO.GetUsername(), entity.GetUsername())
	assert.Equal(t, userDTO.GetLocalCurrency(), entity.GetLocalCurrency())
	assert.Equal(t, userDTO.GetEtherbase(), entity.GetEtherbase())
	assert.Equal(t, userDTO.GetKeystore(), entity.GetKeystore())

	dto := mapper.MapUserEntityToDto(entity)
	assert.Equal(t, entity.GetId(), dto.GetId())
	assert.Equal(t, entity.GetUsername(), dto.GetUsername())
	assert.Equal(t, entity.GetLocalCurrency(), dto.GetLocalCurrency())
	assert.Equal(t, entity.GetEtherbase(), dto.GetEtherbase())
	assert.Equal(t, entity.GetKeystore(), dto.GetKeystore())
}
