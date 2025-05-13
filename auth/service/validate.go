package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaByte/bytemesh/auth/config"
)

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "ERROR: Method not allowed")
		return
	}

	tokenString := r.Header.Get("Authorization")

	if strings.TrimSpace(tokenString) == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Missing Authorization header")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	if err := config.VerifyToken(tokenString); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid Token %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Authorized")
}
