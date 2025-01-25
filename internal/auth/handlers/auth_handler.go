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

// Compile-time check that AuthHandler implements IAuthHandler
var _ IAuthHandler = (*AuthHandler)(nil)

// Login handles POST /auth/login requests and accepts UserCredentials as input
// @Summary Login
// @Description Returns a JWT token for a verified user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "Login credentials"
// @Success 200 {object} models.Token "JWT token"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /auth/login [post]
func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.UserCredentials

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewBadRequestError("Check credentials and try again"), h.errorLogger)
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
