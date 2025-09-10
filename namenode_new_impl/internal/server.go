package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/AdityaByte/namenode/internal/handler"
	"github.com/AdityaByte/namenode/internal/payloads"
)

var DataNodes *payloads.RegisteredDataNodes

type Server struct {
	listenAddr string
	listener   net.Listener
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (server *Server) Start() error {
	listener, err := net.Listen("tcp", server.listenAddr)
	if err != nil {
		return fmt.Errorf("Error while listening to %s: %v", server.listenAddr, err)
	}

	server.listener = listener
	defer server.listener.Close()

	logger.InfoLogger.Printf("Server is listening on %s\n", server.listenAddr)

	server.acceptConnection()

	return nil
}

func (server *Server) acceptConnection() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			logger.ErrorLogger.Println("Error connecting to client", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)

	// Properly closing the resource at the end.
	// Not closing the connection cause of implementing the connection pool.
	// defer func() {
	// 	if err := conn.Close(); err != nil {
	// 		logger.ErrorLogger.Println("Failed to close connection %v\n", err)
	// 	}
	// }()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	mongoRepo, err := LoadMongoRepository()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	requestVerb, err := reader.ReadString('\n')
	if err != nil {
		logger.ErrorLogger.Printf("Failed to read the string, %v\n", err)
		return
	}

	requestVerb = strings.TrimSpace(requestVerb)
	if requestVerb == "" {
		logger.ErrorLogger.Println("Empty request were found")
	}

	logger.InfoLogger.Println("Request Verb:", requestVerb)

	switch requestVerb {
	case "REGISTER":
		// Calling the Registration Handler method.
		decoder := json.NewDecoder(conn)
		var reg_node *payloads.DataNode
		if err := decoder.Decode(&reg_node); err != nil {
			logger.ErrorLogger.Printf("Invalid JSON, %v\n", err)
			return
		}
		handler.NodeRegistrationHandler(conn, reg_node, DataNodes)

	case "HEARTBEAT":
		// Now we need to decode the heartbeat data.
		decoder := json.NewDecoder(conn)
		var heartbeat *payloads.HeartBeat
		if err := decoder.Decode(&heartbeat); err != nil {
			logger.ErrorLogger.Printf("Invalid JSON, %v\n", err)
			return
		}
		handler.HeartBeatHandler(DataNodes, heartbeat)
	case "GET":
		// Request is being sent by the coordinator.
		// Get request usually has a filename.
		fullFileName, err := reader.ReadString('\n')
		if err != nil {
			logger.ErrorLogger.Printf("Failed to read the filename, %v", err)
			return
		}
		fullFileName = strings.TrimSpace(fullFileName)
		if fullFileName == "" {
			logger.ErrorLogger.Println("Empty filename recieved")
			return
		}
		data, err := handler.HandleGetRequest(ctx, mongoRepo, DataNodes, fullFileName)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return
		}

		// At the end when we gets the data we are just forwarding it to the client.
		// Before that I need to send the length of the data too.
		lengthofData := fmt.Sprintf("Data_Length=%d\n", len(data))

		// Now we are sending the data length to string format.
		if _, err := conn.Write([]byte(lengthofData)); err != nil {
			logger.ErrorLogger.Printf("Failed to send the data to coordinator, %v\n", err)
			return
		}

		// Now I have to send the actual data.
		if _, err := conn.Write(data); err != nil {
			logger.ErrorLogger.Printf("Failed to send the data to coordinator, %v\n", err)
			return
		}

		// That's all I have to do in that Get case.

	case "POST":

	default:
		logger.ErrorLogger.Printf("Request Verb : {%s} not found\n", requestVerb)
	}
}
