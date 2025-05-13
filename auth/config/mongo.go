package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func LoadMongoRepository() (*MongoRepository, error) {
	// Logic: for connecting to the mongodb instance.
	var (
		uri            = os.Getenv("MONGO_URI")
		dbName         = os.Getenv("DB_NAME")
		collectionName = os.Getenv("COLLECTION")
	)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to connect to the mongodb instance %s", err)
	}

	// Here we have to ping the instance for checking the connection is established or not.
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("ERROR: Failed to ping %s",err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoRepository{
		Client: client,
		Collection: collection,
	}, nil
}
