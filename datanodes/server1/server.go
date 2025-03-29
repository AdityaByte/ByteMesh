package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)

	var recievedData chunkData
	err := decoder.Decode(&recievedData)

	fmt.Println("File id", recievedData.FileId)
	fmt.Println("File name", recievedData.Filename)
	fmt.Println("File data len", len(recievedData.Data))

	if err != nil {
		fmt.Println("Error while decoding data", err)
		return
	}

	// filename -> foldername
	// fileid -> filename

	// Major change : When we deserliaze the data it automatically reads it and converts to its original form.
	// reader := bufio.NewReader(conn)
	// recievedData.Data = make([]byte, 30*1024)
	// n, err := conn.Read(recievedData.Data)
	// n, err := reader.Read(recievedData.Data)

	// if err != nil {
	// 	if err == io.EOF {
	// 		fmt.Println("Connection is closed by sender, but some data may have been recieved.", err)
	// 	} else {
	// 		fmt.Println("Error while reading data", err)
	// 	}
	// } else {
	// 	fmt.Printf("Recieved %d bytes of data\n", n)
	// }

	err = os.MkdirAll(fmt.Sprintf("storage/%s", strings.Trim(recievedData.Filename, "\n")), os.ModePerm)
	// err = os.MkdirAll("storage/" + strings.Trim(chunkData.filename, "\n"), os.ModePerm)

	if err != nil {
		fmt.Println("Error creating directory", err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("storage/%s/%s", strings.Trim(recievedData.Filename, "\n"), recievedData.FileId), recievedData.Data, 0644)

	if err != nil {
		fmt.Println("Error saving file", err)
	}

	fmt.Println("Chunk data saved successfully")
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
