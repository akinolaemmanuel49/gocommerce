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
var _ HandlerInterface = (*UserHandler)(nil)

// Create handles POST /user requests
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create user
	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the created user
	utils.WriteJSON(w, r, http.StatusCreated, user, h.logger)
}

// Read handles GET /user/:id requests
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if err := utils.ValidateID(ID, "User"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to get user by ID
	user, err := h.userService.RetrieveUserByID(r.Context(), ID)
	switch err {
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("User", "ID", ID), h.errorLogger)
		return
	case nil:
		// No error continue execution
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the user data
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// ReadAll handles GET /users requests
func (h *UserHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Log to stdout
	h.logger.Printf("%v %v", r.Method, r.URL.Path)

	ctx := r.Context()

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

// Update handles PATCH /users/:id requests
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if err := utils.ValidateID(ID, "User"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}
	var req models.UpdateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUserByID(r.Context(), ID, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with updated user
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// Delete handles PATCH /users/:id/delete requests
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if err := utils.ValidateID(ID, "User"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to soft-delete user
	err := h.userService.DeleteUserByID(r.Context(), ID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with confirmation of deletion
	response := map[string]string{"message": fmt.Sprintf("User with ID: %s was successfully deleted", ID)}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
