package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeSaveForCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/content_types")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("ct-name", payload["name"])
		assert.Equal("ct-description", payload["description"])

		fields := payload["fields"].([]interface{})
		assert.Equal(2, len(fields))

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})

		assert.Equal("field1", field1["id"].(string))
		assert.Equal("field1-name", field1["name"].(string))
		assert.Equal("Symbol", field1["type"].(string))

		assert.Equal("field2", field2["id"].(string))
		assert.Equal("field2-name", field2["name"].(string))
		assert.Equal("Symbol", field2["type"].(string))
		assert.Equal(true, field2["disabled"].(bool))

		assert.Equal(field1["id"].(string), payload["displayField"])

		w.WriteHeader(201)
		fmt.Fprintln(w, string(readTestData("content_type.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test space
	space, err := spaceFromTestData("space-1.json")
	assert.Nil(err)

	// webhook
	ct := space.NewContentType()
	ct.Name = "ct-name"
	ct.Type = "ct-type"
	ct.Description = "ct-description"

	fields := []*Field{
		&Field{
			ID:       "field1",
			Name:     "field1-name",
			Type:     "Symbol",
			Required: true,
		},
		&Field{
			ID:       "field2",
			Name:     "field2-name",
			Type:     "Symbol",
			Disabled: true,
		},
	}

	ct.Fields = fields
	ct.DisplayField = fields[0].ID

	err = ct.Save()
	assert.Nil(err)
	assert.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assert.Equal("ct-name", ct.Name)
	assert.Equal("ct-description", ct.Description)
}

func TestContentTypeSaveForUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("ct-name-updated", payload["name"])
		assert.Equal("ct-description-updated", payload["description"])

		fields := payload["fields"].([]interface{})
		assert.Equal(3, len(fields))

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})
		field3 := fields[2].(map[string]interface{})

		assert.Equal("field1", field1["id"].(string))
		assert.Equal("field1-name-updated", field1["name"].(string))
		assert.Equal("String", field1["type"].(string))

		assert.Equal("field2", field2["id"].(string))
		assert.Equal("field2-name-updated", field2["name"].(string))
		assert.Equal("Integer", field2["type"].(string))
		assert.Nil(field2["disabled"])

		assert.Equal("field3", field3["id"].(string))
		assert.Equal("field3-name", field3["name"].(string))
		assert.Equal("Date", field3["type"].(string))

		assert.Equal(field3["id"].(string), payload["displayField"])

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("content_type-updated.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assert.Nil(err)

	ct.Name = "ct-name-updated"
	ct.Type = "ct-type-updated"
	ct.Description = "ct-description-updated"

	field1 := ct.Fields[0]
	field1.Name = "field1-name-updated"
	field1.Type = "String"
	field1.Required = false

	field2 := ct.Fields[1]
	field2.Name = "field2-name-updated"
	field2.Type = "Integer"
	field2.Disabled = false

	field3 := &Field{
		ID:   "field3",
		Name: "field3-name",
		Type: "Date",
	}

	ct.Fields = append(ct.Fields, field3)
	ct.DisplayField = ct.Fields[2].ID

	err = ct.Save()
	assert.Nil(err)
	assert.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assert.Equal("ct-name-updated", ct.Name)
	assert.Equal("ct-description-updated", ct.Description)
	assert.Equal(2, ct.Sys.Version)
}

func TestContentTypeDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assert.Nil(err)

	err = ct.Delete()
	assert.Nil(err)
}
