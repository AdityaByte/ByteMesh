package service

import (
	"context"
	"fmt"
	"namenodeserver/database"
	"namenodeserver/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func pushMetaData(ctx context.Context, metadata *model.MetaData, mongoRepo *database.MongoRepository) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	result, err := mongoRepo.Collection.InsertOne(ctx, metadata)

	if err != nil {
		return err
	}

	fmt.Println("Items inserted to db:", result.InsertedID)
	return nil
}

func fetchMetaData(ctx context.Context, filename string, mongoRepo *database.MongoRepository) (*model.MetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second * 3)
	defer cancel()

	filter := bson.M{"filename" : filename}

	var metaData model.MetaData // Created an instance of metadata for storing it.

	err := mongoRepo.Collection.FindOne(ctx, filter).Decode(&metaData)

	if err != nil {
		return nil, err
	}

	fmt.Println("Meta data found:", metaData)

	return &metaData, nil
}
