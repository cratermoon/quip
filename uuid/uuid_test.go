package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {

	assert := assert.New(t)

	uuid, err := NewUUID()
	assert.Nil(err)
	assert.NotNil(uuid)
	assert.Equal(36, len(uuid), "Wrong length for uuid")
}
