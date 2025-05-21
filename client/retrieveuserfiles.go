package client

import (
	"fmt"
	"os"

	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
	"github.com/AdityaByte/bytemesh/models"
	"github.com/AdityaByte/bytemesh/utils"
)

func RetriveUserFiles(username string) (*[]models.MetaData, error) {
	// So basically it sends the request to the middleware then the middleware makes the connection to the namenode server then we forward that connection to the coordinator it fetches all that content and we recieve it so that's the flow.

	// Here we need to read the file .cred and find the user which has been logged in.
	if utils.CheckEmptyField(username) {
		user, err := os.ReadFile("../.auth/.cred")
		if err != nil {
			return nil, fmt.Errorf("Failed to read the file: %v", err)
		}
		username = string(user)
	}

	// Now here we have to reterive the user files ok.
	response, err := middleware.FetchUserFiles(username)
	if err != nil {
		return nil, err
	}

	logger.InfoLogger.Println("data:", *response)

	// Else here we need to print the data ok.
	// Data that we recieved is a pointer to the models.metadata ok
	for _, data := range *response {
		// So singleData is the single struct - for ex
		logger.InfoLogger.Println(data) // ok
	}

	return response, nil
}
