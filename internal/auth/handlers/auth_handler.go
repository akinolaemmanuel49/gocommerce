package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/services"
	"github.com/akinolaemmanuel49/gocommerce/utils"
)

// NewAuthHandler creates a new instance of an AuthHandler
func NewAuthHandler(authService *services.AuthService, logger *log.Logger, errorLogger *log.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger, errorLogger: errorLogger}
}

var _ IAuthHandler = (*AuthHandler)(nil)

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.UserCredentials

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to generate a JWT token
	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, token, h.logger)
}
