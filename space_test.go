package contentful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSpacesService_Get() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(space.Name)
}

func ExampleSpacesService_List() {
	cma := NewCMA("cma-token")
	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	spaces := collection.ToSpace()
	for _, space := range spaces {
		fmt.Println(space.Sys.ID, space.Name)
	}
}

func ExampleSpacesService_Upsert_create() {
	cma := NewCMA("cma-token")

	space := &Space{
		Name:          "space-name",
		DefaultLocale: "en-US",
	}

	err := cma.Spaces.Upsert(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Upsert_update() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	space.Name = "modified"
	err = cma.Spaces.Upsert(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete() {
	cma := NewCMA("cma-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Spaces.Delete(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete_all() {
	cma := NewCMA("cma-token")

	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, space := range collection.ToSpace() {
		err := cma.Spaces.Delete(space)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestSpacesServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Spaces.List().Next()
	assert.Nil(err)

	spaces := collection.ToSpace()
	assert.Equal(2, len(spaces))
	assert.Equal("id1", spaces[0].Sys.ID)
	assert.Equal("id2", spaces[1].Sys.ID)
}

func TestSpacesServiceList_Pagination(t *testing.T) {
	var err error
	assert := assert.New(t)

	requestCount := 1
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces")
		checkHeaders(r, assert)

		w.WriteHeader(200)
		query := r.URL.Query()
		if requestCount == 1 {
			assert.Equal(query.Get("order"), "-sys.createdAt")
			assert.Equal(query.Get("skip"), "")
			fmt.Fprintln(w, readTestData("spaces.json"))
		} else {
			assert.Equal(query.Get("order"), "-sys.createdAt")
			assert.Equal(query.Get("skip"), "100")
			fmt.Fprintln(w, readTestData("spaces-page-2.json"))
		}
		requestCount++
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Spaces.List().Next()
	assert.Nil(err)

	nextPage, err := collection.Next()
	assert.Nil(err)
	assert.IsType(&Collection{}, nextPage)
}

func TestSpacesServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID)

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("space-1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	space, err := cma.Spaces.Get(spaceID)
	assert.Nil(err)
	assert.Equal("id1", space.Sys.ID)
}

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

	space := &Space{
		Name:          "new space",
		DefaultLocale: "en",
	}

	err := cma.Spaces.Upsert(space)
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

	err = cma.Spaces.Upsert(space)
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

	err = cma.Spaces.Delete(space)
	assert.Nil(err)
}
