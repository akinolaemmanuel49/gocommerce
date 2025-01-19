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
	"go.mongodb.org/mongo-driver/bson"
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
	// Initialize context
	ctx := r.Context()

	// Initialize request body
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

	// Parse query parameters
	query := r.URL.Query()
	id := query.Get("id")

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Check if user is trying to read another user, and if so, check if user is an admin
	var userIDToRead string
	if id != claims.UserID {
		if claims.Role != "admin" {
			errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to update this user"), h.errorLogger)
			return
		} else if id != "" && claims.Role == "admin" {
			userIDToRead = id // Admin can read any user
		} else {
			userIDToRead = claims.UserID // Customer can only read self
		}
	}

	// Call service to get user by ID
	user, err := h.userService.RetrieveUserByID(ctx, userIDToRead)
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

	// Build filter map from query parameters
	filter := map[string]interface{}{}
	if firstName := query.Get("firstName"); firstName != "" {
		filter["firstName"] = bson.M{"$regex": firstName, "$options": "i"}
	}
	if lastName := query.Get("lastName"); lastName != "" {
		filter["lastName"] = bson.M{"$regex": lastName, "$options": "i"}
	}
	if email := query.Get("email"); email != "" {
		filter["email"] = bson.M{"$regex": email, "$options": "i"}
	}
	if country := query.Get("country"); country != "" {
		filter["address.country"] = bson.M{"$regex": country, "$options": "i"}
	}
	if state := query.Get("state"); state != "" {
		filter["address.state"] = bson.M{"$regex": state, "$options": "i"}
	}
	if role := query.Get("role"); role != "" {
		filter["role"] = role
	}

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

	// Parse query parameters
	query := r.URL.Query()
	id := query.Get("id")

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

	// Check if user is trying to update another user, and if so, check if user is an admin
	var userIDToUpdate string
	if id != claims.UserID {
		if claims.Role != "admin" {
			errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to update this user"), h.errorLogger)
			return
		} else if id != "" && claims.Role == "admin" {
			userIDToUpdate = id // Admin can update any user
		} else {
			userIDToUpdate = claims.UserID // Customer can only update self
		}
	}

	// Call service to update user
	user, err := h.userService.UpdateUserByID(ctx, userIDToUpdate, &req)
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

	// Parse query parameters
	query := r.URL.Query()
	id := query.Get("id")

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Check if user is trying to delete another user, and if so, check if user is an admin
	var userIDToDelete string
	if id != claims.UserID {
		if claims.Role != "admin" {
			errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to delete this user"), h.errorLogger)
			return
		} else if id != "" && claims.Role == "admin" {
			userIDToDelete = id // Admin can delete any user
		} else {
			userIDToDelete = claims.UserID // Customer can only delete self
		}
	}

	// Call service to soft-delete user
	err = h.userService.DeleteUserByID(ctx, userIDToDelete)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with confirmation of deletion
	response := map[string]string{"message": fmt.Sprintf("User with ID: %s was successfully deleted", userIDToDelete)}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
