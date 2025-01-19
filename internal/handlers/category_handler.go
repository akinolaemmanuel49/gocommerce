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

func NewCategoryHandler(categoryService *services.CategoryService, logger, errorLogger *log.Logger) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, logger: logger, errorLogger: errorLogger}
}

// Compile-time check that CategoryHandler implements HandlerInterface
var _ ICategoryHandler = (*CategoryHandler)(nil)

// Create handles POST /categories requests and accepts CreateCategory as input
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Initialize request body
	var input models.CreateCategory

	// Parse request body
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create category
	category, err := h.categoryService.CreateCategory(ctx, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusCreated, category, h.logger)
}

// Read handles GET /categories/:id requests
func (h *CategoryHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r)

	// Validate ID
	if err := utils.ValidateID(ID, "Category"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to get category by ID
	category, err := h.categoryService.RetrieveCategoryByID(ctx, ID)
	switch err {
	case nil:
		// No error continue execution
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Category", "ID", ID), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, category, h.logger)
}

// ReadAll handles GET /categories requests with optional filters
func (h *CategoryHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Parse query parameters for filters and pagination
	query := r.URL.Query()
	lastID := query.Get("lastID")
	limitStr := query.Get("limit")

	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Build filter map from query parameters
	filter := map[string]interface{}{}
	if name := query.Get("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// Call service to read all categories
	categories, nextCursor, err := h.categoryService.RetrieveAllCategories(ctx, filter, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"data":       categories,
		"nextCursor": nextCursor,
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PUT /categories/:id requests
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r)

	// Validate the ID
	if err := utils.ValidateID(ID, "Category"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Initialize request body
	var req models.UpdateCategory

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update category
	user, err := h.categoryService.UpdateCategoryByID(ctx, ID, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, user, h.logger)
}

// Delete handles DELETE /categories/:id requests
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r)

	// Validate the ID
	if err := utils.ValidateID(ID, "Category"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to soft-delete user
	err = h.categoryService.DeleteCategoryByID(ctx, ID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": fmt.Sprintf("Category with ID: %s was successfully deleted", ID),
	}
	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
