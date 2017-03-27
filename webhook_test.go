package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookSaveForCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("webhook-name", payload["name"])
		assert.Equal("https://www.example.com/test", payload["url"])
		assert.Equal("username", payload["httpBasicUsername"])
		assert.Equal("password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assert.Equal(2, len(topics))
		assert.Equal("Entry.create", topics[0].(string))
		assert.Equal("ContentType.create", topics[1].(string))

		headers := payload["headers"].([]interface{})
		assert.Equal(2, len(headers))
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assert.Equal("header1", header1["key"].(string))
		assert.Equal("header1-value", header1["value"].(string))

		assert.Equal("header2", header2["key"].(string))
		assert.Equal("header2-value", header2["value"].(string))

		w.WriteHeader(201)
		fmt.Fprintln(w, string(readTestData("webhook.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	webhook := &Webhook{
		Name: "webhook-name",
		URL:  "https://www.example.com/test",
		Topics: []string{
			"Entry.create",
			"ContentType.create",
		},
		HTTPBasicUsername: "username",
		HTTPBasicPassword: "password",
		Headers: []*WebhookHeader{
			&WebhookHeader{
				Key:   "header1",
				Value: "header1-value",
			},
			&WebhookHeader{
				Key:   "header2",
				Value: "header2-value",
			},
		},
	}

	err = cma.Webhooks.Upsert(spaceID, webhook)
	assert.Nil(err)
	assert.Equal("7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assert.Equal("webhook-name", webhook.Name)
	assert.Equal("username", webhook.HTTPBasicUsername)
}

func TestWebhookSaveForUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI")
		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("updated-webhook-name", payload["name"])
		assert.Equal("https://www.example.com/test-updated", payload["url"])
		assert.Equal("updated-username", payload["httpBasicUsername"])
		assert.Equal("updated-password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assert.Equal(3, len(topics))
		assert.Equal("Entry.create", topics[0].(string))
		assert.Equal("ContentType.create", topics[1].(string))
		assert.Equal("Asset.create", topics[2].(string))

		headers := payload["headers"].([]interface{})
		assert.Equal(2, len(headers))
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assert.Equal("header1", header1["key"].(string))
		assert.Equal("updated-header1-value", header1["value"].(string))

		assert.Equal("header2", header2["key"].(string))
		assert.Equal("updated-header2-value", header2["value"].(string))

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("webhook-updated.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData("webhook.json")
	assert.Nil(err)

	webhook.Name = "updated-webhook-name"
	webhook.URL = "https://www.example.com/test-updated"
	webhook.Topics = []string{
		"Entry.create",
		"ContentType.create",
		"Asset.create",
	}
	webhook.HTTPBasicUsername = "updated-username"
	webhook.HTTPBasicPassword = "updated-password"
	webhook.Headers = []*WebhookHeader{
		&WebhookHeader{
			Key:   "header1",
			Value: "updated-header1-value",
		},
		&WebhookHeader{
			Key:   "header2",
			Value: "updated-header2-value",
		},
	}

	err = cma.Webhooks.Upsert(spaceID, webhook)
	assert.Nil(err)
	assert.Equal("7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assert.Equal(1, webhook.Sys.Version)
	assert.Equal("updated-webhook-name", webhook.Name)
	assert.Equal("updated-username", webhook.HTTPBasicUsername)
}

func TestWebhookDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI")
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData("webhook.json")
	assert.Nil(err)

	err = cma.Webhooks.Delete(spaceID, webhook)
	assert.Nil(err)
}
