package contentful

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// EntriesService service
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
func (service *EntriesService) List(spaceID string) (*Collection, error) {
	col := NewCollection(&CollectionOptions{})
	col.c = service.c

	if service.c.Environment == "" {
		return col, errors.New("the environment must be set before calling this method")
	}

	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries"
	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return col, err
	}
	col.req = req

	if ok := service.c.do(req, &col); ok != nil {
		return nil, err
	}
	return col, nil
}

// Get returns a single entry
func (service *EntriesService) Get(spaceID, entryID string) (Entry, error) {
	entry := Entry{
		Sys:    &Sys{},
	}
	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries/" + entryID
	query := url.Values{}

	req, err := service.c.newRequest(http.MethodGet, path, query, nil)
	if err != nil {
		return entry, err
	}

	if ok := service.c.do(req, &entry); ok != nil {
		return entry, err
	}
	return entry, err
}

// Delete the entry
func (service *EntriesService) Delete(spaceID string, entryID string) error {
	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries/" + entryID

	req, err := service.c.newRequest(http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Create a single entry, this method does not publish the entry.
func (service *EntriesService) Create(spaceID string, entry *Entry) error {
	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries/"

	req, err := service.c.newRequest(http.MethodPut, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Publish the entry
func (service *EntriesService) Publish(spaceID string, entry *Entry) error {
	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries/" + entry.Sys.ID + "published"

	req, err := service.c.newRequest(http.MethodPut, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// UnPublish the entry
func (service *EntriesService) UnPublish(spaceID string, entry *Entry) error {
	path := "/spaces/" + spaceID +  "/environments/" + service.c.Environment + "/entries/" + entry.Sys.ID + "published"

	req, err := service.c.newRequest(http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(entry.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}
