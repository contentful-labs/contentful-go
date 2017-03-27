package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// ContentTypesService service
type ContentTypesService service

// ContentType model
type ContentType struct {
	Sys          *Sys     `json:"sys"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Fields       []*Field `json:"fields,omitempty"`
	DisplayField string   `json:"displayField,omitempty"`
}

const (
	FieldTypeText     = "Text"
	FieldTypeArray    = "Array"
	FieldTypeLink     = "Link"
	FieldTypeInteger  = "Integer"
	FieldTypeLocation = "Location"
	FieldTypeBoolean  = "Boolean"
	FieldTypeDate     = "Date"
	FieldTypeObject   = "Object"
)

// Field model
type Field struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	LinkType    string              `json:"linkType,omitempty"`
	Items       *FieldTypeArrayItem `json:"items,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Localized   bool                `json:"localized,omitempty"`
	Disabled    bool                `json:"disabled,omitempty"`
	Omitted     bool                `json:"omitted,omitempty"`
	Validations []FieldValidation   `json:"validations,omitempty"`
}

// FieldTypeArrayItem model
type FieldTypeArrayItem struct {
	Type        string            `json:"type"`
	Validations []FieldValidation `json:"validations"`
	LinkType    string            `json:"linkType,omitempty"`
}

// GetVersion returns entity version
func (ct *ContentType) GetVersion() int {
	version := 1
	if ct.Sys != nil {
		version = ct.Sys.Version
	}

	return version
}

// List return a content type collection
func (service *ContentTypesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/content_types", spaceID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

func (service *ContentTypesService) Get(spaceID, contentTypeID string) (*ContentType, error) {
	path := fmt.Sprintf("/spaces/%s/content_types/%s", spaceID, contentTypeID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var ct ContentType
	if err = service.c.do(req, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}

// Upsert updates or creates a new content type
func (service *ContentTypesService) Upsert(spaceID string, ct *ContentType) error {
	bytesArray, err := json.Marshal(ct)
	if err != nil {
		return err
	}

	var path string
	var method string

	if ct.Sys != nil && ct.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s/content_types/%s", spaceID, ct.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/content_types", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(ct.GetVersion()))

	if err = service.c.do(req, ct); err != nil {
		return err
	}

	return nil
}

// Delete the content_type
func (service *ContentTypesService) Delete(spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s/content_types/%s", spaceID, ct.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = service.c.do(req, nil); err != nil {
		return err
	}

	return nil
}

// Activate the contenttype, a.k.a publish
func (service *ContentTypesService) Activate(spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s/content_types/%s/published", spaceID, ct.Sys.ID)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = service.c.do(req, ct); err != nil {
		return err
	}

	return nil
}

// Deactivate the contenttype, a.k.a unpublish
func (service *ContentTypesService) Deactivate(spaceID string, ct *ContentType) error {
	path := fmt.Sprintf("/spaces/%s/content_types/%s/published", spaceID, ct.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = service.c.do(req, ct); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
