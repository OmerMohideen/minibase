package db

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/OmerMohideen/minibase/logger"
	"github.com/OmerMohideen/minibase/models"
)

func TestCollection_InsertRecord(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)
	record := models.NewRecord()
	collection.InsertRecord(record)
	if len(collection.records) != 1 {
		t.Errorf("InsertRecord() failed: Record not found in the collection")
	}
}

func TestCollection_GetRecordByID(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)
	record := models.NewRecord()
	collection.InsertRecord(record)
	retrievedRecord, err := collection.GetRecordByID(1)
	if err != nil {
		t.Errorf("GetRecordByID() failed: %v", err)
	}
	if retrievedRecord != record {
		t.Errorf("GetRecordByID() failed: Retrieved record doesn't match inserted record")
	}
}

func TestCollection_UpdateRecord(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)
	record := models.NewRecord()
	collection.InsertRecord(record)
	updatedRecord := models.NewRecord()
	err := collection.UpdateRecord(1, updatedRecord)
	if err != nil {
		t.Errorf("UpdateRecord() failed: %v", err)
	}
	retrievedRecord, err := collection.GetRecordByID(1)
	if err != nil {
		t.Errorf("GetRecordByID() failed: %v", err)
	}
	if retrievedRecord != updatedRecord {
		t.Errorf("UpdateRecord() failed: Retrieved record doesn't match updated record")
	}
}

func TestCollection_DeleteRecord(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)
	collection.SetDir(t.TempDir())

	record := models.NewRecord()
	collection.InsertRecord(record)
	err := collection.DeleteRecord(1)
	if err != nil {
		t.Errorf("DeleteRecord() failed: %v", err)
	}
	rec, err := collection.GetRecordByID(1)
	fmt.Println(rec)
	if err == nil {
		t.Error("DeleteRecord() failed: Record still exists in the collection after deletion")
	}
}

func TestCollection_GetRecords(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)
	records := []*models.Record{
		{Fields: map[string]interface{}{"name": "Sajith", "age": 30}},
		{Fields: map[string]interface{}{"name": "Mahinda", "age": 35}},
		{Fields: map[string]interface{}{"name": "Anura", "age": 40}},
	}
	for _, record := range records {
		collection.InsertRecord(record)
	}
	allRecords := collection.GetRecords()
	for _, expectedRecord := range records {
		found := false
		for _, retrievedRecord := range allRecords {
			if reflect.DeepEqual(retrievedRecord.Fields, expectedRecord.Fields) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetRecords() failed: Record %+v not found in the map returned", expectedRecord.Fields)
		}
	}
	if len(allRecords) != len(records) {
		t.Errorf("GetRecords() failed: Length of the returned map does not match the number of inserted records. Expected: %d, Got: %d", len(records), len(allRecords))
	}
}

func TestCollection_FlushRecords(t *testing.T) {
	logger := logger.New(nil, nil)
	collection := NewCollection("test_collection", logger)

	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Sajith", "age": 30}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Mahinda", "age": 35}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Anura", "age": 40}})

	collection.SetDir(t.TempDir())
	err := collection.FlushRecords()
	if err != nil {
		t.Errorf("FlushRecords() failed: Error saving collection data to file: %v", err)
	}
}

func TestCollection_LoadRecord(t *testing.T) {
	id, logger, tempDir := 2, logger.New(nil, nil), t.TempDir()
	collection := NewCollection("test_collection", logger)
	collection.SetDir(tempDir)

	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Sajith", "age": 30}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Mahinda", "age": 35}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Anura", "age": 40}})

	err := collection.FlushRecords()
	if err != nil {
		t.Errorf("LoadRecord() failed: Error saving collection data to file: %v", err)
	}

	newcollection := NewCollection("test_collection", logger)
	newcollection.SetDir(tempDir)

	err = newcollection.LoadRecord(id)
	if err != nil {
		t.Errorf("LoadRecord() failed: Error loading collection record from file: %v", err)
	}

	oldrecord, newrecord := collection.records[id], newcollection.records[id]
	if !reflect.DeepEqual(oldrecord, newrecord) {
		t.Errorf("LoadRecords() failed: Loaded record with ID %d does not match original record", id)
	}
}

func TestCollection_SetDir(t *testing.T) {
	logger, tempDir := logger.New(nil, nil), t.TempDir()
	collection := NewCollection("test_collection", logger)
	collection.SetDir(tempDir)

	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Sajith", "age": 30}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Mahinda", "age": 35}})
	collection.InsertRecord(&models.Record{Fields: map[string]interface{}{"name": "Anura", "age": 40}})

	err := collection.FlushRecords()
	if err != nil {
		t.Errorf("SetDir() failed: Error saving collection data to file: %v", err)
	}

	newcollection := NewCollection("test_collection", logger)
	newcollection.SetDir(tempDir)

	if err != nil {
		t.Errorf("SetDir() failed: Error loading next id from file: %v", err)
	}

	if collection.nextID != newcollection.nextID {
		t.Errorf("SetDir() failed: Error next id is not the same as %d", collection.nextID)
	}
}
