package middleware

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/AdityaByte/bytemesh/coordinator"
	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/models"
)

const (
	sizeOfChunk = 500 // In KB
	nameNode    = ":9004"
)

func decideParts(filesize float64) float64 {
	return math.Ceil(filesize / sizeOfChunk)
}

func CreateChunk(file *os.File) (*[]models.Chunk, string, float64, error) {

	fileData, err := os.ReadFile(file.Name())

	if err != nil {
		return nil, "", 0, err
	}

	logger.InfoLogger.Println("original file size:", len(fileData))
	originalFileSize := float64(len(fileData))
	fileSizeInKb := originalFileSize / 1024

	logger.InfoLogger.Println("File size in kb:", fileSizeInKb)

	parts := decideParts(float64(fileSizeInKb))

	logger.InfoLogger.Println("parts:", parts)

	var chunks []models.Chunk

	var newChunk []byte

	first := 0
	last := int(sizeOfChunk) * 1024

	for i := 0; i < int(parts); i++ {

		logger.InfoLogger.Println("Iteration", i, "First:", first, "Last:", last, "FileData size:", len(fileData))

		if last > len(fileData) {
			last = len(fileData)
		}

		newChunk = fileData[first:last]
		first = last
		last += int(sizeOfChunk) * 1024

		id := "Chunk" + strconv.Itoa(i+1)

		chunks = append(chunks, models.Chunk{
			Id:   id,
			Data: newChunk,
		})
	}

	logger.InfoLogger.Println("File name is :", file.Name())

	if strings.Contains(file.Name(), "../storage") {
		newString := strings.TrimPrefix(file.Name(), "../storage/")
		logger.InfoLogger.Println("New String is:", newString)
		return &chunks, newString, fileSizeInKb, nil
	}

	newString := strings.TrimPrefix(file.Name(), "storage/")
	logger.InfoLogger.Println("New String is:", newString)
	return &chunks, newString, fileSizeInKb, nil
}

func GetChunks(filename string) (*[]byte, error) {
	conn, err := net.Dial("tcp", nameNode)
	defer func() {
		if err := conn.Close(); err != nil {
			logger.ErrorLogger.Println("Failed to close connection", err)
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to namenode", err)
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("Empty filename")
	}

	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString("GET\n" + filename + "\n"); err != nil {
		return nil, fmt.Errorf("Failed to send the request:", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("Flush Failed")
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	response = strings.TrimSpace(response)
	logger.InfoLogger.Println("The response we are getting is ", response)

	if response != "200" {
		return nil, fmt.Errorf("Response is not OK", response)
	}

	decoder := gob.NewDecoder(conn)
	var recievedData models.MetaData
	err = decoder.Decode(&recievedData)

	if err != nil {
		return nil, fmt.Errorf("Error occured while decoding the data", err)
	}

	logger.InfoLogger.Println("metadata is :", recievedData)

	if reflect.DeepEqual(recievedData, models.MetaData{}) {
		return nil, fmt.Errorf("No file found at the server...")
	}

	recievedFileData, err := coordinator.FetchChunks(&recievedData)

	if err != nil {
		return nil, err
	}

	return recievedFileData, nil
}
