package models

import (
	"testing"
)

func TestRecord_GetField(t *testing.T) {
	record := NewRecord()
	record.AddField("name", "Sajith")
	record.AddField("age", 30)
	name, err := record.GetField("name")
	if err != nil {
		t.Errorf("GetField() failed: Error getting field 'name': %v", err)
	}
	if name != "Sajith" {
		t.Errorf("GetField() failed: Expected name to be 'Sajith', got '%s'", name)
	}
}

func TestRecord_Validate(t *testing.T) {
	record := NewRecord()
	record.AddField("name", "Sajith")
	record.AddField("age", 30)
	schema := map[string]string{
		"name":  "string",
		"age":   "int",
		"email": "string",
	}
	err := record.Validate(schema)
	if err == nil {
		t.Errorf("Validate() failed: Expected validation error for missing 'email' field")
	}
	record.AddField("email", "Sajith@example.com")
	err = record.Validate(schema)
	if err != nil {
		t.Errorf("Validate() failed: Validation error: %v", err)
	}
}
