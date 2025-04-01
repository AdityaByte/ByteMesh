package service

import (
	"context"
	"fmt"
	"namenodeserver/database"
	"namenodeserver/model"

	"go.mongodb.org/mongo-driver/bson"
)

func PushMetaData(ctx context.Context, metadata *model.MetaData, mongoRepo *database.MongoRepository) error {

	result, err := mongoRepo.Collection.InsertOne(ctx, metadata)

	if err != nil {
		return err
	}

	fmt.Println("Items inserted to db:", result.InsertedID)
	return nil
}

func FetchMetaData(ctx context.Context, filename string, extension string, mongoRepo *database.MongoRepository) (*model.MetaData, error) {

	fmt.Println("File name is :", filename)
	fmt.Println("File extension is :", extension)

	filter := bson.M{"filename": filename, "fileextension": extension}

	var metaData model.MetaData // Created an instance of metadata for storing it.

	err := mongoRepo.Collection.FindOne(ctx, filter).Decode(&metaData)

	if err != nil {
		return nil, err
	}

	fmt.Println("Meta data found:", metaData)

	return &metaData, nil
}
