package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AdityaByte/bytemesh/client"
	"github.com/AdityaByte/bytemesh/coordinator"
	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
	"github.com/AdityaByte/bytemesh/utils"
)

// The main purpose of the distributed file storage system is to allow users to
// upload files to the cloud and retrieve them when needed.
//
// One of the key advantages of a distributed file storage system is its fault tolerance,
// which is achieved through its architecture.
//
// 1. **NameNode**: Responsible for handling metadata, such as file locations and storage management.
// 2. **DataNode**: We currently have three DataNodes where file chunks are stored.
// 3. **Client**: Can perform two main operations—uploading files or downloading them.
// 4. **Middleware**: Manages file uploads, downloads, and metadata retrieval.
// 5. **Coordinator**: Sends requests to servers, fetches data, and passes it to the middleware for further processing.

// Initializing function which initialize all the environment variables by
// loading out the .env file.

func main() {

	const version = "v1.0.0"
	const author = "@AdityaByte"
	fileLocation := flag.String("u", "", "Upload File Location")
	downloadFileName := flag.String("d", "", "Download File name")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage : %s\n", "bytemesh -u <Upload File Location>  -d <Download file name>")
		fmt.Fprintf(flag.CommandLine.Output(), "Version : %s\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Author : %s\n", author)
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *fileLocation != "" {
		file, err := client.Upload(*fileLocation)
		if err != nil {
			utils.RemoveFile(file.Name()) // If something fails out we remove the file from the local folder.
			logger.ErrorLogger.Fatalf("%v", err)
		}

		chunks, filename, filesize, err := middleware.CreateChunk(file)
		if err != nil {
			utils.RemoveFile(file.Name())
			logger.ErrorLogger.Fatalf("%v", err)
		}

		for i, chunk := range *chunks {
			logger.InfoLogger.Println("Iteration:", i, "Chunk ID:", chunk.Id, "Data Length:", len(chunk.Data))
		}

		if err := coordinator.SendChunks(chunks, filename, filesize); err != nil {
			utils.RemoveFile(file.Name())
			logger.ErrorLogger.Fatalf("%v", err)
		}

		if err := os.Remove(fmt.Sprintf("../storage/%s", filename)); err != nil {
			logger.ErrorLogger.Println("Failed to remove the file from the storage folder:", err)
		}
	}

	if *downloadFileName != "" {
		if err := client.Download(*downloadFileName); err != nil {
			logger.ErrorLogger.Fatalf("%v", err)
		}
	}

}
