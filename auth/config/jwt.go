package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
	"github.com/golang-jwt/jwt/v5"
)

// Here we have to create a jwt token each time the users logs in and we have to save that in a seperate folder.

var (
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

func CreateToken(username string) (string, error) {
	// Here we have to generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", fmt.Errorf("ERROR: Failed to generate the private key, %s", err)
	}

	// here we need to set the private key.
	publicKey = &privateKey.PublicKey

	// Now we have to generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", err
	}

	logger.InfoLogger.Println("Real Token:", tokenString)

	return tokenString, nil
}

// Function for verifying the token.
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("ERROR: Unexpected signing method %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}

	if !token.Valid {
		return fmt.Errorf("ERROR: Invalid token")
	}

	logger.InfoLogger.Println("Valid Token")
	return nil
}
