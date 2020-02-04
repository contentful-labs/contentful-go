package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// AssetsService service
type AssetsService service

// Asset represents a Contentful asset
type Asset struct {
	locale string
	Sys    *Sys         `json:"sys,omitempty"`
	Fields *AssetFields `json:"fields,omitempty"`
}

type AssetFields struct {
	Title       map[string]string `json:"title,omitempty"`
	Description map[string]string `json:"description,omitempty"`
	File        map[string]*File  `json:"file,omitempty"`
}

// File represents a Contentful File
type File struct {
	URL         string       `json:"url,omitempty"`
	UploadURL   string       `json:"upload,omitempty"`
	Details     *FileDetails `json:"details,omitempty"`
	FileName    string       `json:"fileName,omitempty"`
	ContentType string       `json:"contentType,omitempty"`
}

type FileDetails struct {
	Size  int          `json:"size,omitempty"`
	Image *ImageFields `json:"image,omitempty"`
}

type ImageFields struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// GetVersion returns entity version
func (asset *Asset) GetVersion() int {
	version := 1
	if asset.Sys != nil {
		version = asset.Sys.Version
	}

	return version
}

// List returns asset collection
func (service *AssetsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/assets", spaceID)
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

// Get returns a single asset entity
func (service *AssetsService) Get(spaceID, assetID string) (*Asset, error) {
	path := fmt.Sprintf("/spaces/%s/assets/%s", spaceID, assetID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var asset Asset
	if err := service.c.do(req, &asset); err != nil {
		return nil, err
	}

	return &asset, nil
}

func (service *AssetsService) Create(spaceID string, asset *Asset) error {
	bytesArray, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	var path, method string
	if asset.Sys != nil && asset.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s/assets/%s", spaceID, asset.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/assets", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	return service.c.do(req, asset)
}

// Upsert updates or creates a new asset entity
func (service *AssetsService) Upsert(spaceID string, asset *Asset) error {
	bytesArray, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	var path string
	var method string

	if asset.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s/assets/%s", spaceID, asset.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/assets", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(asset.GetVersion()))

	return service.c.do(req, asset)
}

// Delete sends delete request
func (service *AssetsService) Delete(spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s/assets/%s", spaceID, asset.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Process the asset
func (service *AssetsService) Process(spaceID string, asset *Asset) error {
	// TODO remove hardcoding of locale
	path := fmt.Sprintf("/spaces/%s/assets/%s/files/%s/process", spaceID, asset.Sys.ID, "en-US")
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// Publish the asset
func (service *AssetsService) Publish(spaceID string, asset *Asset) error {
	path := fmt.Sprintf("/spaces/%s/assets/%s/published", spaceID, asset.Sys.ID)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(asset.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, asset)
}
