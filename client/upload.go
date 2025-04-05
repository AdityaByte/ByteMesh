package client

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AdityaByte/bytemesh/utils"
)

func Upload(filelocation string) (*os.File, error) {

	log.Println("File Location :", filelocation)

	srcFile, err := os.Open(filelocation)

	if err != nil {
		return nil, fmt.Errorf("Error in opening the file")
	}

	defer srcFile.Close()

	dirPath := "../storage/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			log.Println("Failed to created directory", err)
			return nil, err
		} else {
			log.Println("File created successfully")
		}
	} else {
		log.Println("Directory already exists.")
	}

	destPath := "../storage/" + utils.Getfilename(srcFile.Name()) // for debuggger
	// destPath := "storage/" + utils.Getfilename(srcFile.Name())
	destFile, err := os.Create(destPath)

	if err != nil {
		return nil, fmt.Errorf("Error in creating the file at the destination path %w", err)
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		return nil, fmt.Errorf("Error while copying file")
	}

	return destFile, nil
}
