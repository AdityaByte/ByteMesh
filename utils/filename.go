package utils

import (
	"strings"

	"github.com/AdityaByte/bytemesh/logger"
)

func Getfilename(filepath string) string {
	arr := strings.Split(filepath, "\\")
	filename := arr[len(arr)-1]
	logger.InfoLogger.Println("filename is ", filename)
	return filename
}
