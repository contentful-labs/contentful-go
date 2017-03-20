package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	cma      *Contentful
	c        *Contentful
	CMAToken = "b4c0n73n7fu1"
	spaceID  = "id1"
)

func readTestData(fileName string) string {
	path := "testdata/" + fileName
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}

func checkHeaders(req *http.Request, assert *assert.Assertions) {
	assert.Equal("Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal("application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
}

func spaceFromTestData(fileName string) (*Space, error) {
	content := readTestData(fileName)
	space := Space{
		c: cma,
	}
	err := json.NewDecoder(strings.NewReader(content)).Decode(&space)
	if err != nil {
		return nil, err
	}

	return &space, nil
}

func webhookFromTestData(fileName string) (*Webhook, error) {
	// get test space
	space, err := spaceFromTestData("space-1.json")
	if err != nil {
		return nil, err
	}

	content := readTestData(fileName)
	webhook := Webhook{
		c: cma,
		s: space,
	}
	err = json.NewDecoder(strings.NewReader(content)).Decode(&webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

func setup() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture := strings.Replace(r.URL.Path, "/", "-", -1)
		fixture = strings.TrimLeft(fixture, "-")
		var path string

		if e := r.URL.Query().Get("error"); e != "" {
			path = "testdata/error-" + e + ".json"
		} else {
			if r.Method == "GET" {
				path = "testdata/" + fixture + ".json"
			}

			if r.Method == "POST" {
				path = "testdata/" + fixture + "-new.json"
			}

			if r.Method == "PUT" {
				path = "testdata/" + fixture + "-updated.json"
			}
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintln(w, string(file))
		return
	})

	server = httptest.NewServer(handler)

	c = NewCMA(CMAToken)
	c.BaseURL = server.URL
}

func teardown() {
	server.Close()
	c = nil
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	method := "GET"
	path := "/some/path"
	query := url.Values{}
	query.Add("foo", "bar")
	query.Add("faz", "zoo")

	expectedURL, _ := url.Parse(c.BaseURL)
	expectedURL.Path = path
	expectedURL.RawQuery = query.Encode()

	req, err := c.newRequest(method, path, query, nil)
	assert.Nil(err)
	assert.Equal(req.Header.Get("Authorization"), "Bearer "+CMAToken)
	assert.Equal(req.Header.Get("Content-Type"), "application/vnd.contentful.management.v1+json")
	assert.Equal(req.Method, method)
	assert.Equal(req.URL.String(), expectedURL.String())

	method = "POST"
	type RequestBody struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	bodyData := RequestBody{
		Name: "test",
		Age:  10,
	}
	body, _ := json.Marshal(bodyData)
	req, err = c.newRequest(method, path, query, bytes.NewReader(body))
	assert.Nil(err)
	assert.Equal(req.Header.Get("Authorization"), "Bearer "+CMAToken)
	assert.Equal(req.Header.Get("Content-Type"), "application/vnd.contentful.management.v1+json")
	assert.Equal(req.Method, method)
	assert.Equal(req.URL.String(), expectedURL.String())
	defer req.Body.Close()
	var requestBody RequestBody
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	assert.Nil(err)
	assert.Equal(requestBody, bodyData)
}

func TestHandleError(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	method := "GET"
	path := "/some/path"
	requestID := "request-id"
	query := url.Values{}
	errResponse := ErrorResponse{
		Sys: &Sys{
			ID:   "AccessTokenInvalid",
			Type: "Error",
		},
		Message:   "Access token is invalid",
		RequestID: requestID,
	}

	marshaled, _ := json.Marshal(errResponse)
	errResponseReader := bytes.NewReader(marshaled)
	errResponseReadCloser := ioutil.NopCloser(errResponseReader)

	req, _ := c.newRequest(method, path, query, nil)
	responseHeaders := http.Header{}
	responseHeaders.Add("X-Contentful-Request-Id", requestID)
	res := &http.Response{
		Header:     responseHeaders,
		StatusCode: http.StatusUnauthorized,
		Body:       errResponseReadCloser,
		Request:    req,
	}

	err := c.handleError(req, res)
	assert.IsType(AccessTokenInvalidError{}, err)
	assert.Equal(req, err.(AccessTokenInvalidError).APIError.req)
	assert.Equal(res, err.(AccessTokenInvalidError).APIError.res)
	assert.Equal(&errResponse, err.(AccessTokenInvalidError).APIError.err)
}

func TestGetSpace(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, err := c.GetSpace(spaceID)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("Space", space.Sys.Type)
	assert.Equal("id1", space.Sys.ID)
	assert.Equal("id1", space.ID())
	assert.Equal("Contentful Example API", space.Name)
}

func TestGetSpaces(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	col, err := c.GetSpaces().Next()
	if err != nil {
		assert.Fail(err.Error())
	}
	spaces := col.ToSpace()

	assert.Equal("Space", spaces[0].Sys.Type)
	assert.Equal("id1", spaces[0].Sys.ID)
	assert.Equal("id1", spaces[0].ID())
	assert.Equal("Contentful Example API", spaces[0].Name)
}
