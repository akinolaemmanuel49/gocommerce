package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IDatabase interface defines the methods of the Database struct
type IDatabase interface {
	Connect(uri string, dbName string) (*mongo.Database, error)
}

// Database defines a structure for the database instance
type Database struct {
	database *mongo.Database
}

// Connect method for Database provides implementation for the database connection
func (db *Database) Connect(uri string, dbName string) (*mongo.Database, error) {
	// SetServerAPIOptions method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		// SetMaxPoolSize(1000). // Max pool size for free tier is 100
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

// NewDatabase creates a new instance of Database
func NewDatabase() IDatabase {
	return &Database{}
}
