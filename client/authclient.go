package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/AdityaByte/bytemesh/datanodes/server2/logger"
	"github.com/AdityaByte/bytemesh/utils"
)

const authServerURL = "http://localhost:8080"

func SignUp(username string, password string) error {

	logger.InfoLogger.Printf("Username %s and Password %s\n", username, password)
	// 1st Check: Fields should not be empty.
	if utils.CheckEmptyField(username) || utils.CheckEmptyField(password) {
		return fmt.Errorf("ERROR: Field's should not be empty")
	}

	// 2nd Check: Password should be greater than six characters
	if len(password) <= 6 {
		return fmt.Errorf("ERROR: Password should be greater than 6 characters")
	}

	// Dude have to send the the post request to the auth server
	var jsonBuffer bytes.Buffer
	err := json.NewEncoder(&jsonBuffer).Encode(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return fmt.Errorf("ERROR: Failed to encode the data to json %v", err)
	}

	// Creating a new http request post.
	req, err := http.NewRequest("POST", authServerURL+"/signup", &jsonBuffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Now we have to send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Now we have to print the response.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.InfoLogger.Println("Status Code:", resp.Status)
	logger.InfoLogger.Println("Response Body:", string(body))

	return nil
}

func LogIn(username string, password string) error {

	if utils.CheckEmptyField(username) || utils.CheckEmptyField(password) {
		return fmt.Errorf("ERROR: Field's are empty")
	}

	var jsonBuffer bytes.Buffer
	json.NewEncoder(&jsonBuffer).Encode(map[string]string{
		"username": username,
		"password": password,
	})

	// Creating a new get request
	req, err := http.NewRequest("POST", authServerURL+"/login", &jsonBuffer)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to create a new request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Now we have to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to send the request to the auth server %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to read the response %v", err)
	}

	if resp.StatusCode == 200 {
		tokenString := resp.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		// Here we need to save the token locally.
		if err := os.MkdirAll("../.auth/", os.ModePerm); err != nil {
			return fmt.Errorf("ERROR: Failed to create the .auth directory %v", err)
		}

		// Now we need to save the auth token locally ok
		file, err := os.Create("../.auth/.jwt-token")
		if err != nil {
			return fmt.Errorf("ERROR: Failed to create the  token file %v", err)
		}

		// Now we need to write the content
		if _, err := file.WriteString(tokenString); err != nil {
			return fmt.Errorf("ERROR: Failed to write the tokenstring locally %v", err)
		}

		// We can also create a credentials and save the username too.
		file2, err := os.Create("../.auth/.cred")
		if err != nil {
			return fmt.Errorf("ERROR: Failed to create the .cred file %v", err)
		}

		// Now we need to write the content
		if _, err := file2.WriteString(username); err != nil {
			return fmt.Errorf("ERROR: Failed to write the cred string locally %v", err)
		}

		logger.InfoLogger.Println("Token saved successfully.")

	}

	logger.InfoLogger.Println("Response Code:", resp.Status)
	logger.InfoLogger.Println("Response Body:", string(body))

	return nil
}

func LogOut() error {
	if err := os.RemoveAll("../.auth"); err != nil {
		return fmt.Errorf("ERROR: Failed to log out %v", err)
	}

	return nil
}

func ValidateToken() error {
	data, err := os.ReadFile("../.auth/.jwt-token")
	if err != nil {
		return fmt.Errorf("Failed to read the token %v", err)
	}

	// Now we have to convert the token data to string
	tokenString := string(data)

	// logger.InfoLogger.Println("Fetch token from storage:", tokenString)

	// Now we have to create a request and send that to the authserver
	req, err := http.NewRequest("GET", authServerURL+"/validate", nil)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to create a new request %v", err)
	}

	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to do the request %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to read the response %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	logger.InfoLogger.Println("Response Status:", resp.Status)
	logger.InfoLogger.Println("Response Body:", string(body))

	return nil
}
