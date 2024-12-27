package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(categoryRepository *repositories.CategoryRepository, logger *log.Logger) *CategoryService {
	return &CategoryService{categoryRepository: categoryRepository, logger: logger}
}

// GetAllCategories retrieves paginated categories with optional filters
func (s *CategoryService) GetAllCategories(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Category, string, error) {
	categories, nextCursor, err := s.categoryRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching categories: %w", err)
	}

	return categories, nextCursor, nil
}

// CreateCategory creates a new instance of a category and commits it to the database
func (s *CategoryService) CreateCategory(ctx context.Context, newCategory *models.CreateCategory) (*models.Category, error) {
	// Transform CreateCategory to Category
	category := &models.Category{
		Name:        newCategory.Name,
		Description: newCategory.Description,
		ParentID:    newCategory.ParentID,
		Image:       newCategory.Image,
		CommonFields: models.CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert category into the database
	result, err := s.categoryRepository.Insert(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("error creating new category: %w", err)
	}

	// Convert InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	category.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return category, nil
}
