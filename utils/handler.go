package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
)

// ValidateID checks to see if the path provides an :id value
func ValidateID(ID, entity string) error {
	if ID == "" {
		return errors.NewValidationError("ID", fmt.Sprintf("%s ID is required", entity))
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	if rw, ok := w.(*errors.ErrorResponseWriter); ok && !rw.Written {
		rw.WriteHeader(statusCode)
	} else {
		w.WriteHeader(statusCode)
	}

	logger.Printf("%s %d %s [User-Agent: %s]", r.Method, statusCode, r.URL.Path, r.UserAgent())
	json.NewEncoder(w).Encode(data)
}
