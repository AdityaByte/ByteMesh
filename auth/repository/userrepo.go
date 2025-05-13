package repository

import (
	"context"
	"fmt"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/auth/model"
	"github.com/AdityaByte/bytemesh/auth/utils"
	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Function for checking the username exists or not.
func CheckUserExists(username string, repo *config.MongoRepository) error {
	// Here we have to check that the username exists or not
	var existingUser model.User

	err := repo.Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&existingUser)

	if err == nil {
		return fmt.Errorf("ERROR: Username already exists! Try another one")
	} else if err != mongo.ErrNoDocuments {
		return err
	} 

	return nil // No user exists for that username.
}

func ensureIndexes(repo *config.MongoRepository) error {
	_, err := repo.Collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}


func CreateUser(user *model.User, repo *config.MongoRepository) error {
	// For that we have to hash the password and other things too.
	encryptedPassword, err := utils.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	n, err := repo.Collection.InsertOne(context.TODO(), model.User{
		Username: user.Username,
		Password: encryptedPassword,
	})

	// While creating the username we have to associate an index to it also.
	if err = ensureIndexes(repo); err != nil {
		return fmt.Errorf("ERROR: Indexing Failed %v", err)
	}

	if err != nil {
		return fmt.Errorf("ERROR: Failed to insert the document %v", err)
	}

	logger.InfoLogger.Println("Item Inserted", n)
	return nil
}

func FindUser(user *model.User, repo *config.MongoRepository) error {
	var dbuser model.User
	repo.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&dbuser)

	// Here we have to compare the hashed password with the password
	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(user.Password)); err != nil {
		return fmt.Errorf("ERROR: Invalid Credentials")
	}
	return nil
}