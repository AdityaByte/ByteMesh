package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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
	fmt.Println("Data node server 1 is listening on ", s.listenAddr)
	s.acceptConnection()

	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
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

	// Firstly we have to check the request type ok
	reader := bufio.NewReader(conn)
	requestType, err := reader.ReadString('\n')
	
	if err != nil {
		fmt.Println("Failed to fetch the request type: ", err)
		return
	}

	requestType = strings.TrimSpace(requestType)
	fmt.Println("Request type is ", requestType)

	switch(requestType) {
	case "GET":
		if err := handleGetRequest(conn, reader); err != nil {
			log.Fatalf("%v", err)
		}
	case "POST":
		handlePostRequest(reader)
	default:
		fmt.Println("Request type not found:", requestType)
		return
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

func handleGetRequest(conn net.Conn, reader *bufio.Reader) error {

	
	filename, err := reader.ReadString('\n')
	
	if err != nil {
		return fmt.Errorf("Failed to read the filename %v", err)
	}

	fmt.Println("filename is", filename)

	isEmpty := utils.CheckEmptyField(filename)
	
	if isEmpty {
		return fmt.Errorf("Filename is empty")
	}

	filename = strings.TrimSpace(filename)

	chunkId, err := reader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("Failed to read chunkid: %v", err)
	}

	isEmpty = utils.CheckEmptyField(chunkId)

	if isEmpty {
		return fmt.Errorf("Chunk Id is empty")
	}

	chunkId = strings.TrimSpace(chunkId)

	fmt.Println("ChunkId is ", chunkId)

	data, err := getBytes(filename, chunkId)

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("Error sending data %v", err)
	}
	writer.Flush()

	fmt.Println("Data sent successfully")
	return nil
}

func getBytes(filename string, chunkId string) ([]byte, error) {
	
	data, err := os.ReadFile(fmt.Sprintf("storage/%s/%s", filename, chunkId))
	if err != nil {
		return nil, fmt.Errorf("Failed to read the file: %v", err)
	}

	return data, nil
}

func handlePostRequest(reader *bufio.Reader) error {
	decoder := gob.NewDecoder(reader)
	var recievedData chunkData
	if err := decoder.Decode(&recievedData); err != nil {
		return fmt.Errorf("Failed to decode the data %v", err)
	}
	fmt.Println("FileId", recievedData.FileId)
	fmt.Println("Filename", recievedData.Filename)
	fmt.Println("Data length:", len(recievedData.Data))

	err := os.MkdirAll(fmt.Sprintf("storage/%s", strings.Trim(recievedData.Filename, "\n")), os.ModePerm)

	if err != nil {
		return fmt.Errorf("Failed to create directory %v", err) 
	}

	err = os.WriteFile(fmt.Sprintf("storage/%s/%s", strings.Trim(recievedData.Filename, "\n"), recievedData.FileId), recievedData.Data, 0644)

	if err != nil {
		return fmt.Errorf("Failed to save the file %v", err)
	}

	fmt.Println("Chunk data saved successfully")

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
		fmt.Println("Server failed to start", err)
		os.Exit(1)
	}
}
