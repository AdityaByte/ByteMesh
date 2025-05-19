package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/auth/model"
	"github.com/AdityaByte/bytemesh/auth/repository"
	"github.com/AdityaByte/bytemesh/datanodes/server3/logger"
)

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

	token, err := config.CreateToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if err := os.MkdirAll("../.auth/", os.ModePerm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create the directory %v", err)
		return
	}

	file, err := os.Create("../.auth/.jwt_token")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create the file %v", err)
		return
	}

	if _, err := file.WriteString(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to write in the string %v", err)
		return
	}

	logger.InfoLogger.Println("Token has been created to .auth")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Valid credentials")
}
