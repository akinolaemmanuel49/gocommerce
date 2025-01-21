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

// Create handles POST /categories requests and accepts CreateCategory as input [ADMIN]
// @Security BearerAuth
// @Summary Create a new category.
// @Description This endpoint creates a new category, this is an admin only endpoint.
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategory true "Category Details"
// @Success 201 {object} models.Category "Created category"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Not Found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /categories [post]
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

// Read handles GET /categories/:id requests [PUBLIC]
// @Summary Read a category
// @Description This endpoint fetches a single category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category "Returned category"
// @Failure 400 "Invalid Category ID"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /categories [get]
func (h *CategoryHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r)

	// Validate the ID
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

// ReadAll handles GET /categories requests with optional filters [PUBLIC]
// @Summary Read all categories
// @Description This endpoint fetches a list of categories with cursor based pagination, optionally filtered by name.
// @Tags Categories
// @Accept json
// @Produce json
// @Param name query string false "Filter products by name"
// @Param lastID query string false "Last category id in a page"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.MultipleEntityClientResponse "Returned products and next cursor"
// @Failure 500 "Internal Server Error"
// @Router /categories/all [get]
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

// Update handles PUT /categories/:id requests [ADMIN]
// @Security BearerAuth
// @Summary Update category
// @Description This endpoint updates a single category, this is an admin only endpoint.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string false "Category ID"
// @Success 200 {object} models.Category "Updated category"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /categories [put]
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

// Delete handles DELETE /categories/:id requests [ADMIN]
// @Security BearerAuth
// @Summary Delete category
// @Description This endpoint deletes a single category, this is an admin only endpoint.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string false "Category ID"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /categories [delete]
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
