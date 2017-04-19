package gcloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDBAdd(t *testing.T) {

	assert := assert.New(t)

	ctx, done, err := aetest.NewContext()
	assert.Nil(err)
	defer done()
	fmt.Printf("I haz a context %+v\n", ctx)
	kit, err := NewKit()
	assert.Nil(err)

	q, err := kit.DBAdd("text", "Don't take any wooden nickels", "newquips")
	assert.Nil(err)
	assert.NotNil(q)

	fmt.Println(q)
}
