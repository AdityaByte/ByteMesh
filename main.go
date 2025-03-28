package main

import (
	"fmt"
	"os"

	"github.com/AdityaByte/bytemesh/client"
	"github.com/AdityaByte/bytemesh/middleware"
	"github.com/AdityaByte/bytemesh/server"
)

// Program flow
// main.go - entry point client communicate with it 
// then the request is being sent to the client's upload function 
// which uploads the file in a local storage location
// After that the middleware creates the parts like how many parts are being created of it 
// and the chunk then everything goes correctly it returns a chunk array to me
// which i then passed to the server function which takes care of the distribution of the chunks
// to the datanodes right now we have 3 datanodes 
// if everything goes correctly then a chunk data nodes are stored at data nodes server correctly.

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage go run . <filelocation>")
		return
	}

	fileLocation := os.Args[1]

	file, err := client.Upload(fileLocation)

	if err != nil {
		fmt.Println(err)
		return
	}

	chunks, filename, err := middleware.CreateChunk(file)
	

	if err != nil {
		fmt.Println(err)
	}

	if err := server.SendChunks(chunks, filename); err != nil {
		fmt.Println("Error sending chunks to the server nodes", err)
		return
	}

	for _, chunk := range *chunks {
		fmt.Println("chunk id:", chunk.Id)
		fmt.Println("chunk datasize:", len(chunk.Data), "byte")
		fmt.Println("------------")
	}
}
