package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type metaData struct {
	Filename string
	FileExtension string
	Location map[string]string
}

type Server struct {
	listenAddr string
	ln         net.Listener
}

func NewServer(la string) *Server {
	return &Server{
		listenAddr: la,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("Error while listening to %s: %v", s.listenAddr, err)
	}

	s.ln = ln

	defer s.ln.Close()

	fmt.Println("Name node server is listening on", s.listenAddr)

	s.acceptConnection()

	return nil
} 

func (s *Server) acceptConnection() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Error connecting to client", err)
			return
		}

		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	defer conn.Close()
	metaData := &metaData{
		Location: make(map[string]string),
	}

	decoder := gob.NewDecoder(conn)
	if err := decoder.Decode(metaData); err != nil {
		fmt.Println("Error while decoding the metadata", err)
		return
	}

	fmt.Println("Filename is:", metaData.Filename)
	fmt.Println("FileExtension is:", metaData.FileExtension)
	fmt.Println("meta data location:", metaData.Location)

	if metaData.Location != nil {
		for key,value := range metaData.Location {
			fmt.Printf("%s -> %s\n", key, value)
		}
	} else {
		fmt.Println("key value data doesn't exists.")
	}
}

func main() {
	const addr = ":9004"
	server := NewServer(addr)
	if err := server.Start(); err != nil {
		fmt.Println("Server failed to start", err)
		os.Exit(1)
	}
}