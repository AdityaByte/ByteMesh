package main

import (
	"fmt"
	"os"

	"github.com/AdityaByte/bytemesh/client"
	"github.com/AdityaByte/bytemesh/middleware"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage go run . <filelocation>")
		return
	}

	fileLocation := os.Args[1]

	file, err := client.Upload(fileLocation);

	if err != nil {
		fmt.Println(err)
		return
	}

	chunks, err := middleware.CreateChunk(file)

	if err != nil {
		fmt.Println(err)
	}


	for _, chunk := range *chunks {
		fmt.Println("chunk id:", chunk.Id)
		fmt.Println("chunk datasize:", len(chunk.Data) , "byte")
		fmt.Println("------------")
	}
}
