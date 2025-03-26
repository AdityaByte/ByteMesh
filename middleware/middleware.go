package middleware

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/AdityaByte/bytemesh/chunk"
)

var sizeOfChunk float64 = 30 // In kb

func decideParts(filesize float64) float64 { 
	return math.Ceil(filesize / sizeOfChunk)
}

func CreateChunk(file *os.File) (*[]chunk.Chunk, error) {

	fileData, err := os.ReadFile(file.Name())

	if err != nil {
		return nil, err
	}


	fmt.Println("original file size:", len(fileData))
	fmt.Println(len(fileData) )
	originalFileSize := float64(len(fileData))
	fileSizeInKb := originalFileSize / 1024

	fmt.Println("File size:", fileSizeInKb)

	parts := decideParts(float64(fileSizeInKb))

	fmt.Println("parts:", parts)

	var chunks []chunk.Chunk

	var newChunk []byte

	first := 0
	last := int(sizeOfChunk) * 1024

	for i := 0; i < int(parts); i++ {

		if last > len(fileData) {
			last = len(fileData)
		}

		newChunk = fileData[first:last]
		first = last
		last += int(sizeOfChunk) * 1024

		id := "Chunk" + strconv.Itoa(i+1)

		chunks = append(chunks, chunk.Chunk{
			Id: id,
			Data: newChunk,
		})
	}

	return &chunks,nil
}
