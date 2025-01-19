package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	custom_errors "github.com/akinolaemmanuel49/gocommerce/common/errors"
	auth_middlewares "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
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

// IsAuthorized checks if the user is authorized to access a protected route
func IsAuthorized(ctx context.Context) (*models.JWTClaims, error) {
	claims := auth_middlewares.GetClaims(ctx)
	if claims == nil {
		return nil, custom_errors.NewUnauthorizedError()
	}
	return claims, nil
}

func IsAdmin(ctx context.Context) (*models.JWTClaims, error) {
	claims, err := IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if claims.Role != "admin" {
		return nil, custom_errors.NewForbiddenError()
	}
	return claims, nil
}
