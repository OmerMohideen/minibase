package models

import "fmt"

// Record represents a record with customizable fields.
type Record struct {
	ID     int
	Fields map[string]interface{}
}

// This function creates a new record.
// To store data inside using Fields.
func NewRecord() *Record {
	return &Record{
		ID:     0,
		Fields: make(map[string]interface{}),
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

// This function checks if the record has all required fields.
func (r *Record) Validate(requiredFields []string) error {
	for _, field := range requiredFields {
		if _, ok := r.Fields[field]; !ok {
			return fmt.Errorf("required field '%s' is missing", field)
		}
	}
	return nil
}
