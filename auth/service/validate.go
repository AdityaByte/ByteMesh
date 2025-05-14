package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaByte/bytemesh/auth/config"
	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
)

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "ERROR: Method not allowed")
		return
	}

	tokenString := r.Header.Get("Authorization")
	logger.InfoLogger.Println("Token String that the server gets :", tokenString)

	if strings.TrimSpace(tokenString) == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Missing Authorization header")
		return
	}

	// tokenString = tokenString[len("Bearer "):] // If we comments out this then in this case the error should be fixed out.

	logger.InfoLogger.Println("Token String that the server gets :", tokenString)

	if err := config.VerifyToken(tokenString); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid Token %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Authorized")
}
