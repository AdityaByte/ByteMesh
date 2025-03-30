package main

import (
	"bufio"
	"context"
	"encoding/gob"
	"fmt"
	"namenodeserver/database"
	"namenodeserver/model"
	"namenodeserver/service"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

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

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error closing the connection %v\n", err)
		}
	} ()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mongoRepo, err := database.LoadMongoRepository()
	if err != nil {
		fmt.Println(err)
		return
	}

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
		if err := handleGetRequest(ctx, conn, reader, mongoRepo); err != nil {
			fmt.Println(err)
		}
	case "POST":
		if err := handlePostRequest(ctx, conn, reader, mongoRepo); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("The particular request %s is not found\n", requestType)
	}

	// Old code for saving data
	// metaData := &model.MetaData{
	// 	Location: make(map[string]string),
	// }

	// decoder := gob.NewDecoder(conn)
	// if err := decoder.Decode(metaData); err != nil {
	// 	fmt.Println("Error while decoding the metadata", err)
	// 	return
	// }

	// fmt.Println("Filename is:", metaData.Filename)
	// fmt.Println("FileExtension is:", metaData.FileExtension)
	// fmt.Println("meta data location:", metaData.Location)

	// if metaData.Location != nil {
	// 	for key,value := range metaData.Location {
	// 		fmt.Printf("%s -> %s\n", key, value)
	// 	}
	// } else {
	// 	fmt.Println("key value data doesn't exists.")
	// }
}


func handleGetRequest(ctx context.Context, conn net.Conn, reader *bufio.Reader, mongoRepo *database.MongoRepository) error {
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error occured while reading the filename", err)
	}
	filename = strings.TrimSpace(filename)

	if filename == "" {
		return fmt.Errorf("Filename is empty")
	}

	fmt.Println("Filename is ", filename)

	metadata, err := service.FetchMetaData(ctx, filename, mongoRepo)

	if err != nil {
		return err
	}

	// Dummy data for testing.
	// if filename == "dfs-flowchart.png" {
	// 	fmt.Println("Now i am here in the filename block")
	// 	newMetaData := &model.MetaData{
	// 		Filename: "dfs-flowchart",
	// 		FileExtension: "png",
	// 		Location: map[string]string{
	// 			"Node1": "chunk1",
	// 			"Node0": "chunk2",
	// 	},
	// }

	conn.Write([]byte("200\n"))

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(metadata)

	if  err != nil {
		return err
	}

	return nil
}

func handlePostRequest(ctx context.Context, conn net.Conn, reader *bufio.Reader, mongoRepo *database.MongoRepository) error {
	defer conn.Close()
	metaData := &model.MetaData{
		Location: make(map[string]string),
	}
	
	decoder := gob.NewDecoder(reader)
	if err := decoder.Decode(metaData); err != nil {
		return fmt.Errorf("Error occured while decoding the data %v", err)
	}

	if err := service.PushMetaData(ctx, metaData, mongoRepo); err != nil {
		return err
	}

	_, err := conn.Write([]byte("200\n"))

	if err != nil {
		return err
	}

	return nil
}

func init() {
	if err:=godotenv.Load(); err != nil {
		fmt.Println("No .env file found.")
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