package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"os"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestDatabase holds the MongoDB client and database reference
type TestDatabase struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// SetupTestDatabase sets up a MongoDB test database
func SetupTestDatabase(uri, dbName string) (*TestDatabase, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return the test database
	return &TestDatabase{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// TeardownTestDatabase drops the test database and closes the connection
func (tdb *TestDatabase) TeardownTestDatabase() error {
	if tdb.Database != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := tdb.Database.Drop(ctx); err != nil {
			return err
		}
	}

	if tdb.Client != nil {
		return tdb.Client.Disconnect(context.Background())
	}

	return nil
}

// SeedDatabase seeds a collection with data from a JSON file
func (tdb *TestDatabase) SeedDatabase(collectionName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var docs []interface{}
	if err := json.Unmarshal(data, &docs); err != nil {
		return err
	}

	collection := tdb.Database.Collection(collectionName)
	_, err = collection.InsertMany(context.Background(), docs)
	return err
}

func GetRequestAuthenticated(t *testing.T, jwtSecretKey []byte, userID string, role string, req *http.Request) *http.Request {
	validTokenString, _ := utils.GenerateJWT(jwtSecretKey, userID, role)

	type contextKey string
	const UserClaimsKey contextKey = "userClaims"

	req.Header.Set("Authorization", "Bearer "+validTokenString)
	req.Header.Set("Content-Type", "application/json")

	// Parse and validate the token
	claims, err := utils.ParseJWT(jwtSecretKey, validTokenString)
	if err != nil {
		t.Fatalf("Failed to extract claims from token")
	}

	ctx := context.WithValue(req.Context(), UserClaimsKey, claims)
	req = req.WithContext(ctx)
	return req
}
