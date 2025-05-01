package client

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
)

// The Download() function takes the name of the file with extension
// as parameter and sends the request to the namenode server so if the data exists
// at the server and the filename matches with any of the file which is at the server then we
// will get the request with request code 200 ok.

const addr = ":9004"

func Download(filename string) error {

	logger.InfoLogger.Println("Filename:", filename)

	data, err := middleware.GetChunks(filename)
	if err != nil {
		return err
	}

	logger.InfoLogger.Println("Downloaded Data length:", len(*data))

	if _, err = os.Stat("../download/"); os.IsNotExist(err) {
		if err := os.Mkdir("../download/", os.ModePerm); err != nil {
			logger.ErrorLogger.Println("Failed to create the download directory:", err)
		} else {
			logger.InfoLogger.Println("Download Directory created successfully.")
		}
	} else {
		logger.InfoLogger.Println("Download Directory already exists.")
	}

	file, err := os.Create("../download/" + filename)

	if err != nil {
		return fmt.Errorf("Error creating file %v", err)
	}

	reader := bytes.NewReader(*data)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("Failed to create file %v", err)
	}

	return nil
}
