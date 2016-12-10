package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveLocale(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	locale := space.NewLocale()
	locale.Name = "Deutsch"
	locale.Code = "de"
	locale.FallbackCode = "de"

	req, err := locale.getSaveReq()
	assert.Nil(err)
	assert.Equal("POST", req.Method)
	assert.Equal("/spaces/"+spaceID+"/locales", req.URL.Path)
	assert.Equal("0", req.Header.Get("X-Contentful-Version"))
}
