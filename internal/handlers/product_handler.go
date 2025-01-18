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
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProductHandler(productService *services.ProductService, logger, errorLogger *log.Logger) *ProductHandler {
	return &ProductHandler{productService: productService, logger: logger, errorLogger: errorLogger}
}

// Compile-time check that ProductHandler implements HandlerInterface
var _ IProductHandler = (*ProductHandler)(nil)

// Create handles POST /products requests and accepts CreateProduct as input
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateProduct

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to create product
	product, err := h.productService.CreateProduct(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Respond with the created product
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
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Product", "ID", id), h.errorLogger)
		return
	case nil:
		// No error, proceed
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the product data
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// ReadAll handles GET /products requests with optional filters
func (h *ProductHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
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

	products, nextCursor, err := h.productService.RetrieveAllProducts(ctx, filter, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       products,
		"nextCursor": nextCursor,
	}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PUT /products/:id requests
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	// Validate the ID
	if err := utils.ValidateID(id, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

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

	// Respond with the updated product
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// Delete handles DELETE /products/:id requests
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
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

	// Respond with confirmation of deletion
	response := map[string]string{"message": "Product successfully deleted"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
