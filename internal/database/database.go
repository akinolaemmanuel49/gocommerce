package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the interface
type IDatabase interface {
	Connect(uri string, dbName string) (*mongo.Database, error)
}

// Struct implementing the interface
type Database struct {
	database *mongo.Database
}

// Implement the Connect method
func (db *Database) Connect(uri string, dbName string) (*mongo.Database, error) {
	// SetServerAPIOptions method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		// SetMaxPoolSize(500). // Max pool size for free tier is 100
		SetServerAPIOptions(serverAPI)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Assign the connected database to the struct
	db.database = client.Database(dbName)
	return db.database, nil
}

// Factory function to create a new Database
func NewDatabase() IDatabase {
	return &Database{}
}
