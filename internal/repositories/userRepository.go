package repositories

import (
	"context"
	"fmt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db.Collection("users")),
	}
}

// FindAll retrieves users based on filters and implements cursor-based pagination.
func (r *UserRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.User, string, error) {
	query := bson.M{}
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid lastID: %v", err)
		}
		query["_id"] = bson.M{"$gt": objID} // Fetch users with IDs greater than lastID
	}

	options := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"_id": 1}) // Sort by _id in ascending order for cursor-based pagination

	cursor, err := r.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	// Decode users
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, "", err
	}

	// Determine the next cursor (last _id in the result)
	var nextCursor string
	if len(users) > 0 {
		nextCursor = users[len(users)-1].ID // Assuming models.User.ID corresponds to the MongoDB `_id`
	}

	return users, nextCursor, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}

	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, ID string) (*models.User, error) {
	var user models.User
	filter := bson.M{"_id": ID}

	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}
