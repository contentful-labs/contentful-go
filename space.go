package contentful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

//Space model
type Space struct {
	c             *Contentful
	Sys           *Sys   `json:"sys,omitempty"`
	Name          string `json:"name,omitempty"`
	DefaultLocale string `json:"defaultLocale,omitempty"`
}

// ID returns the sys.id
func (s *Space) ID() string {
	return s.Sys.ID
}

// GetLocales returns a locales collection
func (s *Space) GetLocales() *Collection {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/locales"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = s.c
	col.s = s
	col.req = req

	return col
}

// GetLocale returns a single locale entity
func (s *Space) GetLocale(localeID string) (*Locale, error) {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/locales/" + localeID

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	locale := &Locale{
		c: s.c,
		s: s,
	}

	if err := s.c.do(req, locale); err != nil {
		return nil, err
	}

	return locale, nil
}

// NewContentType creates a new content type object
func (s *Space) NewContentType() *ContentType {
	return &ContentType{
		c: s.c,
		s: s,
		Sys: &Sys{
			Type:    "ContentType",
			Version: 1,
		},
	}
}

// NewAPIKey creates an apikey instance
func (s *Space) NewAPIKey() *APIKey {
	return &APIKey{
		c: s.c,
		s: s,
	}
}

// GetAPIKeys returns a apikey collection
func (s *Space) GetAPIKeys() *Collection {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/api_keys"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = s.c
	col.s = s
	col.req = req

	return col
}

// GetAPIKey returns a single apikey entity
func (s *Space) GetAPIKey(apiKeyID string) (*APIKey, error) {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/api_keys/" + apiKeyID

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	apiKey := &APIKey{
		c: s.c,
		s: s,
	}

	if err := s.c.do(req, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// NewWebhook creates an apikey instance
func (s *Space) NewWebhook() *Webhook {
	return &Webhook{
		c:   s.c,
		s:   s,
		Sys: &Sys{},
	}
}

// GetWebhooks returns a webhook collection
func (s *Space) GetWebhooks() *Collection {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/webhook_definitions"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = s.c
	col.s = s
	col.req = req

	return col
}

// GetWebhook returns a single webhook entity
func (s *Space) GetWebhook(webhookID string) (*Webhook, error) {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/webhook_definitions/" + webhookID

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	webhook := s.NewWebhook()

	if err := s.c.do(req, webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

// NewLocale creates a new locale struct
func (s *Space) NewLocale() *Locale {
	return &Locale{
		c: s.c,
		s: s,
		Sys: &Sys{
			Type:    "Locale",
			Version: 0,
		},
	}
}

// GetEntries return entries for the space
func (s *Space) GetEntries() *Entries {
	return &Entries{
		Query: *NewQuery(),
		c:     s.c,
		space: s,
	}
}

// GetAssets returns a assets collection
func (s *Space) GetAssets() *Collection {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/assets"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = s.c
	col.s = s
	col.req = req

	return col
}

// GetAsset returns a asset
func (s *Space) GetAsset(assetID string) (*Asset, error) {
	return s.GetAssetWithLocale(assetID, "en-US")
}

// GetAssetWithLocale returns an asset with given locele
func (s *Space) GetAssetWithLocale(assetID string, locale string) (*Asset, error) {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/assets/" + assetID
	query := url.Values{}
	query.Add("locale", locale)

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var rawAsset map[string]interface{}
	if err := s.c.do(req, &rawAsset); err != nil {
		return nil, err
	}

	asset, err := newAssetFromAPIResponse(rawAsset, locale)
	if err != nil {
		return nil, err
	}

	asset.c = s.c
	asset.s = s

	return asset, nil
}

// NewAsset creates a new asset struct
func (s *Space) NewAsset() *Asset {
	return &Asset{
		c:      s.c,
		s:      s,
		locale: "en-US",
		Sys: &Sys{
			Version: 1,
		},
		Fields: &FileFields{},
	}
}

// GetContentType return a content type
func (s *Space) GetContentType(contentTypeID string) (*ContentType, error) {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/content_types/" + contentTypeID

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	contentType := &ContentType{
		c: s.c,
		s: s,
	}

	if err := s.c.do(req, contentType); err != nil {
		return nil, err
	}

	return contentType, nil
}

// GetContentTypes return a collection
func (s *Space) GetContentTypes() *Collection {
	method := "GET"
	path := "/spaces/" + s.Sys.ID + "/content_types"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection(&CollectionOptions{})
	col.c = s.c
	col.s = s
	col.req = req

	return col
}

func (s *Space) getSaveReq() (*http.Request, error) {
	bytesArray, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var path string
	var method string

	if s.Sys.CreatedAt != "" {
		path = "/spaces/" + s.ID()
		method = "PUT"
	} else {
		path = "/spaces"
		method = "POST"
	}

	req, err := s.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, err
	}

	version := strconv.Itoa(s.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return req, nil
}

// Save saves the space
func (s *Space) Save() error {
	req, err := s.getSaveReq()
	if err != nil {
		return err
	}

	if ok := s.c.do(req, s); ok != nil {
		return ok
	}

	return nil
}

// Delete the space
func (s *Space) Delete() error {
	path := "/spaces/" + s.ID()
	method := "DELETE"

	req, err := s.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(s.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = s.c.do(req, nil); err != nil {
		return err
	}

	return nil
}
