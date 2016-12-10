package contentful

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Webhook model
type Webhook struct {
	c                 *Contentful
	s                 *Space
	Sys               *Sys             `json:"sys,omitempty"`
	Name              string           `json:"name,omitempty"`
	URL               string           `json:"url,omitempty"`
	Topics            []string         `json:"topics,omitempty"`
	HTTPBasicUsername string           `json:"httpBasicUsername,omitempty"`
	HTTPBasicPassword string           `json:"httpBasicPassword,omitempty"`
	Headers           []*WebhookHeader `json:"headers,omitempty"`
}

// WebhookHeader model
type WebhookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Save the webhook
func (wh *Webhook) Save() error {
	bytesArray, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	var path string
	var method string

	if wh.Sys.CreatedAt != "" {
		path = "/spaces/" + wh.s.ID() + "/webhook_definitions/" + wh.Sys.ID
		method = "PUT"
	} else {
		path = "/spaces/" + wh.s.ID() + "/webhook_definitions"
		method = "POST"
	}

	req, err := wh.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	version := strconv.Itoa(wh.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if ok := wh.c.do(req, wh); ok != nil {
		return ok
	}

	return nil
}

// Delete the webhook
func (wh *Webhook) Delete() error {
	path := "/spaces/" + wh.s.ID() + "/webhook_definitions/" + wh.Sys.ID
	method := "DELETE"

	req, err := wh.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(wh.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = wh.c.do(req, nil); err != nil {
		return err
	}

	return nil
}
