package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// GetAllUser handles GET /users requests
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Log to stdout
	h.logger.Printf("%v %v", r.Method, r.URL.Path)

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
	// if country := query.Get("country"); country != "" {
	// 	filter["country"] = country
	// }

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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to create user
	user, err := h.userService.CreateUser(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if ID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Call service to get user by ID
	user, err := h.userService.GetUserByID(r.Context(), ID)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user exists
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user data
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, ID string) {
	// Validate the ID
	if ID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	var input models.UpdateUser

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUser(r.Context(), ID, &input)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with updated user
	writeJSON(w, http.StatusOK, user)
}
