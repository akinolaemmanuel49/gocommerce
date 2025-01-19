package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	auth_middlewares "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
)

// ValidateID checks to see if the path provides an :id value :: Possibly deprecated
func ValidateID(ID, entity string) error {
	if ID == "" {
		return errors.NewValidationError("ID", fmt.Sprintf("%s ID is required", entity))
	}
	return nil
}

// WriteJSON writes a JSON response to the client
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

// IsAuthorized checks if the user is authorized to access a protected route
func IsAuthorized(ctx context.Context) (*models.JWTClaims, error) {
	claims := auth_middlewares.GetClaims(ctx)
	if claims == nil {
		return nil, errors.NewAuthorizationError("")
	}
	return claims, nil
}

// IsAdmin checks if the user is an admin
func IsAdmin(ctx context.Context) (*models.JWTClaims, error) {
	claims, err := IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if claims.Role != "admin" {
		return nil, errors.NewForbiddenError("You are not authorized to access this resource")
	}
	return claims, nil
}
