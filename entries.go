package contentful

//Entries model
type Entries struct {
	Query
	c     *Contentful
	space *Space
}

//Locale localizes the entries
func (es *Entries) Locale(locale string) *Entries {
	es.Query.Locale(locale)
	return es
}

func (es *Entries) initEntry(entry *Entry) *Entry {
	entry.c = es.c
	entry.space = es.space
	entry.locale = es.Query.locale
	return entry
}

//Get returns a single entryu
func (es *Entries) Get(id string) (entry *Entry, err error) {
	path := "/spaces/" + es.space.Sys.ID + "/entries/" + id
	query := es.Query.Values()

	req, err := es.c.newRequest("GET", path, query, nil)
	if err != nil {
		return &Entry{}, err
	}

	if ok := es.c.do(req, &entry); ok != nil {
		return &Entry{}, err
	}

	return es.initEntry(entry), err
}

// All returns all entries
func (es *Entries) All() ([]Entry, error) {
	// u, _ := url.Parse("/spaces/" + es.space.Sys.ID + "/entries")
	// u.RawQuery = es.Query.String()

	path := "/spaces/" + es.space.Sys.ID + "/entries"
	query := es.Query.Values()

	req, err := es.c.newRequest("GET", path, query, nil)
	if err != nil {
		return []Entry{}, err
	}

	response := struct {
		Sys   Sys
		Total int
		Skip  int
		Limit int
		Items []Entry
	}{}

	if ok := es.c.do(req, &response); ok != nil {
		return []Entry{}, err
	}

	entries := []Entry{}
	for _, entry := range response.Items {
		entries = append(entries, *es.initEntry(&entry))
	}

	return entries, nil
}
