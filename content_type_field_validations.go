package contentful

import "encoding/json"

// FieldValidation interface
type FieldValidation interface{}

// FieldValidationLink model
type FieldValidationLink struct {
	LinkContentType []string `json:"linkContentType,omitempty"`
}

// FieldValidationUnique model
type FieldValidationUnique struct {
	Unique bool `json:"unique"`
}

// FieldValidationPredefinedValues model
type FieldValidationPredefinedValues struct {
	In           []interface{} `json:"in,omitempty"`
	ErrorMessage string        `json:"message"`
}

// FieldValidationRange model
type FieldValidationRange struct {
	Min          float64
	Max          float64
	ErrorMessage string
}

// MarshalJSON for custom json marshaling
func (v *FieldValidationRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Range *struct {
			Min float64 `json:"min,omitempty"`
			Max float64 `json:"max,omitempty"`
		} `json:"range,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Range: &struct {
			Min float64 `json:"min,omitempty"`
			Max float64 `json:"max,omitempty"`
		}{
			Min: v.Min,
			Max: v.Max,
		},
		Message: v.ErrorMessage,
	})
}

// UnmarshalJSON for custom json unmarshaling
func (v *FieldValidationRange) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	rangeData := payload["range"].(map[string]interface{})

	if val, ok := rangeData["min"].(float64); ok {
		v.Min = val
	}

	if val, ok := rangeData["max"].(float64); ok {
		v.Max = val
	}

	if val, ok := payload["message"].(string); ok {
		v.ErrorMessage = val
	}

	return nil
}

// FieldValidationSize model
type FieldValidationSize struct {
	Min          int
	Max          int
	ErrorMessage string
}

// MarshalJSON for custom json marshaling
func (v *FieldValidationSize) MarshalJSON() ([]byte, error) {
	type size struct {
		Min int `json:"min,omitempty"`
		Max int `json:"max,omitempty"`
	}

	return json.Marshal(&struct {
		Size    *size  `json:"size,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Size:    &size{v.Min, v.Max},
		Message: v.ErrorMessage,
	})
}

// UnmarshalJSON for custom json unmarshaling
func (v *FieldValidationSize) UnmarshalJSON(data []byte) error {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	sizeData := payload["size"].(map[string]interface{})

	if val, ok := sizeData["min"].(float64); ok {
		v.Min = int(val)
	}

	if val, ok := sizeData["max"].(float64); ok {
		v.Max = int(val)
	}

	if val, ok := payload["message"].(string); ok {
		v.ErrorMessage = val
	}

	return nil
}
