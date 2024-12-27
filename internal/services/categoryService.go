package services

import (
	"context"
	"fmt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(categoryRepository *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepository: categoryRepository}
}

// GetAllCategories retrieves paginated categories with optional filters
func (s *CategoryService) GetAllCategories(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Category, string, error) {
	categories, nextCursor, err := s.categoryRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching categories: %w", err)
	}

	return categories, nextCursor, nil
}
