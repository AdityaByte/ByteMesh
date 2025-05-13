package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/auth/model"
	"github.com/AdityaByte/bytemesh/auth/repository"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid Request Type")
		return
	}

	var user model.User

	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	// Here we have to Load the Mongo Repository
	repo, err := config.LoadMongoRepository()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	if err := repository.CheckUserExists(user.Username, repo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Else the username is unique so we have to save it.
	if err := repository.CreateUser(&user, repo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Signup successful")
}
