package main

import (
	"bufio"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"namenodeserver/database"
	"namenodeserver/health"
	"namenodeserver/logger"
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

	logger.InfoLogger.Println("Name Node is listening on", s.listenAddr)

	s.acceptConnection()

	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			logger.ErrorLogger.Println("Error connecting to client", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	// Firstly of all we have to sent the OK response to the coordinator
	// for telling that everything goes right to the server.

	defer func() {
		if err := conn.Close(); err != nil {
			logger.ErrorLogger.Println("Failed to close connection %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mongoRepo, err := database.LoadMongoRepository()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	reader := bufio.NewReader(conn)
	requestType, err := reader.ReadString('\n')

	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	requestType = strings.TrimSpace(requestType)

	if requestType == "" {
		logger.ErrorLogger.Println("Request Type not found.")
		return
	}

	fmt.Println("Request type:", requestType)

	switch requestType {
	case "GET":
		if err := handleGetRequest(ctx, conn, reader, mongoRepo); err != nil {
			logger.ErrorLogger.Println(err)
		}
	case "POST":
		if err := handlePostRequest(ctx, conn, reader, mongoRepo); err != nil {
			logger.ErrorLogger.Println(err)
		}
	case "HEALTH":
		if err := health.Health(conn, reader); err != nil {
			logger.ErrorLogger.Println(err)
		}
	case "GETALL":
		if err := handleGetAllRequest(ctx, conn, reader, mongoRepo); err != nil {
			logger.ErrorLogger.Println(err)
		}

	default:
		logger.ErrorLogger.Printf("Request Type : {%s} not found\n", requestType)
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
		return fmt.Errorf("Error occured while reading the filename: %v", err)
	}
	filename = strings.TrimSpace(filename)

	if filename == "" {
		return fmt.Errorf("ERROR: Empty Field - Filename")
	}

	nameAndExtension := strings.Split(filename, ".")

	name := nameAndExtension[0]
	extension := nameAndExtension[len(nameAndExtension)-1]

	metadata, err := service.FetchMetaData(ctx, name, extension, mongoRepo)

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
	if err = encoder.Encode(metadata); err != nil {
		return fmt.Errorf("Failed to encode the data: %v", err)
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

func handleGetAllRequest(ctx context.Context, conn net.Conn, reader *bufio.Reader, repo *database.MongoRepository) error {
	defer conn.Close()
	user, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ERROR: Failed to read the username: %v", err)
	}
	// Now what we have to do we have to send that request to the service handler ok.

	user = strings.TrimSpace(user)

	data, err := service.FetchUserSpecificMetaData(ctx, user, repo)
	if err != nil {
		return err
	}

	// Now we need to encode the data with the help of gob encoder.
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(&data); err != nil {
		return fmt.Errorf("ERROR: Failed to encode the data: %v", err)
	}

	logger.InfoLogger.Println("Data sent successfully.")
	return nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found.")
	}
}

func main() {
	const addr = ":9004"
	server := NewServer(addr)
	if err := server.Start(); err != nil {
		logger.ErrorLogger.Fatalf("Server failed to start %v", err)
		os.Exit(1)
	}
}
