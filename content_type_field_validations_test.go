package contentful

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldValidationLink(t *testing.T) {
	var err error
	assert := assert.New(t)

	validation := &FieldValidationLink{
		LinkContentType: []string{"test", "test2"},
	}

	data, err := json.Marshal(validation)
	assert.Nil(err)
	assert.Equal("{\"linkContentType\":[\"test\",\"test2\"]}", string(data))
}

func TestFieldValidationUnique(t *testing.T) {
	var err error
	assert := assert.New(t)

	validation := &FieldValidationUnique{
		Unique: false,
	}

	data, err := json.Marshal(validation)
	assert.Nil(err)
	assert.Equal("{\"unique\":false}", string(data))
}

func TestFieldValidationPredefinedValues(t *testing.T) {
	var err error
	assert := assert.New(t)

	validation := &FieldValidationPredefinedValues{
		In:           []interface{}{5, 10, "string", 6.4},
		ErrorMessage: "error message",
	}

	data, err := json.Marshal(validation)
	assert.Nil(err)
	assert.Equal("{\"in\":[5,10,\"string\",6.4],\"message\":\"error message\"}", string(data))
}

func TestFieldValidationRange(t *testing.T) {
	var err error
	assert := assert.New(t)

	validation := &FieldValidationRange{
		Min:          60,
		Max:          100,
		ErrorMessage: "error message",
	}

	data, err := json.Marshal(validation)
	assert.Nil(err)
	assert.Equal("{\"range\":{\"min\":60,\"max\":100},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationRange
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assert.Nil(err)
	assert.Equal(float64(60), validationCheck.Min)
	assert.Equal(float64(100), validationCheck.Max)
	assert.Equal("error message", validationCheck.ErrorMessage)
}
