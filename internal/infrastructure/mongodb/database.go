package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents a MongoDB database connection
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewDatabase creates a new MongoDB database connection
func NewDatabase(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Database{
		Client:   client,
		Database: db,
	}, nil
}

// Close closes the MongoDB connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return d.Client.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.Database.Collection(name)
}

// Connect creates a new MongoDB connection
func Connect(ctx context.Context, uri, dbName string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Database{
		Client:   client,
		Database: db,
	}, nil
}

// Disconnect closes the MongoDB connection
func (d *Database) Disconnect(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}
