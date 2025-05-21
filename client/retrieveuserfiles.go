package client

import (
	"fmt"
	"os"

	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
)

func RetriveUserFiles() error {
	// So basically it sends the request to the middleware then the middleware makes the connection to the namenode server then we forward that connection to the coordinator it fetches all that content and we recieve it so that's the flow.

	// Here we need to read the file .cred and find the user which has been logged in.
	user, err := os.ReadFile("../.auth/.cred")
	if err != nil {
		return fmt.Errorf("Failed to read the file: %v", err)
	}

	// Now here we have to reterive the user files ok.
	data, err := middleware.FetchUserFiles(string(user))
	if err != nil {
		return err
	}

	fmt.Println("data:", *data)

	// Else here we need to print the data ok.
	// Data that we recieved is a pointer to the models.metadata ok
	for singleData := range *data {
		// So singleData is the single struct - for ex
		logger.InfoLogger.Println(singleData) // ok
	}

	return nil
}