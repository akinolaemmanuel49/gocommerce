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

func NewProductHandler(productService *services.ProductService, logger, errorLogger *log.Logger) *ProductHandler {
	return &ProductHandler{productService: productService, logger: logger, errorLogger: errorLogger}
}

// Compile-time check that ProductHandler implements HandlerInterface
var _ IProductHandler = (*ProductHandler)(nil)

// Create handles POST /products requests and accepts CreateProduct as input
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Initialize request body
	var input models.CreateProduct

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to create product
	product, err := h.productService.CreateProduct(ctx, &input)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusCreated, product, h.logger)
}

// Read handles GET /products/:id requests
func (h *ProductHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
	// Validate the ID
	if err := utils.ValidateID(id, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

	// Call service to get product by ID
	product, err := h.productService.RetrieveProductByID(r.Context(), id)
	switch err {
	case nil:
		// No error, proceed
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Product", "ID", id), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the product data
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// ReadAll handles GET /products requests with optional filters
func (h *ProductHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
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

	// Build filter map
	filter := map[string]interface{}{}
	if name := query.Get("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if category := query.Get("category"); category != "" {
		filter["category"] = bson.M{"$regex": category, "$options": "i"}
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

	// Call service to read all products
	products, nextCursor, err := h.productService.RetrieveAllProducts(ctx, filter, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"data":       products,
		"nextCursor": nextCursor,
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PUT /products/:id requests
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Validate the ID
	if err := utils.ValidateID(id, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

	// Initialize request body
	var input models.UpdateProduct

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update product
	product, err := h.productService.UpdateProductByID(r.Context(), id, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// Delete handles DELETE /products/:id requests
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Validate the ID
	if err := utils.ValidateID(id, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

	// Call service to delete product
	if err := h.productService.DeleteProductByID(r.Context(), id); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": fmt.Sprintf("Product with ID: %s was successfully deleted", id),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
