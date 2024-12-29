package errors

import (
	"encoding/json"
	"log"
	"net/http"

	l "github.com/akinolaemmanuel49/gocommerce/log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleError maps error types to HTTP responses
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	errorLogger := l.SetupLogger("service.log", "ERROR")
	if err == nil {
		writeErrorResponse(w, r, http.StatusInternalServerError, "an unexpected error occurred", errorLogger)
		return
	}

	switch err.(type) {
	case *ValidationError:
		writeErrorResponse(w, r, http.StatusBadRequest, err.Error(), errorLogger)
	case *ConflictError:
		writeErrorResponse(w, r, http.StatusConflict, err.Error(), errorLogger)
	case *NotFoundError:
		writeErrorResponse(w, r, http.StatusNotFound, err.Error(), errorLogger)
	case *InternalServerError:
		writeErrorResponse(w, r, http.StatusInternalServerError, "internal server error", errorLogger) // Mask internal details
	default:
		log.Printf("Unhandled error: %v", err)
		writeErrorResponse(w, r, http.StatusInternalServerError, "an unexpected error occurred", errorLogger)
	}
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, errorLogger *log.Logger) {
	w.WriteHeader(statusCode)
	errorLogger.Printf("%s %s [User-Agent: %s]: %v", r.Method, r.URL.Path, r.UserAgent(), message)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		http.Error(w, `{"error":"failed to encode error response"}`, http.StatusInternalServerError)
	}
}
