package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
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

	reader := bufio.NewReader(conn)
	requestType, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	requestType = strings.TrimSpace(requestType)

	if requestType == "" {
		fmt.Println("Request type is not specified..")
		return
	}

	fmt.Println("Requst type:", requestType)

	switch(requestType) {
	case "GET":
		fmt.Println("Hey i am here get")
		filename, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error occured while reading the filename", err)
		}
		filename = strings.TrimSpace(filename)

		if filename == "" {
			fmt.Println("Filename is empty")
			return
		}

		fmt.Println("Filename is ", filename)

		if filename == "dfs-flowchart.png" {
			fmt.Println("Now i am here in the filename block")
			newMetaData := &metaData{
				Filename: "dfs-flowchart",
				FileExtension: "png",
				Location: map[string]string{
					"Node1": "chunk1",
					"Node0": "chunk2",
				},
			}

			conn.Write([]byte("200\n"))

			encoder := gob.NewEncoder(conn)
			err = encoder.Encode(newMetaData)

			if  err != nil {
				fmt.Println("Error occured while encoding the data", err)
				return
			}
		}
	case "POST":
		fmt.Println("Post request")
		return
	default:
		fmt.Printf("The particular request %s is not found\n", requestType)
		return
	}

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