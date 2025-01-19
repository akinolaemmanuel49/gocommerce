package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService, logger, errorLogger *log.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger, errorLogger: errorLogger}
}

// Compile-time check that UserHandler implements HandlerInterface
var _ IUserHandler = (*UserHandler)(nil)

// Create handles POST /user requests [PUBLIC]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create user
	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the created user
	utils.WriteJSON(w, r, http.StatusCreated, user, h.logger)
}

// Read handles GET /user requests [CUSTOMER]
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to get user by ID
	user, err := h.userService.RetrieveUserByID(ctx, claims.UserID)
	switch err {
	case nil:
		// No error continue execution
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("User", "ID", claims.UserID), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the user data
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// ReadAll handles GET /users requests [ADMIN]
func (h *UserHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Parse query parameters for pagination
	query := r.URL.Query()
	lastID := query.Get("lastID")
	limitStr := query.Get("limit")

	limit := 10 // Default value
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// TODO:: more filters
	// Build filter map
	filter := map[string]interface{}{}
	if firstName := query.Get("firstName"); firstName != "" {
		filter["firstName"] = firstName
	}
	if lastName := query.Get("lastName"); lastName != "" {
		filter["lastName"] = lastName
	}
	if email := query.Get("email"); email != "" {
		filter["email"] = email
	}
	// if country := query.Get("country"); country != "" {
	// 	filter["country"] = country
	// }

	users, nextCursor, err := h.userService.RetrieveAllUsers(ctx, filter, lastID, limit)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	response := map[string]interface{}{
		"data":      users,
		"nextCusor": nextCursor,
	}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PUT /users requests [CUSTOMER]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Initialize request body
	var req models.UpdateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUserByID(ctx, claims.UserID, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with updated user
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// Delete handles DELETE /users [CUSTOMER]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to soft-delete user
	err = h.userService.DeleteUserByID(ctx, claims.UserID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with confirmation of deletion
	response := map[string]string{"message": fmt.Sprintf("User with ID: %s was successfully deleted", claims.UserID)}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
