package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/utils"
)

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// Compile-time check that UserHandler implements HandlerInterface
var _ HandlerInterface = (*UserHandler)(nil)

// Create handles POST /user requests
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("*", "Invalid request body"))
		return
	}

	// Call service to create user
	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// Respond with the created user
	utils.WriteJSON(w, r, http.StatusCreated, user)
}

// Read handles GET /user/:id requests
func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if ID == "" {
		errors.HandleError(w, r, errors.NewValidationError("ID", "User ID is required"))
		return
	}

	// Call service to get user by ID
	user, err := h.userService.RetrieveUserByID(r.Context(), ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// Respond with the user data
	utils.WriteJSON(w, r, http.StatusOK, user)
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
		errors.HandleError(w, r, err)
		return
	}

	response := map[string]interface{}{
		"data":      users,
		"nextCusor": nextCursor,
	}
	utils.WriteJSON(w, r, http.StatusOK, response)
}

// Update handles PATCH /users/:id requests
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if ID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	var req models.UpdateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("*", "Invalid request body"))
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUserByID(r.Context(), ID, &req)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	// Respond with updated user
	utils.WriteJSON(w, r, http.StatusOK, user)
}

// Delete handles PATCH /users/:id/delete requests
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented") // TODO
}
