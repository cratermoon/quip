package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {

	assert := assert.New(t)

	kit, err := NewKit()
	assert.Nil(err)

	q, err := kit.SDBTakeFirst("text", "newquips")
	assert.Nil(err)
	assert.NotNil(q)

	fmt.Println(q)
}
