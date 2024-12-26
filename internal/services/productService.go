package services

import (
	"context"
	"fmt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
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
