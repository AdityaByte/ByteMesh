package utils

import (
	"fmt"
	"strings"
)

func Getfilename(filepath string) string {
	arr := strings.Split(filepath, "\\")
	filename := arr[len(arr)-1]
	fmt.Println("filename is ", filename)
	return filename
}
