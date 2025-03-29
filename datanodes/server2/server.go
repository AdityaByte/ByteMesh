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

	err = os.MkdirAll(fmt.Sprintf("storage/%s", strings.Trim(recievedData.Filename, "\n")), os.ModePerm)

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

func main() {

	const listenAddr = ":9002"

	server := NewServer(listenAddr)
	if err := server.Start(); err != nil {
		fmt.Println("Server failed to start", err)
		os.Exit(1)
	}
}
