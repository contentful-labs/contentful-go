package contentful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

// File model
type File struct {
	Name        string      `json:"fileName,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	URL         string      `json:"url,omitempty"`
	UploadURL   string      `json:"upload,omitempty"`
	Detail      *FileDetail `json:"details,omitempty"`
}

// FileDetail model
type FileDetail struct {
	Size  int        `json:"size,omitempty"`
	Image *FileImage `json:"image,omitempty"`
}

// FileImage model
type FileImage struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

// FileFields model
type FileFields struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	File        *File  `json:"file,omitempty"`
}

// Asset model
type Asset struct {
	c      *Contentful
	s      *Space
	locale string
	Sys    *Sys        `json:"sys"`
	Fields *FileFields `json:"fields"`
}

func (a *Asset) updateInternals(response map[string]interface{}) error {
	fileName := response["fields"].(map[string]interface{})["file"].(map[string]interface{})["fileName"]
	localized := true

	if fileName == nil {
		localized = false
	}

	newAsset := map[string]interface{}{}

	if localized == false {
		newAsset["sys"] = response["sys"]
		newAsset["fields"] = map[string]interface{}{}
		fields := newAsset["fields"].(map[string]interface{})

		// file
		fields["file"] = response["fields"].(map[string]interface{})["file"].(map[string]interface{})[a.locale]
		// title
		title := response["fields"].(map[string]interface{})["title"]
		if title != nil {
			fields["title"] = title.(map[string]interface{})[a.locale]
		}

		// description
		description := response["fields"].(map[string]interface{})["description"]
		if description != nil {
			fields["description"] = description.(map[string]interface{})[a.locale]
		}
	} else {
		newAsset = response
	}

	byteArray, err := json.Marshal(newAsset)
	if err != nil {
		return err
	}

	json.NewDecoder(bytes.NewReader(byteArray)).Decode(a)

	return nil
}

func newAssetFromAPIResponse(response map[string]interface{}, locale string) (*Asset, error) {
	asset := &Asset{
		locale: locale,
	}

	err := asset.updateInternals(response)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// ID return sys.id
func (a *Asset) ID() string {
	return a.Sys.ID
}

func (a *Asset) makeSaveReq() (*http.Request, error) {
	payload := map[string]interface{}{
		"sys": "",
		"fields": map[string]interface{}{
			"title":       map[string]string{},
			"description": map[string]string{},
			"file":        map[string]interface{}{},
		},
	}

	payload["sys"] = a.Sys
	fields := payload["fields"].(map[string]interface{})

	// title
	title := fields["title"].(map[string]string)
	title[a.locale] = a.Fields.Title

	// description
	description := fields["description"].(map[string]string)
	description[a.locale] = a.Fields.Description

	// file
	file := fields["file"].(map[string]interface{})
	file[a.locale] = a.Fields.File

	bytesArray, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var path string
	var method string

	if a.Sys.CreatedAt != "" {
		path = "/spaces/" + a.s.ID() + "/assets/" + a.ID()
		method = "PUT"
	} else {
		path = "/spaces/" + a.s.ID() + "/assets"
		method = "POST"
	}

	req, err := a.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, err
	}

	version := strconv.Itoa(a.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return req, nil
}

// Save saves the asset
func (a *Asset) Save() error {
	req, err := a.makeSaveReq()
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err = a.c.do(req, &response); err != nil {
		return err
	}

	err = a.updateInternals(response)
	if err != nil {
		return err
	}

	return nil
}

// Delete sends delete request
func (a *Asset) Delete() error {
	path := "/spaces/" + a.s.ID() + "/assets/" + a.ID()
	method := "DELETE"

	req, err := a.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(a.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = a.c.do(req, nil); err != nil {
		return err
	}

	return nil
}

func (a *Asset) getLatestVersion() (int, error) {
	path := "/spaces/" + a.s.ID() + "/assets/" + a.ID()
	method := "GET"

	req, err := a.c.newRequest(method, path, nil, nil)
	if err != nil {
		return 0, err
	}

	var response map[string]interface{}
	if err = a.c.do(req, &response); err != nil {
		return 0, err
	}

	sys := response["sys"].(map[string]interface{})
	version := int(sys["version"].(float64))

	return version, nil
}

func (a *Asset) updateToLatestVersoin() error {
	version, err := a.getLatestVersion()
	if err != nil {
		return err
	}

	a.Sys.Version = version

	return nil
}

// Process the asset
func (a *Asset) Process() error {
	path := "/spaces/" + a.s.ID() + "/assets/" + a.ID() + "/files/" + a.locale + "/process"
	method := "PUT"

	req, err := a.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(a.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = a.c.do(req, nil); err != nil {
		return err
	}

	return nil
}

// Publish published the asset
func (a *Asset) Publish() error {
	err := a.updateToLatestVersoin()
	if err != nil {
		return err
	}

	path := "/spaces/" + a.s.ID() + "/assets/" + a.ID() + "/published"
	method := "PUT"

	req, err := a.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(a.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	var response map[string]interface{}
	if err = a.c.do(req, &response); err != nil {
		return err
	}

	a.updateInternals(response)

	return nil
}
