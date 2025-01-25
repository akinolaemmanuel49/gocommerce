package repositories

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CategoryRepository represents data layer access for category collections
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db.Collection("categories")),
	}
}

// FindByID retrieves a category by its ID
func (r *CategoryRepository) FindByID(ctx context.Context, ID string) (*models.Category, error) {
	objID, err := utils.StringToObjectID(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "isDeleted": false}
	var category models.Category

	if err := r.Collection.FindOne(ctx, filter).Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &category, nil
}

// FindAll retrieves categories based on filters and implements cursor-based pagination
func (r *CategoryRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Category, string, error) {
	query := bson.M{}
	filter["isDeleted"] = false
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objID, err := utils.StringToObjectID(lastID)
		if err != nil {
			return nil, "", errors.NewValidationError("nextCursor", "must be a valid ObjectID")
		}
		query["_id"] = bson.M{"$gt": objID} // Fetch categories with IDs greater than lastID
	}

	options := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"_id": 1}) // Sort by _id in ascending order

	cursor, err := r.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	// Decode categories
	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, "", err
	}

	// Determine the next cursor (last _id in the result)
	var nextCursor string
	if len(categories) > 0 {
		nextCursor = categories[len(categories)-1].ID // Assuming models.Category.ID corresponds to the MongoDB `_id`
	}

	return categories, nextCursor, nil
}
