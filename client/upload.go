package client

import (
	"fmt"
	"io"
	"os"

	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/utils"
)

func Upload(filelocation string) (*os.File, error) {

	logger.InfoLogger.Println("File Location :", filelocation)

	srcFile, err := os.Open(filelocation)

	if err != nil {
		return nil, fmt.Errorf("Failed to open file")
	}

	defer srcFile.Close()

	dirPath := "../storage/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			logger.ErrorLogger.Println("Failed to created directory", err)
			return nil, err
		} else {
			logger.InfoLogger.Println("File created successfully")
		}
	} else {
		logger.InfoLogger.Println("Directory already exists.")
	}

	destPath := "../storage/" + utils.Getfilename(srcFile.Name()) // for debuggger
	// destPath := "storage/" + utils.Getfilename(srcFile.Name())
	destFile, err := os.Create(destPath)

	if err != nil {
		return nil, fmt.Errorf("Failed to create file at the desired location %w", err)
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		return nil, fmt.Errorf("Failed to copy file")
	}

	return destFile, nil
}
