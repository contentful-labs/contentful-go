package contentful

import "reflect"

//Entry model
type Entry struct {
	c      *Contentful
	space  *Space
	locale string
	Sys    *Sys `json:"sys"`
	Fields map[string]interface{}
}

// Get returns the entry's keys
func (e *Entry) Get(key string) *EntryField {
	col, _ := e.space.GetContentTypes().Next()
	contentTypes := col.ToContentType()

	ef := EntryField{
		c:     e.c,
		space: e.space,
		value: e.Fields[key],
	}

	for _, ct := range contentTypes {
		if ct.Sys.ID == e.Sys.ContentType.Sys.ID {
			for _, field := range ct.Fields {
				if field.ID == key {
					ef.dataType = field.Type
				}
			}
		}
	}

	return &ef
}

// EntryField model
type EntryField struct {
	c        *Contentful
	space    *Space
	value    interface{}
	dataType string
}

//String converts interface to string
func (ef *EntryField) String() string {
	return ef.value.(string)
}

//LString returns the given lovale
func (ef *EntryField) LString(locale string) string {
	m := ef.value.(map[string]interface{})

	if val, ok := m[locale]; ok {
		return val.(string)
	}

	panic("no such a locale")
}

//Integer converts interface to integer
func (ef *EntryField) Integer() int {
	return int(ef.value.(float64))
}

//LInteger converts interface to integer
func (ef *EntryField) LInteger(locale string) int {
	m := ef.value.(map[string]interface{})

	if val, ok := m[locale]; ok {
		return int(val.(float64))
	}

	panic("no such a locale")
}

//Array converts interface to slice
func (ef *EntryField) Array() []string {
	res := []string{}

	switch reflect.TypeOf(ef.value).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(ef.value)

		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface().(string))
		}
	}

	return res
}

//LArray converts interface to slice
func (ef *EntryField) LArray(locale string) []string {
	m := ef.value.(map[string]interface{})

	if val, ok := m[locale]; ok {
		res := []string{}

		switch reflect.TypeOf(val).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(val)

			for i := 0; i < s.Len(); i++ {
				res = append(res, s.Index(i).Interface().(string))
			}
		}

		return res
	}

	panic("no such a locale")
}

//LinkID returns link model
func (ef *EntryField) LinkID() string {
	m := ef.value.(map[string]interface{})
	sys := m["sys"].(map[string]interface{})
	return sys["id"].(string)
}

//LLinkID returns link model
func (ef *EntryField) LLinkID(locale string) string {
	m := ef.value.(map[string]interface{})

	if val, ok := m[locale]; ok {
		m := val.(map[string]interface{})
		sys := m["sys"].(map[string]interface{})
		return sys["id"].(string)
	}

	panic("no such a locale")
}

//LinkType returns link model
func (ef *EntryField) LinkType() string {
	m := ef.value.(map[string]interface{})
	sys := m["sys"].(map[string]interface{})
	return sys["linkType"].(string)
}

//LLinkType returns link model
func (ef *EntryField) LLinkType(locale string) string {
	m := ef.value.(map[string]interface{})

	if val, ok := m[locale]; ok {
		m := val.(map[string]interface{})
		sys := m["sys"].(map[string]interface{})
		return sys["linkType"].(string)
	}

	panic("no such a locale")
}

//Asset returns the linked asset
func (ef *EntryField) Asset() *Asset {
	if ef.LinkType() != "Asset" {
		panic("you can only convert asset types")
	}

	asset, _ := ef.space.GetAsset(ef.LinkID())
	return asset
}

//LAsset returns the linked asset
func (ef *EntryField) LAsset(locale string) *Asset {
	if ef.LLinkType(locale) != "Asset" {
		panic("you can only convert asset types")
	}

	asset, _ := ef.space.GetAsset(ef.LLinkID(locale))
	return asset
}

//Entry returns the linked entry
func (ef *EntryField) Entry() *Entry {
	if ef.LinkType() != "Entry" {
		panic("you can only convert entry types")
	}

	entry, _ := ef.space.GetEntries().Get(ef.LinkID())
	return entry
}

//LEntry returns the linked entry
func (ef *EntryField) LEntry(locale string) *Entry {
	if ef.LLinkType(locale) != "Entry" {
		panic("you can only convert entry types")
	}

	entry, _ := ef.space.GetEntries().Get(ef.LLinkID(locale))
	return entry
}
