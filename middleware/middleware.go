package middleware

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/AdityaByte/bytemesh/coordinator"
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

	fmt.Println("original file size:", len(fileData))
	fmt.Println(len(fileData))
	originalFileSize := float64(len(fileData))
	fileSizeInKb := originalFileSize / 1024

	fmt.Println("File size:", fileSizeInKb)

	parts := decideParts(float64(fileSizeInKb))

	fmt.Println("parts:", parts)

	var chunks []models.Chunk

	var newChunk []byte

	first := 0
	last := int(sizeOfChunk) * 1024

	for i := 0; i < int(parts); i++ {

		log.Println("Iteration", i, "First:", first, "Last:", last, "FileData size:", len(fileData))

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

	fmt.Println("File name is :", file.Name())

	if strings.Contains(file.Name(), "../storage") {
		newString := strings.TrimPrefix(file.Name(), "../storage/")
		log.Println("New String is:", newString)
		return &chunks, newString, fileSizeInKb, nil
	}

	newString := strings.TrimPrefix(file.Name(), "storage/")
	log.Println("New String is:", newString)
	return &chunks, newString, fileSizeInKb, nil
}

func GetChunks(filename string) (*[]byte, error) {
	conn, err := net.Dial("tcp", nameNode)
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error occured while closing the connection..", err)
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to namenode", err)
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("File name is empty")
	}

	writer := bufio.NewWriter(conn)
	writer.WriteString("GET\n" + filename + "\n")
	writer.Flush()

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	response = strings.TrimSpace(response)
	fmt.Println("The response we are getting is ", response)

	if response != "200" {
		return nil, fmt.Errorf("Response is not OK", response)
	}

	decoder := gob.NewDecoder(conn)
	var recievedData models.MetaData
	err = decoder.Decode(&recievedData)

	if err != nil {
		return nil, fmt.Errorf("Error occured while decoding the data", err)
	}

	fmt.Println("metadata is :", recievedData)

	if reflect.DeepEqual(recievedData, models.MetaData{}) {
		return nil, fmt.Errorf("No file found at the server...")
	}

	recievedFileData, err := coordinator.FetchChunks(&recievedData)

	if err != nil {
		return nil, err
	}

	return recievedFileData, nil
}
