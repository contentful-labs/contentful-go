package contentful

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// APIKey model
type APIKey struct {
	c             *Contentful
	s             *Space
	Sys           *Sys            `json:"sys,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description,omitempty"`
	AccessToken   string          `json:"accessToken,omitempty"`
	Policies      []*APIKeyPolicy `json:"policies,omitempty"`
	PreviewAPIKey *PreviewAPIKey  `json:"preview_api_key,omitempty"`
}

// APIKeyPolicy model
type APIKeyPolicy struct {
	Effect  string `json:"effect,omitempty"`
	Actions string `json:"actions,omitempty"`
}

// PreviewAPIKey model
type PreviewAPIKey struct {
	Sys *Sys
}

// Save the apikey
func (ak *APIKey) Save() error {
	bytesArray, err := json.Marshal(ak)
	if err != nil {
		return err
	}

	var path string
	var method string

	if ak.Sys.CreatedAt != "" {
		path = "/spaces/" + ak.s.ID() + "/api_keys/" + ak.Sys.ID
		method = "PUT"
	} else {
		path = "/spaces/" + ak.s.ID() + "/api_keys"
		method = "POST"
	}

	req, err := ak.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	version := strconv.Itoa(ak.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if ok := ak.c.do(req, ak); ok != nil {
		return ok
	}

	return nil
}

// Delete the locale
func (ak *APIKey) Delete() error {
	path := "/spaces/" + ak.s.ID() + "/api_keys/" + ak.Sys.ID
	method := "DELETE"

	req, err := ak.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ak.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = ak.c.do(req, nil); err != nil {
		return err
	}

	return nil
}
