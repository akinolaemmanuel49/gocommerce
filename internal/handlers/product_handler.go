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

// Create handles POST /products requests and accepts CreateProduct as input [ADMIN]
// @Security BearerAuth
// @Summary Create a new product.
// @Description This endpoint creates a new product, this is an admin only endpoint.
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.CreateProduct true "Product Details"
// @Success 201 {object} models.Product "Created product"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Not Found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /products [post]
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

// Read handles GET /products/:id requests [PUBLIC]
// @Summary Read a product
// @Description This endpoint fetches a single product.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product "Returned product"
// @Failure 400 "Invalid Product ID"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /products [get]
func (h *ProductHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r)

	// Validate the ID
	if err := utils.ValidateID(ID, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

	// Call service to get product by ID
	product, err := h.productService.RetrieveProductByID(ctx, ID)
	switch err {
	case nil:
		// No error, proceed
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Product", "ID", ID), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// ReadAll handles GET /products requests with optional filters [PUBLIC]
// @Summary Read all products
// @Description This endpoint fetches a list of products with cursor based pagination, optionally filtered by name, category, priceMin, priceMax
// @Tags Products
// @Accept json
// @Produce json
// @Param name query string false "Filter products by name"
// @Param category query string false "Filter products by category"
// @Param priceMin query float64 false "Filter products by a set minimum price"
// @Param priceMax query float64 false "Filter products by a set maximum price"
// @Param lastID query string false "Last product id in a page"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.MultipleEntityClientResponse "Returned products and next cursor"
// @Failure 500 "Internal Server Error"
// @Router /products/all [get]
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

// Update handles PUT /products/:id requests [ADMIN]
// @Security BearerAuth
// @Summary Update product
// @Description This endpoint updates a single product, this is an admin only endpoint.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string false "Product ID"
// @Success 200 {object} models.Product "Updated product"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /products [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	if err := utils.ValidateID(ID, "Product"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
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
	product, err := h.productService.UpdateProductByID(ctx, ID, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, product, h.logger)
}

// Delete handles DELETE /products/:id requests [ADMIN]
// @Security BearerAuth
// @Summary Delete product
// @Description This endpoint deletes a single product, this is an admin only endpoint.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string false "Product ID"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /products [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	if err := utils.ValidateID(ID, "Product"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to delete product
	if err := h.productService.DeleteProductByID(ctx, ID); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": fmt.Sprintf("Product with ID: %s was successfully deleted", ID),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
