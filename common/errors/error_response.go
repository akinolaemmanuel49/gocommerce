package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleError maps error types to HTTP responses
func HandleError(w http.ResponseWriter, r *http.Request, err error, errorLogger *log.Logger) {
	if err == nil {
		WriteErrorResponse(w, r, http.StatusInternalServerError, "an unexpected error occurred", errorLogger)
		return
	}

	switch err.(type) {
	case *ValidationError:
		WriteErrorResponse(w, r, http.StatusBadRequest, err.Error(), errorLogger)
	case *BadRequestError:
		WriteErrorResponse(w, r, http.StatusBadRequest, err.Error(), errorLogger)
	case *ConflictError:
		WriteErrorResponse(w, r, http.StatusConflict, err.Error(), errorLogger)
	case *NotFoundError:
		WriteErrorResponse(w, r, http.StatusNotFound, err.Error(), errorLogger)
	case *AuthorizationError:
		WriteErrorResponse(w, r, http.StatusUnauthorized, err.Error(), errorLogger)
	case *ForbiddenError:
		WriteErrorResponse(w, r, http.StatusForbidden, err.Error(), errorLogger)
	case *InternalServerError:
		WriteErrorResponse(w, r, http.StatusInternalServerError, "Internal server error", errorLogger) // Mask internal details
	default:
		errorLogger.Printf("Unhandled error: %v", err)
		WriteErrorResponse(w, r, http.StatusInternalServerError, "An unexpected error occurred", errorLogger)
	}
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, errorLogger *log.Logger) {
	w.WriteHeader(statusCode)
	errorLogger.Printf("%s %d %s [User-Agent: %s]: %v", r.Method, statusCode, r.URL.Path, r.UserAgent(), message)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		http.Error(w, `{"error":"failed to encode error response"}`, http.StatusInternalServerError)
	}
}
