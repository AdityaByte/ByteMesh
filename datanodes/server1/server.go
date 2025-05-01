package main

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
	"github.com/AdityaByte/bytemesh/utils"
)

type chunkData struct {
	Filename string
	FileId   string
	Data     []byte
}

type Server struct {
	listenAddr string
	listener   net.Listener
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.listener = listener
	defer s.listener.Close()
	logger.InfoLogger.Println("Data Node 1 is listening on:", s.listenAddr)
	s.acceptConnection()

	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.ErrorLogger.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// In the handleConnection we firslty we have to define a proper schema that in which form is the request
// is been sent -> Request is been of two types
// 1.Get Request -> for getting out some resources from the server
// 2.Post Request -> This one is for pushing out some resources to the server.

/*
1. For Post request
Request Type - POST
Headers  -> file name and other things
Body - actual data

2. For Get request
Request Type - GET
headers -> filename and chunkid
*/

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		requestType, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				logger.InfoLogger.Println("Client Disconnected")
				return
			}
			logger.ErrorLogger.Println("Failed to fetch the request type: %v", err)
			return
		}

		requestType = strings.TrimSpace(requestType)
		logger.InfoLogger.Println("Request type:", requestType)

		switch requestType {
		case "GET":
			if err := handleGetRequest(reader, writer); err != nil {
				logger.ErrorLogger.Println("GET Failed: %v", err)
			}
		case "POST":
			if err := handlePostRequest(reader, writer); err != nil {
				logger.ErrorLogger.Println("POST Failed: %v", err)
			}
		case "HEALTH":
			logger.InfoLogger.Println("Handling health check..")
			if err := Health(conn, reader, writer); err != nil {
				logger.ErrorLogger.Println("ERROR: %v\n", err)
			}
		default:
			writer.WriteString("Error: Invalid Request\n")
			writer.Flush()
			return
		}

	}

	// decoder := gob.NewDecoder(conn)

	// var recievedData chunkData
	// err = decoder.Decode(&recievedData)

	// if err != nil {
	// 	fmt.Println("Error while decoding data", err)
	// 	return
	// }

	// fmt.Println("File id", recievedData.FileId)
	// fmt.Println("File name", recievedData.Filename)
	// fmt.Println("File data len", len(recievedData.Data))

	// // filename -> foldername
	// // fileid -> filename

	// // Major change : When we deserliaze the data it automatically reads it and converts to its original form.
	// // reader := bufio.NewReader(conn)
	// // recievedData.Data = make([]byte, 30*1024)
	// // n, err := conn.Read(recievedData.Data)
	// // n, err := reader.Read(recievedData.Data)

	// // if err != nil {
	// // 	if err == io.EOF {
	// // 		fmt.Println("Connection is closed by sender, but some data may have been recieved.", err)
	// // 	} else {
	// // 		fmt.Println("Error while reading data", err)
	// // 	}
	// // } else {
	// // 	fmt.Printf("Recieved %d bytes of data\n", n)
	// // }

	// err = os.MkdirAll(fmt.Sprintf("storage/%s", strings.Trim(recievedData.Filename, "\n")), os.ModePerm)
	// // err = os.MkdirAll("storage/" + strings.Trim(chunkData.filename, "\n"), os.ModePerm)

	// if err != nil {
	// 	fmt.Println("Error creating directory", err)
	// 	return
	// }

	// err = os.WriteFile(fmt.Sprintf("storage/%s/%s", strings.Trim(recievedData.Filename, "\n"), recievedData.FileId), recievedData.Data, 0644)

	// if err != nil {
	// 	fmt.Println("Error saving file", err)
	// }

	// fmt.Println("Chunk data saved successfully")
}

func handleGetRequest(reader *bufio.Reader, writer *bufio.Writer) error {

	filename, err := reader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("ERROR: Failed to read the filename %v", err)
	}

	logger.InfoLogger.Println("Filename:", filename)

	isEmpty := utils.CheckEmptyField(filename)

	if isEmpty {
		return fmt.Errorf("ERROR: Filename is empty")
	}

	filename = strings.TrimSpace(filename)

	chunkId, err := reader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("ERROR: Failed to read chunkid: %v", err)
	}

	isEmpty = utils.CheckEmptyField(chunkId)

	if isEmpty {
		return fmt.Errorf("ERROR: Chunk Id is empty")
	}

	chunkId = strings.TrimSpace(chunkId)

	logger.InfoLogger.Println("ChunkId is ", chunkId)

	data, err := getBytes(filename, chunkId)

	if err != nil {
		return err
	}

	// I am going to do a small change here like using the length prefixed protocol
	// firstly we are sending the size of the chunk so that the client must read all the data as per the size.

	chunkSize := uint32(len(data))
	if err := binary.Write(writer, binary.BigEndian, chunkSize); err != nil {
		return fmt.Errorf("ERROR: Failed to send the chunk size: %v", err)
	}

	nn, err := writer.Write(data)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to send data %v", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ERROR: Failed to Flush out the data: %v", err)
	}

	logger.InfoLogger.Println("Length of the data:", nn)
	logger.InfoLogger.Println("Data sent successfully")
	return nil
}

func getBytes(filename string, chunkId string) ([]byte, error) {

	data, err := os.ReadFile(fmt.Sprintf("storage/%s/%s", filename, chunkId))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to read the file: %v", err)
	}

	return data, nil
}

func handlePostRequest(reader *bufio.Reader, writer *bufio.Writer) error {
	decoder := gob.NewDecoder(reader)
	var recievedData chunkData
	if err := decoder.Decode(&recievedData); err != nil {
		writer.WriteString("ERROR: Failed to Decode Data\n")
		writer.Flush()
		return err
	}

	logger.InfoLogger.Printf("Saving chunk %s, Size %d bytes", recievedData.FileId, len(recievedData.Data))

	tempPath := fmt.Sprintf("storage/%s/%s.tmp", recievedData.Filename, recievedData.FileId)
	finalPath := fmt.Sprintf("storage/%s/%s", recievedData.Filename, recievedData.FileId)

	if err := os.MkdirAll(fmt.Sprintf("storage/%s", strings.Trim(recievedData.Filename, "\n")), os.ModePerm); err != nil {
		writer.WriteString("ERROR: Failed to create directory\n")
		writer.Flush()
		return err
	}

	if err := os.WriteFile(tempPath, recievedData.Data, 0644); err != nil {
		writer.WriteString("ERROR: Failed to write the chunk\n")
		writer.Flush()
		return err
	}

	// If everything goes right at last we have to rename it to the finalPath name
	if err := os.Rename(tempPath, finalPath); err != nil {
		writer.WriteString("Failed to commit the chunk\n")
		writer.Flush()
		return err
	}

	writer.WriteString("OK\n")
	writer.Flush()
	return nil
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	reader := bufio.NewReader(conn)

// 	chunkName, err := reader.ReadString('\n')
// 	if err != nil {
// 		if err.Error() == "EOF" {
// 			fmt.Println("No chunks recieved, closing connection")
// 			return
// 		}
// 		fmt.Println("Error reading chunk name:", err)
// 		return
// 	}

// 	chunkName = strings.TrimSpace(chunkName)
// 	if chunkName == "" {
// 		fmt.Println("Recieved an empty chunk name closing connection.")
// 		return
// 	}

// 	chunkData := make([]byte, 30*1024)
// 	n, err := reader.Read(chunkData)
// 	if err != nil {
// 		fmt.Println("Error reading chunk data:", err)
// 		return
// 	}

// 	// Saving the chunk to the storage
// 	err = os.MkdirAll("storage/", os.ModePerm)
// 	if err != nil {
// 		fmt.Println("Error creating directory", err)
// 		return
// 	}

// 	err = os.WriteFile("storage/"+chunkName, chunkData[:n], 0644)
// 	if err != nil {
// 		fmt.Println("Error saving file", err)
// 		return
// 	}

// 	fmt.Println("Chunk saved successfully..")

// }

func main() {

	const listenAddr = ":9001"

	server := NewServer(listenAddr)
	if err := server.Start(); err != nil {
		logger.ErrorLogger.Fatalf("Server failed to start %v", err)
		os.Exit(1)
	}
}
