package contentful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

// Locale model
type Locale struct {
	c            *Contentful
	s            *Space
	Sys          *Sys   `json:"sys"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	FallbackCode string `json:"fallbackCode"`
	Default      bool   `json:"default"`
	Optional     bool   `json:"optional,omitempty"`
	CDA          bool   `json:"contentDeliveryApi,omitempty"`
	CMA          bool   `json:"contentManagementApi,omitempty"`
}

// Save saved the locale in given space
func (l *Locale) Save() error {
	req, err := l.getSaveReq()
	if err != nil {
		return err
	}

	if ok := l.c.do(req, l); ok != nil {
		return ok
	}

	return nil
}

func (l *Locale) getSaveReq() (*http.Request, error) {
	bytesArray, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	var path string
	var method string

	if l.Sys.CreatedAt != "" {
		path = "/spaces/" + l.s.ID() + "/locales/" + l.Sys.ID
		method = "PUT"
	} else {
		path = "/spaces/" + l.s.ID() + "/locales"
		method = "POST"
	}

	req, err := l.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, err
	}

	version := strconv.Itoa(l.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return req, nil
}

// Delete the locale
func (l *Locale) Delete() error {
	path := "/spaces/" + l.s.ID() + "/locales/" + l.Sys.ID
	method := "DELETE"

	req, err := l.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(l.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = l.c.do(req, nil); err != nil {
		return err
	}

	return nil
}
