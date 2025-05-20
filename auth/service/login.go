package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/auth/model"
	"github.com/AdityaByte/bytemesh/auth/repository"
	"github.com/AdityaByte/bytemesh/datanodes/server3/logger"
)

type token struct {
	Token string `json:token`
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content.Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	var user model.User

	json.NewDecoder(r.Body).Decode(&user)

	repo, err := config.LoadMongoRepository()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if err := repository.FindUser(&user, repo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// When the credentials are valid we have to generate a JWT Token which validates the each upcoming requests such as upload and download.

	// Creating the token.
	token, err := config.CreateToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	logger.InfoLogger.Println("Token has been generated.")

	// Sending the token back to the client with header.
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Valid credentials")
}
