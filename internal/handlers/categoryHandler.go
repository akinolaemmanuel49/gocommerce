package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

func NewCategoryHandler(categoryService *services.CategoryService, logger *log.Logger) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, logger: logger}
}

// GetAllCategories handles GET /categories requests with optional filter
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Log to stdout
	h.logger.Printf("%v %v", r.Method, r.URL.Path)

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

	// Build filter map
	// TODO

	categories, nextCursor, err := h.categoryService.GetAllCategories(ctx, nil, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       categories,
		"nextCursor": nextCursor,
	}
	writeJSON(w, http.StatusOK, response)
}

// CreateCategory handles POST /categories requests and accepts CreateCategory as input
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCategory

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to create category
	category, err := h.categoryService.CreateCategory(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	// Respond with the created category
	writeJSON(w, http.StatusCreated, category)
}
