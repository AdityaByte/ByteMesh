package utils

import (
	"os"

	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
)

// In this utility function we have to delete the file which was reside in the
// storage folder if something went wrong.
func RemoveFile(filelocation string) {
	if err := os.Remove(filelocation); err != nil {
		logger.ErrorLogger.Printf("Failed to remove the file: %s, error: %v", filelocation, err)
	}
	logger.InfoLogger.Println("File removed successfully.")
}
