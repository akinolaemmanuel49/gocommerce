package log

import (
	"io"
	"log"
	"os"
)

// SetupLogger initializes the logging system.
// It will log to both stdout (for handlers) and a file for (services/repositories)
func SetupLogger(fileName string) *log.Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	return log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime)
}
