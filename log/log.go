package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Mutex to protect file access
var fileMutex sync.Mutex

// FileWriter is a custom writer that ensures thread-safe writes to a file
type FileWriter struct {
	file *os.File
}

// NewLogFileWriter creates a new instance of LogFileWriter
func NewLogFileWriter(fileName string) (*FileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	return &FileWriter{file: file}, nil
}

// Write implements the io.Writer interface and ensures that log file writes are thread-safe
func (w *FileWriter) Write(p []byte) (n int, err error) {
	// Lock the mutex before writing to the file
	fileMutex.Lock()
	defer fileMutex.Unlock()

	return w.file.Write(p)
}

// SetupLogger initializes the logging system.
// It will log to both stdout (for handlers) and a file for (services/repositories)
func SetupLogger(fileName, logLevel string) (*log.Logger, error) {
	// Create the log file writer
	fileWriter, err := NewLogFileWriter(fileName)
	if err != nil {
		return nil, err
	}

	// Create a multi-writer that writes to both stdout and the file
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// Create and return the logger
	return log.New(multiWriter, fmt.Sprintf("%-7s: ", logLevel), log.Ldate|log.Ltime), nil
}
