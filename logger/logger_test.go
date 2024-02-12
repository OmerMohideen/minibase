package logger

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLogger_Info(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "info_log_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	logger := New(tmpfile, tmpfile)
	logger.Info("This is an informational message")
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read log content: %v", err)
	}
	expectedOutput := fmt.Sprintf("INFO: %s %s This is an informational message\n", time.Now().Format("2006/01/02"), time.Now().Format("15:04:05"))
	if string(content) != expectedOutput {
		t.Errorf("Info() failed: Info log output does not match expected. Got: %s, Expected: %s", content, expectedOutput)
	}
}

func TestLogger_Error(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "error_log_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	logger := New(tmpfile, tmpfile)
	logger.Error("This is an error message")
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read log content: %v", err)
	}
	expectedOutput := fmt.Sprintf("ERROR: %s %s This is an error message\n", time.Now().Format("2006/01/02"), time.Now().Format("15:04:05"))
	if string(content) != expectedOutput {
		t.Errorf("Error() failed: Error log output does not match expected. Got: %s, Expected: %s", content, expectedOutput)
	}
}
