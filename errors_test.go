package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundErrorResponse(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, string(readTestData("error-notfound.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test space
	_, err = cma.Spaces.Get("unknown-space-id")
	assert.NotNil(err)
	_, ok := err.(NotFoundError)
	assert.Equal(true, ok)
	notFoundError := err.(NotFoundError)
	assert.Equal(404, notFoundError.APIError.res.StatusCode)
	assert.Equal("request-id", notFoundError.APIError.err.RequestID)
	assert.Equal("The resource could not be found.", notFoundError.APIError.err.Message)
	assert.Equal("Error", notFoundError.APIError.err.Sys.Type)
	assert.Equal("NotFound", notFoundError.APIError.err.Sys.ID)
}
