package utils

import (
	"sort"

	"github.com/OmerMohideen/minibase/models"
)

// This function gets the range which the id belongs.
// Returns the starting range and ending range.
// example: ID: 3, MAX_CHUNK: 500 -> RANGE: 1-500
func GetChunkRange(recordID, chunkSize int) (start, end int) {
	fileIndex := (recordID - 1) / chunkSize
	start = fileIndex*chunkSize + 1
	end = start + chunkSize - 1
	return start, end
}

// This function splits the map into chunks using
// the specified chunkSize. Returns an array of chunks.
func ChunkMap(records map[int]*models.Record, chunkSize int) [][]*models.Record {
	var chunks [][]*models.Record

	keys := make([]int, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}
		chunkKeys := keys[i:end]

		chunk := make([]*models.Record, 0, len(chunkKeys))
		for _, key := range chunkKeys {
			chunk = append(chunk, records[key])
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}
