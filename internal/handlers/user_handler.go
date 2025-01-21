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
// @Summary Create a new user
// @Description This endpoint creates a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User Details"
// @Success 201 {object} models.User "Created user"
// @Failure 400 "Invalid Request Body"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /users [post]
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

	// Write response to client
	utils.WriteJSON(w, r, http.StatusCreated, user, h.logger)
}

// Read handles GET /user requests [CUSTOMER]
// @Security BearerAuth
// @Summary Read a user
// @Description This endpoint fetches a single user, optionally fetch a user by id if the user is an admin.
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string false "User ID"
// @Success 200 {object} models.User "Returned user"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /users [get]
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
	if id == "" || id == claims.UserID {
		userIDToRead = claims.UserID // Read id from claims if id query is empty
	} else if id != "" && claims.Role == "admin" {
		userIDToRead = id // Admin can read any user
	} else {
		errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to read this user"), h.errorLogger) // Only admin can read any user
		return
	}

	// DEBUG
	fmt.Println("DEBUG!!!")
	fmt.Println(id)
	fmt.Println("READING")
	fmt.Println(userIDToRead)

	// Call service to read user by ID
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

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// ReadAll handles GET /users requests [ADMIN]
// @Security BearerAuth
// @Summary Read all users
// @Description This endpoint fetches a list of users with cursor based pagination, optionally filtered by firstName, lastName, email, country, state, role
// @Tags Users
// @Accept json
// @Produce json
// @Param firstName query string false "Filter users by first name"
// @Param lastName query string false "Filter users by last name"
// @Param email query string false "Filter users by email"
// @Param country query string false "Filter users by country"
// @Param state query string false "Filter users by state"
// @Param role query string false "Filter users by role"
// @Param lastID query string false "Last user id in a page"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.MultipleEntityClientResponse "Returned users and next cursor"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /users/all [get]
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

	// Call service to read all users
	users, nextCursor, err := h.userService.RetrieveAllUsers(ctx, filter, lastID, limit)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.MultipleEntityClientResponse{
		Data:       users,
		NextCursor: nextCursor,
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PUT /users requests [CUSTOMER]
// @Security BearerAuth
// @Summary Update user
// @Description This endpoint updates a single user, optionally update a user by id if the user is an admin.
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string false "User ID"
// @Success 200 {object} models.User "Updated user"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /users [put]
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
	if id == "" || id == claims.UserID {
		userIDToUpdate = claims.UserID // Read id from claims if id query is empty
	} else if id != "" && claims.Role == "admin" {
		userIDToUpdate = id // Admin can update any user
	} else {
		errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to update this user"), h.errorLogger) // Only admin can update any user
		return
	}

	// Call service to update user
	user, err := h.userService.UpdateUserByID(ctx, userIDToUpdate, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// Delete handles DELETE /users [CUSTOMER]
// @Security BearerAuth
// @Summary Delete user
// @Description This endpoint deletes a single user, optionally delete a user by id if the user is an admin.
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string false "User ID"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /users [delete]
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
	if id == "" || id == claims.UserID {
		userIDToDelete = claims.UserID // Read id from claims if id query is empty
	} else if id != "" && claims.Role == "admin" {
		userIDToDelete = id // Admin can delete any user
	} else {
		errors.HandleError(w, r, errors.NewForbiddenError("You are not authorized to delete this user"), h.errorLogger) // Only admin can delete any user
		return
	}

	// Call service to soft-delete user
	err = h.userService.DeleteUserByID(ctx, userIDToDelete)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: fmt.Sprintf("User with ID: %s was successfully deleted", userIDToDelete),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
