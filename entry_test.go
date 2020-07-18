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

func ExampleEntriesService_Upsert_create() {
	cma := NewCMA("cma-token")

	entry := &Entry{
		Sys: &Sys{
			ID: "MyEntry",
			ContentType: &ContentType{
				Sys: &Sys{
					ID: "MyContentType",
				},
			},
		},
		Fields: map[string]interface{}{
			"Description": map[string]string{
				"en-US": "Some example content...",
			},
		},
	}

	err := cma.Entries.Upsert("space-id", entry)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleEntryService_Upsert_update() {
	cma := NewCMA("cma-token")

	entry, err := cma.Entries.Get("space-id", "entry-id")
	if err != nil {
		log.Fatal(err)
	}

	entry.Fields["Description"] = map[string]interface{}{
		"en-US": "modified entry content",
	}

	err = cma.Entries.Upsert("space-id", entry)
	if err != nil {
		log.Fatal(err)
	}
}

func TestEntrySaveForCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/entries")
		assert.Equal(r.Header["X-Contentful-Content-Type"], []string{"MyContentType"})
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		assert.NotNil(payload["fields"])
		fields := payload["fields"].(map[string]interface{})

		assert.Equal(fields["Description"], map[string]interface{}{"en-US": "Some test content..."})

		w.WriteHeader(201)
		fmt.Fprintln(w, string(readTestData("entry_3.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry := &Entry{
		Sys: &Sys{
			ContentType: &ContentType{
				Sys: &Sys{
					ID: "MyContentType",
				},
			},
		},
		Fields: map[string]interface{}{
			"Description": map[string]string{
				"en-US": "Some test content...",
			},
		},
	}

	err = cma.Entries.Upsert("id1", entry)
	assert.Nil(err)
	assert.Equal("foocat", entry.Sys.ID)
}