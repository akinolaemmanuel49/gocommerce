package services

import (
	"context"
	"fmt"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepository *repositories.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(productRepository *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

// GetAllProducts retrieves paginated products with optional filters
func (s *ProductService) GetAllProducts(ctx context.Context, filter map[string]interface{}, lastId string, limit int) ([]models.Product, string, error) {
	products, nextCursor, err := s.productRepository.FindAll(ctx, filter, lastId, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching products: %w", err)
	}

	return products, nextCursor, nil
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, ID string) (*models.Product, error) {
	product, err := s.productRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching product: %w", err)
	}

	return product, nil
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
		return nil, fmt.Errorf("error creating new product: %w", err)
	}

	// Convert result.InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	product.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return product, nil
}
