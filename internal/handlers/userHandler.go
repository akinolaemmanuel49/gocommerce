package handlers

import (
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUser handles GET /users requests
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters for pagination
	query := r.URL.Query()
	lastID := query.Get("lastID")
	limitStr := query.Get("limit")

	limit := 10 // Default value
	if limitStr == "" {
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

	users, nextCursor, err := h.userService.GetAllUsers(ctx, filter, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      users,
		"nextCusor": nextCursor,
	}
	writeJSON(w, http.StatusOK, response)
}
