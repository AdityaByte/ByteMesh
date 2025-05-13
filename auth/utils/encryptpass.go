package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ERROR: Failed to encrypt the password %v", err)
	}

	return string(encryptedPassword), nil
}