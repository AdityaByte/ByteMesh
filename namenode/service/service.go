package service

import (
	"context"
	"fmt"
	"namenodeserver/database"
	"namenodeserver/logger"
	"namenodeserver/model"

	"go.mongodb.org/mongo-driver/bson"
)

func PushMetaData(ctx context.Context, metadata *model.MetaData, mongoRepo *database.MongoRepository) error {

	result, err := mongoRepo.Collection.InsertOne(ctx, metadata)

	if err != nil {
		return err
	}

	logger.InfoLogger.Println("Items inserted to db:", result.InsertedID)
	return nil
}

func FetchMetaData(ctx context.Context, filename string, extension string, mongoRepo *database.MongoRepository) (*model.MetaData, error) {

	logger.InfoLogger.Println("FileName:", filename, "Extension:", extension)

	filter := bson.M{"filename": filename, "fileextension": extension}

	var metaData model.MetaData // Created an instance of metadata for storing it.

	err := mongoRepo.Collection.FindOne(ctx, filter).Decode(&metaData)

	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to find the document: %v", err)
	}

	logger.InfoLogger.Println("Meta data found:", metaData)

	return &metaData, nil
}

func FetchUserSpecificMetaData(ctx context.Context, user string, mongoRepo *database.MongoRepository) (*[]model.MetaData, error) {

	logger.InfoLogger.Println("Owner name:", user)

	filter := bson.M{"owner": user}

	var results []model.MetaData

	cursor, err := mongoRepo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to find the document: %v", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("ERROR: Didn't get the result: %v", err)
	}

	fmt.Println("result:", results)

	// Here i want to print out the result ok.
	for data := range results {
		fmt.Println(data)
	}

	return &results, nil
}
