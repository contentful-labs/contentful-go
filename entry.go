package contentful

import (
	"fmt"
	"net/url"
)

// EntriesService servıce
type EntriesService service

//Entry model
type Entry struct {
	locale string
	Sys    *Sys `json:"sys"`
	Fields map[string]interface{}
}

// GetVersion returns entity version
func (entry *Entry) GetVersion() int {
	version := 1
	if entry.Sys != nil {
		version = entry.Sys.Version
	}

	return version
}

// GetEntryKey returns the entry's keys
func (service *EntriesService) GetEntryKey(entry *Entry, key string) (*EntryField, error) {
	ef := EntryField{
		value: entry.Fields[key],
	}

	col, err := service.c.ContentTypes.List(entry.Sys.Space.Sys.ID).Next()
	if err != nil {
		return nil, err
	}

	for _, ct := range col.ToContentType() {
		if ct.Sys.ID != entry.Sys.ContentType.Sys.ID {
			continue
		}

		for _, field := range ct.Fields {
			if field.ID != key {
				continue
			}

			ef.dataType = field.Type
		}
	}

	return &ef, nil
}

// List returns entries collection
func (service *EntriesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/entries", spaceID)
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

// Get returns a single entry
func (service *EntriesService) Get(spaceID, entryID string) (*Entry, error) {
	path := fmt.Sprintf("/spaces/%s/entries/%s", spaceID, entryID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Entry{}, err
	}

	var entry Entry
	if ok := service.c.do(req, &entry); ok != nil {
		return nil, err
	}

	return &entry, err
}

// Delete the entry
func (service *EntriesService) Delete(spaceID string, entryID string) error {
	path := fmt.Sprintf("/spaces/%s/entries/%s", spaceID, entryID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Activate the entry, a.k.a publish
func (service *EntriesService) Activate(spaceID string, entryID string) error {
	path := fmt.Sprintf("/spaces/%s/entries/%s/published", spaceID, entryID)
	method := "PUT"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Deactivate the entry, a.k.a unpublish
func (service *EntriesService) Deactivate(spaceID string, entryID string) error {
	path := fmt.Sprintf("/spaces/%s/entries/%s/published", spaceID, entryID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
