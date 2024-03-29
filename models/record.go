// This package is used to handle Record.
//
// The package includes creating a Record, adding and getting
// fields from the record.
package models

import (
	"fmt"
	"reflect"
	"time"
)

// Record represents a record with customizable fields.
type Record struct {
	ID        int                    `json:"id"`
	Fields    map[string]interface{} `json:"fields"`
	ExpiresAt time.Time              `json:"-"`
	Flushed   bool                   `json:"-"`
}

// This function creates a new record.
// To store data inside using Fields.
func NewRecord() *Record {
	return &Record{
		ID:      0,
		Fields:  make(map[string]interface{}),
		Flushed: false,
	}
}

// This function adds a field to the record.
func (r *Record) AddField(name string, value interface{}) {
	r.Fields[name] = value
}

// This function retrieves a field from the record.
func (r *Record) GetField(name string) (interface{}, error) {
	value, ok := r.Fields[name]
	if !ok {
		return nil, fmt.Errorf("field '%s' does not exist", name)
	}
	return value, nil
}

// This function checks if the record has all required fields
// with specified types.
func (r *Record) Validate(schema map[string]string) error {
	for fieldName, expectedType := range schema {
		value, ok := r.Fields[fieldName]
		if !ok {
			return fmt.Errorf("required field '%s' is missing", fieldName)
		}
		if reflect.TypeOf(value).Kind().String() != expectedType {
			return fmt.Errorf("field '%s' has incorrect type, expected %s", fieldName, expectedType)
		}
	}
	return nil
}
