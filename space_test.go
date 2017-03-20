package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaceSaveForCreate(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("new space", payload["name"])
		assert.Equal("en", payload["defaultLocale"])

		w.WriteHeader(201)
		fmt.Fprintln(w, string(readTestData("spaces-newspace.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space := cma.NewSpace()
	space.Name = "new space"
	space.DefaultLocale = "en"
	err := space.Save()
	assert.Nil(err)
	assert.Equal("newspace", space.Sys.ID)
	assert.Equal("new space", space.Name)
	assert.Equal("en", space.DefaultLocale)
}

func TestSpaceSaveForUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/newspace")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("changed-space-name", payload["name"])
		assert.Equal("de", payload["defaultLocale"])

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("spaces-newspace-updated.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := spaceFromTestData("spaces-newspace.json")
	assert.Nil(err)

	space.Name = "changed-space-name"
	space.DefaultLocale = "de"
	err = space.Save()
	assert.Nil(err)
	assert.Equal("changed-space-name", space.Name)
	assert.Equal("de", space.DefaultLocale)
	assert.Equal(2, space.Sys.Version)
}

func TestSpaceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID)
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := spaceFromTestData("spaces-" + spaceID + ".json")
	assert.Nil(err)

	err = space.Delete()
	assert.Nil(err)
}
