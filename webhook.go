package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// WebhooksService service
type WebhooksService service

// Webhook model
type Webhook struct {
	Sys               *Sys             `json:"sys,omitempty"`
	Name              string           `json:"name,omitempty"`
	URL               string           `json:"url,omitempty"`
	Topics            []string         `json:"topics,omitempty"`
	HTTPBasicUsername string           `json:"httpBasicUsername,omitempty"`
	HTTPBasicPassword string           `json:"httpBasicPassword,omitempty"`
	Headers           []*WebhookHeader `json:"headers,omitempty"`
	Filters           []WebhookFilter  `json:"filters,omitempty"`
}

// UnmarshalJSON for custom json unmarshaling
func (webhook *Webhook) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	if val, ok := payload["sys"]; ok {
		byteArray, err := json.Marshal(val)
		if err != nil {
			return nil
		}

		var sys Sys
		if err := json.Unmarshal(byteArray, &sys); err != nil {
			return err
		}

		webhook.Sys = &sys
	}

	if val, ok := payload["name"]; ok && val != nil {
		webhook.Name = val.(string)
	}

	if val, ok := payload["url"]; ok && val != nil {
		webhook.URL = val.(string)
	}

	if val, ok := payload["topics"]; ok && val != nil {
		byteArray, err := json.Marshal(val)
		if err != nil {
			return nil
		}

		var topics []string
		if err := json.Unmarshal(byteArray, &topics); err != nil {
			return err
		}

		webhook.Topics = topics
	}

	if val, ok := payload["httpBasicUsername"]; ok && val != nil {
		webhook.HTTPBasicUsername = val.(string)
	}

	if val, ok := payload["headers"]; ok && val != nil {
		byteArray, err := json.Marshal(val)
		if err != nil {
			return nil
		}

		var headers []*WebhookHeader
		if err := json.Unmarshal(byteArray, &headers); err != nil {
			return err
		}

		webhook.Headers = headers
	}

	if val, ok := payload["filters"]; ok && val != nil {
		filters, err := ParseFilters(val.([]interface{}))
		if err != nil {
			return err
		}

		webhook.Filters = filters
	}

	return nil
}

// ParseFilters converts json representation to go struct
func ParseFilters(data []interface{}) (filters []WebhookFilter, err error) {
	for _, value := range data {
		var byteArray []byte
		var filterOpCondition bool
		var filter map[string]interface{}

		if filterMap, ok := value.(map[string]interface{}); ok {
			byteArray, err = json.Marshal(filterMap)
			if err != nil {
				return nil, err
			}

			filterOpCondition = true
			filter = filterMap
		}

		if notFilter, ok := filter["not"].(map[string]interface{}); ok {
			byteArray, err = json.Marshal(notFilter)
			if err != nil {
				return nil, err
			}

			filterOpCondition = false
			filter = notFilter
		}

		if _, ok := filter["equals"]; ok {
			var webhookFilterEquals WebhookFilterEquals
			err = json.Unmarshal(byteArray, &webhookFilterEquals)
			if err != nil {
				return nil, err
			}

			webhookFilterEquals.Condition = filterOpCondition
			filters = append(filters, webhookFilterEquals)
		}

		if _, ok := filter["in"]; ok {
			var webhookFilterIn WebhookFilterIn
			err = json.Unmarshal(byteArray, &webhookFilterIn)
			if err != nil {
				return nil, err
			}

			webhookFilterIn.Condition = filterOpCondition
			filters = append(filters, webhookFilterIn)
		}

		if _, ok := filter["regexp"]; ok {
			var webhookFilterRegexp WebhookFilterRegexp
			err = json.Unmarshal(byteArray, &webhookFilterRegexp)
			if err != nil {
				return nil, err
			}

			webhookFilterRegexp.Condition = filterOpCondition
			filters = append(filters, webhookFilterRegexp)
		}
	}

	return filters, nil
}

// WebhookHeader model
type WebhookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetVersion returns entity version
func (webhook *Webhook) GetVersion() int {
	version := 1
	if webhook.Sys != nil {
		version = webhook.Sys.Version
	}

	return version
}

// List returns webhooks collection
func (service *WebhooksService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/webhook_definitions", spaceID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single webhook entity
func (service *WebhooksService) Get(spaceID, webhookID string) (*Webhook, error) {
	path := fmt.Sprintf("/spaces/%s/webhook_definitions/%s", spaceID, webhookID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := service.c.do(req, &webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

// Upsert updates or creates a new entity
func (service *WebhooksService) Upsert(spaceID string, webhook *Webhook) error {
	bytesArray, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	var path string
	var method string

	if webhook.Sys != nil && webhook.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s/webhook_definitions/%s", spaceID, webhook.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/webhook_definitions", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(webhook.GetVersion()))

	return service.c.do(req, webhook)
}

// Delete the webhook
func (service *WebhooksService) Delete(spaceID string, webhook *Webhook) error {
	path := fmt.Sprintf("/spaces/%s/webhook_definitions/%s", spaceID, webhook.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(webhook.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
