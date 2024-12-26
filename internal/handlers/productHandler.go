package handlers

import (
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetAllProducts handles GET /products requests with optional filters
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters for filters and pagination
	query := r.URL.Query()
	lastID := query.Get("last_id")
	limitStr := query.Get("limit")

	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Build filter map
	filter := map[string]interface{}{}
	if category := query.Get("category"); category != "" {
		filter["category"] = category
	}
	if priceMin := query.Get("priceMin"); priceMin != "" {
		if min, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filter["price"] = map[string]interface{}{"$gte": min}
		}
	}
	if priceMax := query.Get("priceMax"); priceMax != "" {
		if max, err := strconv.ParseFloat(priceMax, 64); err == nil {
			if priceFilter, exists := filter["price"].(map[string]interface{}); exists {
				priceFilter["$lte"] = max
			} else {
				filter["price"] = map[string]interface{}{"$lte": max}
			}
		}
	}

	products, nextCursor, err := h.productService.GetAllProducts(ctx, filter, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       products,
		"nextCursor": nextCursor,
	}
	writeJSON(w, http.StatusOK, response)
}
