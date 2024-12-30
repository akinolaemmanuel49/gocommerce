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

type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db.Collection("users")),
	}
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, ID string) (*models.User, error) {
	objectID, err := utils.StringToObjectID(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var user models.User

	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User", "ID", ID)
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{"email": email}
	var user models.User

	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User", "email", email)
		}
		return nil, err
	}

	return &user, nil
}

// FindAll retrieves users based on filters and implements cursor-based pagination.
func (r *UserRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.User, string, error) {
	query := bson.M{}
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objectID, err := utils.StringToObjectID(lastID)
		if err != nil {
			return nil, "", err
		}
		query["_id"] = bson.M{"$gt": objectID} // Fetch users with IDs greater than lastID
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
