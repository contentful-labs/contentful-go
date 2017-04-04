package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalesServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/locales")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("spaces.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Locales.List(spaceID).Next()
	assert.Nil(err)
}

func TestLocalesServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := cma.Locales.Get(spaceID, "4aGeQYgByqQFJtToAOh2JJ")
	assert.Nil(err)
	assert.Equal("U.S. English", locale.Name)
	assert.Equal("en-US", locale.Code)
}

func TestLocalesServiceUpsertCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("German (Austria)", payload["name"])
		assert.Equal("de-AT", payload["code"])

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale := &Locale{
		Name: "German (Austria)",
		Code: "de-AT",
	}

	err = cma.Locales.Upsert(spaceID, locale)
	assert.Nil(err)
}

func TestLocalesServiceUpsertUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("modified-name", payload["name"])
		assert.Equal("modified-code", payload["code"])

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("locale_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := localeFromTestData("locale_1.json")
	assert.Nil(err)

	locale.Name = "modified-name"
	locale.Code = "modified-code"

	err = cma.Locales.Upsert(spaceID, locale)
	assert.Nil(err)
}

func TestLocalesServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test locale
	locale, err := localeFromTestData("locale_1.json")
	assert.Nil(err)

	// delete locale
	err = cma.Locales.Delete(spaceID, locale)
	assert.Nil(err)
}
