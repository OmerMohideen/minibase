package main

import (
	"math/rand"
	"os"
	"testing"

	"github.com/OmerMohideen/minibase/internal/db"
	"github.com/OmerMohideen/minibase/internal/logger"
	"github.com/OmerMohideen/minibase/internal/models"
)

const LIMIT = 10000

// Benchmark insert function and flush function
func BenchmarkInsertRecords(b *testing.B) {
	collection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	collection.SetDir(b.TempDir())
	b.ResetTimer()

	for i := 0; i < LIMIT; i++ {
		record := models.NewRecord()
		record.AddField("age", rand.Intn(100))
		record.AddField("name", "Mahinda")
		collection.InsertRecord(record)
	}
	collection.FlushRecords()
}

// Benchmark getting non cached records
func BenchmarkGetRecordByID(b *testing.B) {
	tempDir := b.TempDir()

	collection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	collection.SetDir(tempDir)

	for i := 0; i < LIMIT; i++ {
		record := models.NewRecord()
		record.AddField("age", rand.Intn(100))
		record.AddField("name", "Mahinda")
		collection.InsertRecord(record)
	}
	collection.FlushRecords()

	newcollection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	newcollection.SetDir(tempDir)

	b.ResetTimer()

	for _, record := range collection.GetRecords() {
		newcollection.GetRecordByID(record.ID)
	}
}

// Benchmark updating non cached records
func BenchmarkUpdateRecord(b *testing.B) {
	tempDir := b.TempDir()

	collection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	collection.SetDir(tempDir)

	for i := 0; i < LIMIT; i++ {
		record := models.NewRecord()
		record.AddField("age", rand.Intn(100))
		record.AddField("name", "Mahinda")
		collection.InsertRecord(record)
	}
	collection.FlushRecords()

	newcollection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	newcollection.SetDir(tempDir)

	b.ResetTimer()

	for _, record := range collection.GetRecords() {
		newcollection.UpdateRecord(record.ID, record)
	}
}

// Benchmark updating non cached records
func BenchmarkDeleteRecord(b *testing.B) {
	tempDir := b.TempDir()

	collection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	collection.SetDir(tempDir)

	for i := 0; i < LIMIT; i++ {
		record := models.NewRecord()
		record.AddField("age", rand.Intn(100))
		record.AddField("name", "Mahinda")
		collection.InsertRecord(record)
	}
	collection.FlushRecords()

	newcollection := db.NewCollection("minibase", logger.New(os.Stdout, os.Stderr))
	newcollection.SetDir(tempDir)

	b.ResetTimer()

	for _, record := range collection.GetRecords() {
		newcollection.DeleteRecord(record.ID)
	}
}
