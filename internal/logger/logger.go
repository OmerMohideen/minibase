package logger

import (
	"log"
	"os"
)

// Logger represents a custom logger.
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// This function creates a new custom logger.
func New(infoHandle, errorHandle *os.File) *Logger {
	return &Logger{
		infoLog:  log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime),
		errorLog: log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime),
	}
}

// This function logs informational messages.
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

// This function logs error messages.
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}
