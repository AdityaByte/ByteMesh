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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, fmt.Errorf("Error occured while connecting to the client", err)
	}

	defer client.Disconnect(ctx)

	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION"))

	return &MongoRepository{
		Client: client,
		Collection: collection,
	} , nil
}
