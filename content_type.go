package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

//ContentType model
type ContentType struct {
	c            *Contentful
	s            *Space
	Sys          *Sys     `json:"sys"`
	Type         string   `json:"type,omitempty"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Fields       []*Field `json:"fields,omitempty"`
	DisplayField string   `json:"displayField,omitempty"`
}

// Field model
type Field struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Items     *Schema `json:"items,omitempty"`
	Required  bool    `json:"required,omitempty"`
	Localized bool    `json:"localized,omitempty"`
	Disabled  bool    `json:"disabled,omitempty"`
	Omitted   bool    `json:"omitted,omitempty"`
}

// Schema model
type Schema struct {
	Type        string   `json:"type"`
	Validations []string `json:"validations"`
	LinkType    string   `json:"linkType,omitempty"`
}

// AddField to the content type
func (ct *ContentType) AddField(field *Field) {
	ct.Fields = append(ct.Fields, field)
}

// Save the content type to contentful
func (ct *ContentType) Save() error {
	bytesArray, err := json.Marshal(ct)
	if err != nil {
		return err
	}

	var path string
	var method string

	if ct.Sys.CreatedAt != "" {
		path = "/spaces/" + ct.s.ID() + "/content_types/" + ct.Sys.ID
		method = "PUT"
	} else {
		path = "/spaces/" + ct.s.ID() + "/content_types"
		method = "POST"
	}

	req, err := ct.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = ct.c.do(req, ct); err != nil {
		return err
	}

	return nil
}

// Activate the contenttype, a.k.a publish
func (ct *ContentType) Activate() error {
	path := "/spaces/" + ct.s.ID() + "/content_types/" + ct.Sys.ID + "/published"
	method := "PUT"

	req, err := ct.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = ct.c.do(req, ct); err != nil {
		return err
	}

	return nil
}

// Deactivate the contenttype, a.k.a unpublish
func (ct *ContentType) Deactivate() error {
	path := "/spaces/" + ct.s.ID() + "/content_types/" + ct.Sys.ID + "/published"
	method := "DELETE"

	req, err := ct.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = ct.c.do(req, ct); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Delete the content_type
func (ct *ContentType) Delete() (err error) {
	if ct.Sys.PublishedAt != "" {
		err = ct.Deactivate()
		if err != nil {
			return err
		}
	}

	path := "/spaces/" + ct.s.ID() + "/content_types/" + ct.Sys.ID
	method := "DELETE"

	req, err := ct.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(ct.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	if err = ct.c.do(req, nil); err != nil {
		return err
	}

	return nil
}
