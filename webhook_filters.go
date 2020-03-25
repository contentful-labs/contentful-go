package contentful

import (
	"encoding/json"
)

const (
	WebhookFilterDocContentType = "sys.contentType.sys.id"
	WebhookFilterDocEnvironment = "sys.environment.sys.id"
	WebhookFilterDocEntity      = "sys.id"
)

// WebhookFilter interface
type WebhookFilter interface{}

// WebhookFilterEquals model
type WebhookFilterEquals struct {
	Condition bool
	Doc       string
	Equals    string
}

// MarshalJSON for custom json marshaling
func (f WebhookFilterEquals) MarshalJSON() ([]byte, error) {
	if !f.Condition {
		return json.Marshal(&map[string]map[string][]interface{}{
			"not": {
				"equals": {
					f.Equals,
					struct {
						Doc string `json:"doc"`
					}{
						f.Doc,
					},
				},
			},
		})
	}
	return json.Marshal(&map[string][]interface{}{
		"equals": {
			f.Equals,
			struct {
				Doc string `json:"doc"`
			}{
				f.Doc,
			},
		},
	})
}

// UnmarshalJSON for custom json unmarshaling
func (f *WebhookFilterEquals) UnmarshalJSON(data []byte) error {
	var payload map[string][]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	equals := payload["equals"]

	for _, item := range equals {
		if value, ok := item.(string); ok {
			f.Equals = value
		}

		if valueMap, ok := item.(map[string]interface{}); ok {
			f.Doc = valueMap["doc"].(string)
		}
	}

	return nil
}

// WebhookFilterIn model
type WebhookFilterIn struct {
	Condition bool
	Doc       string
	In        []string
}

// MarshalJSON for custom json marshaling
func (f WebhookFilterIn) MarshalJSON() ([]byte, error) {
	if !f.Condition {
		return json.Marshal(&map[string]map[string][]interface{}{
			"not": {
				"in": {
					f.In,
					struct {
						Doc string `json:"doc"`
					}{
						f.Doc,
					},
				},
			},
		})
	}

	return json.Marshal(&map[string][]interface{}{
		"in": {
			f.In,
			struct {
				Doc string `json:"doc"`
			}{
				f.Doc,
			},
		},
	})
}

// UnmarshalJSON for custom json unmarshaling
func (f *WebhookFilterIn) UnmarshalJSON(data []byte) error {
	var payload map[string][]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	in := payload["in"]

	for _, item := range in {
		if values, ok := item.([]interface{}); ok {
			var in []string
			for _, value := range values {
				in = append(in, value.(string))
			}
			f.In = in
		}

		if valueMap, ok := item.(map[string]interface{}); ok {
			f.Doc = valueMap["doc"].(string)
		}
	}

	return nil
}

// WebhookFilterRegexp model
type WebhookFilterRegexp struct {
	Condition bool
	Doc       string
	Pattern   string
}

// MarshalJSON for custom json marshaling
func (f WebhookFilterRegexp) MarshalJSON() ([]byte, error) {
	if !f.Condition {
		return json.Marshal(&map[string]map[string][]interface{}{
			"not": {
				"regexp": {
					struct {
						Pattern string `json:"pattern"`
					}{
						f.Pattern,
					},
					struct {
						Doc string `json:"doc"`
					}{
						f.Doc,
					},
				},
			},
		})
	}

	return json.Marshal(&map[string][]interface{}{
		"regexp": {
			struct {
				Pattern string `json:"pattern"`
			}{
				f.Pattern,
			},
			struct {
				Doc string `json:"doc"`
			}{
				f.Doc,
			},
		},
	})
}

// UnmarshalJSON for custom json unmarshaling
func (f *WebhookFilterRegexp) UnmarshalJSON(data []byte) error {
	var payload map[string][]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	regexp := payload["regexp"]

	for _, item := range regexp {
		if valueMap, ok := item.(map[string]interface{}); ok {
			if value, ok := valueMap["doc"]; ok {
				f.Doc = value.(string)
			}
			if value, ok := valueMap["pattern"]; ok {
				f.Pattern = value.(string)
			}
		}
	}

	return nil
}
