package client

import (
	"fmt"
	"io"
	"os"

	"github.com/AdityaByte/bytemesh/utils"
)

func Upload(filelocation string) (*os.File, error) {
	srcFile, err := os.Open(filelocation)

	if err != nil {
		return nil, fmt.Errorf("Error in opening the file")
	}

	defer srcFile.Close()

	destPath := "storage/" + utils.Getfilename(srcFile.Name())
	destFile, err := os.Create(destPath)

	if err != nil {
		return nil, fmt.Errorf("Error in creating the file at the destination path")
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	if err != nil {
		return nil, fmt.Errorf("Error while copying file")
	}

	return destFile, nil
}
