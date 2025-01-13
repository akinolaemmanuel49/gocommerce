package services

import (
	"context"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(categoryRepository *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepository: categoryRepository}
}

// CreateCategory creates a new instance of a category and commits it to the database
func (s *CategoryService) CreateCategory(ctx context.Context, newCategory *models.CreateCategory) (*models.Category, error) {
	// Transform CreateCategory to Category
	category := models.NewCategory(newCategory)

	// Insert category into the database
	result, err := s.categoryRepository.Insert(ctx, category)
	if err != nil {
		return nil, err
	}

	// Convert InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.NewValidationError("ID", "failed to convert ObjectID to string")
	}
	category.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return category, nil
}

// RetrieveCategoryByID retrieves a category by ID
func (s *CategoryService) RetrieveCategoryByID(ctx context.Context, ID string) (*models.Category, error) {
	category, err := s.categoryRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// RetrieveAllCategories retrieves paginated categories with optional filters
func (s *CategoryService) RetrieveAllCategories(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Category, string, error) {
	categories, nextCursor, err := s.categoryRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", err
	}

	return categories, nextCursor, nil
}

// UpdateCategoryByID updates an instance of a category and commits it to the database
func (s *CategoryService) UpdateCategoryByID(ctx context.Context, ID string, updatedCategory *models.UpdateCategory) (*models.Category, error) {
	// Check for existing category
	existingCategory, err := s.categoryRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	// Transform UpdateCategory to Category
	category := models.CategoryUpdate(updatedCategory, existingCategory)

	update := bson.M{
		"$set": category,
	}

	// Update category in database
	_, err = s.categoryRepository.Update(ctx, ID, update)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategoryByID sets the IsDeleted flag for a category instance to true (performs a soft-delete)
func (s *CategoryService) DeleteCategoryByID(ctx context.Context, ID string) error {
	// Check for existing category
	existingCategory, err := s.categoryRepository.FindByID(ctx, ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	if existingCategory != nil {
		// Apply transformation, set category IsDeleted field to true
		category := &models.Category{
			CommonFields: models.CommonFields{
				IsDeleted: true,
				UpdatedAt: time.Now(),
			},
		}

		_, err = s.categoryRepository.Update(ctx, ID, category)
		if err != nil {
			return err
		}
	}

	return nil
}
