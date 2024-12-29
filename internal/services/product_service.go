package services

import (
	"context"
	"log"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewProductService creates a new instance of ProductService
func NewProductService(productRepository *repositories.ProductRepository, logger *log.Logger) *ProductService {
	return &ProductService{productRepository: productRepository, logger: logger}
}

// CreateProduct creates a new instance of a product and commits it to the database
func (s *ProductService) CreateProduct(ctx context.Context, newProduct *models.CreateProduct) (*models.Product, error) {
	// Transform CreateProduct to Product
	product := &models.Product{
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Images:      newProduct.Images,
		CategoryID:  newProduct.CategoryID,
		Brand:       newProduct.Brand,
		CommonFields: models.CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert product into the database
	result, err := s.productRepository.Insert(ctx, product)
	if err != nil {
		return nil, err
	}

	// Convert result.InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	product.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return product, nil
}

// RetrieveProductByID retrieves a product by its ID
func (s *ProductService) RetrieveProductByID(ctx context.Context, ID string) (*models.Product, error) {
	product, err := s.productRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// RetrieveAllProducts retrieves paginated products with optional filters
func (s *ProductService) RetrieveAllProducts(ctx context.Context, filter map[string]interface{}, lastId string, limit int) ([]models.Product, string, error) {
	products, nextCursor, err := s.productRepository.FindAll(ctx, filter, lastId, limit)
	if err != nil {
		return nil, "", err
	}

	return products, nextCursor, nil
}

// UpdateProductByID updates an instance of a product and commits it to the database
func (s *ProductService) UpdateProductByID(ctx context.Context, ID string, updatedProduct *models.UpdateProduct) (*models.Product, error) {
	// Check for existing product
	existingProduct, err := s.productRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	// Transform UpdateProduct to Product
	product := &models.Product{
		Name:        utils.IfNotNil(updatedProduct.Name, existingProduct.Name),
		Description: utils.IfNotNil(updatedProduct.Description, existingProduct.Description),
		Price:       utils.IfNotNil(updatedProduct.Price, existingProduct.Price),
		Brand:       utils.IfNotNil(updatedProduct.Brand, existingProduct.Brand),
		CommonFields: models.CommonFields{
			UpdatedAt: time.Now(),
		},
	}

	_, err = s.productRepository.Update(ctx, ID, product)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	updatedRecord, err := s.productRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return updatedRecord, nil
}

// DeleteProductByID sets the IsDeleted flag for a product instance to true (performs a soft-delete)
func (s *ProductService) DeleteProductByID(ctx context.Context, ID string) error {
	// Check for existing product
	existingProduct, err := s.productRepository.FindByID(ctx, ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	// Set existingProduct IsDeleted field to true
	existingProduct.IsDeleted = true
	existingProduct.UpdatedAt = time.Now()
	_, err = s.productRepository.Update(ctx, ID, existingProduct)
	if err != nil {
		return err
	}

	return nil
}
