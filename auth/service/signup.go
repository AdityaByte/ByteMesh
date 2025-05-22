package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/auth/model"
	"github.com/AdityaByte/bytemesh/auth/repository"
	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Sign up request..")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User

	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	// Here we have to Load the Mongo Repository
	repo, err := config.LoadMongoRepository()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.CheckUserExists(user.Username, repo); err != nil {
		logger.ErrorLogger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Else the username is unique so we have to save it.
	if err := repository.CreateUser(&user, repo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.InfoLogger.Println("Sign up successfull")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Signup successful")
}
