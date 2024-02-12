// This package is used to handle Collection.
//
// The package includes creating a Collection and CRUD operations for records.
package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	l "github.com/OmerMohideen/minibase/logger"
	"github.com/OmerMohideen/minibase/models"
	"github.com/OmerMohideen/minibase/utils"
)

// Maximum number of records saved in a file.
// This is used for partitioning.
const MAX_CHUNK = 500

// Collection represents a collection in the database.
type Collection struct {
	mu      sync.Mutex
	name    string
	dir     string
	records map[int]*models.Record
	logger  *l.Logger
	nextID  int
}

// This function creates a new collection.
// If you want to open a collection use the name of the collection
// and use SetDir() to update its directory path.
func NewCollection(name string, logger *l.Logger) *Collection {
	collection := Collection{
		name:    name,
		records: make(map[int]*models.Record),
		logger:  logger,
		nextID:  1,
	}
	dir, _ := os.Getwd()
	collection.SetDir(dir)
	return &collection
}

// This function updates the directory of the collection.
// Use this function to open an existing collection or make one
// in a specific directory.
func (c *Collection) SetDir(dir string) {
	c.dir = dir
}

// This function gets all records from the collection
// which exists in the memory.
func (c *Collection) GetRecords() map[int]*models.Record {
	return c.records
}

// This function inserts a record into the collection.
// Note that this is saved in the memory and it is required
// to flush the records in order to save them.
func (c *Collection) InsertRecord(record *models.Record) {
	c.mu.Lock()
	defer c.mu.Unlock()

	record.ID = c.nextID
	c.records[c.nextID] = record
	c.nextID++
}

// This function gets the record by its id if available
// in the memory or pulls from the storage and caches it.
func (c *Collection) GetRecordByID(id int) (*models.Record, error) {
	c.mu.Lock()
	record, ok := c.records[id]
	c.mu.Unlock()
	if ok {
		return record, nil
	}

	if err := c.LoadRecord(id); err != nil {
		return nil, fmt.Errorf("error loading record: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	record, ok = c.records[id]
	if !ok {
		return nil, fmt.Errorf("record with ID '%d' not found even after loading", id)
	}
	return record, nil
}

// This function updates a record in the collection.
// This only updates the record in memory and it is
// required to flush the records in order to save.
func (c *Collection) UpdateRecord(id int, newRecord *models.Record) error {
	c.mu.Lock()
	_, ok := c.records[id]
	c.mu.Unlock()

	newRecord.ID = id
	if ok {
		c.records[id] = newRecord
		return nil
	}

	if err := c.LoadRecord(id); err != nil {
		return fmt.Errorf("error loading record: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok = c.records[id]
	if !ok {
		return fmt.Errorf("record with ID '%d' does not exist", id)
	}
	c.records[id] = newRecord
	return nil
}

// This function deletes a record from the collection.
// It deletes the record from the cache if exists and
// from the storage as well.
func (c *Collection) DeleteRecord(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.records[id]

	if ok {
		delete(c.records, id)
	}

	path := filepath.Join(c.dir, c.name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	min, max := utils.GetChunkRange(id, MAX_CHUNK)
	filename := fmt.Sprintf("%d-%d.json", min, max)
	file, err := os.Open(filepath.Join(c.dir, c.name, filename))
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var records []models.Record
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("error decoding data: %v", err)
	}

	var updatedRecords []models.Record
	found := false
	for _, record := range records {
		if record.ID == id {
			found = true
		} else {
			updatedRecords = append(updatedRecords, record)
		}
	}

	if !found {
		return fmt.Errorf("record with ID %d not found", id)
	}

	file, err = os.OpenFile(filepath.Join(c.dir, c.name, filename), os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(updatedRecords); err != nil {
		return fmt.Errorf("error encoding data: %v", err)
	}
	return nil
}

// This function saves the collection data to the storage.
// It partitiones the record based on its id and uses MAX_CHUNK
// as the maximum records limited to save per JSON file.
func (c *Collection) FlushRecords() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := filepath.Join(c.dir, c.name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	chunks := utils.ChunkMap(c.records, MAX_CHUNK)
	max := MAX_CHUNK
	min := 1
	for _, chunk := range chunks {
		filename := fmt.Sprintf("%d-%d.json", min, max)
		file, err := os.Create(filepath.Join(c.dir, c.name, filename))
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(chunk); err != nil {
			return fmt.Errorf("error encoding data: %v", err)
		}

		max += MAX_CHUNK
		min += MAX_CHUNK
	}

	return nil
}

// This function loads the specified record using its id
// from the storage to the memory.
func (c *Collection) LoadRecord(id int) error {
	c.mu.Lock()
	_, exists := c.records[id]
	c.mu.Unlock()
	if exists {
		return nil
	}

	path := filepath.Join(c.dir, c.name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}

	min, max := utils.GetChunkRange(id, MAX_CHUNK)
	filename := fmt.Sprintf("%d-%d.json", min, max)
	file, err := os.Open(filepath.Join(c.dir, c.name, filename))
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var records []models.Record
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("error decoding data: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert float64 fields to integers
	for _, record := range records {
		for name, value := range record.Fields {
			if floatValue, ok := value.(float64); ok {
				intValue, err := strconv.Atoi(fmt.Sprintf("%.0f", floatValue))
				if err != nil {
					return fmt.Errorf("error converting float64 to int for field %s: %v", name, err)
				}
				record.Fields[name] = intValue
			}
		}
	}

	for _, r := range records {
		if id == r.ID {
			c.records[id] = &r
			break
		}
	}
	return nil
}
