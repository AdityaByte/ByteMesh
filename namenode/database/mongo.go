package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func LoadMongoRepository() (*MongoRepository, error) {

	var (
		uri = os.Getenv("MONGO_URI")
		dbName = os.Getenv("DB_NAME")
		collectionName = os.Getenv("COLLECTION")
	)

	if uri == "" || dbName == "" || collectionName == "" {
		return nil, fmt.Errorf("Environment variables are empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(50).SetMinPoolSize(5).SetServerSelectionTimeout(30 * time.Second).SetSocketTimeout(60 * time.Second).SetConnectTimeout(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to mongodb %w", err)
	}

	// Verifying the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoRepository{
		Client: client,
		Collection: collection,
	} , nil
}
