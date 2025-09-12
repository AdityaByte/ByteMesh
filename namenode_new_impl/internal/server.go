// Package internal provides utilities for managing metadata.
package internal

import (
	"bufio"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/AdityaByte/namenode/internal/database"
	"github.com/AdityaByte/namenode/internal/handler"
	"github.com/AdityaByte/namenode/internal/model"
	"github.com/AdityaByte/namenode/internal/payloads"
	"github.com/AdityaByte/namenode/logger"
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

// handleConnection handle coordinator incoming requests and returns a response.
func handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)

	// Properly closing the resource at the end.
	// Not closing the connection cause of implementing the connection pool.
	defer func() {
		if err := conn.Close(); err != nil {
			logger.ErrorLogger.Println("Failed to close connection %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	mongoRepo, err := database.LoadMongoRepository()
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

		final_metadata, err := handler.HandleGetRequest(ctx, *mongoRepo, DataNodes, fullFileName)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return
		}

		// Encoding and sending the metadata as a go binary object.
		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(&final_metadata); err != nil {
			logger.ErrorLogger.Printf("Failed to encode the metadata payload, %v\n", err)
			return
		}

		logger.InfoLogger.Println("Metadata sents successfully.")

	case "POST":

		var metadata model.MetaData

		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(&metadata); err != nil {
			logger.ErrorLogger.Printf("Failed to decode the data, %v\n", err)
		}

		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*15)
		defer cancel()

		if err := handler.HandlePostRequest(ctx, mongoRepo, metadata); err != nil {
			logger.ErrorLogger.Println(err.Error())
		}

		logger.InfoLogger.Println("Metadata inserted successfully to the database")

	case "HEALTH":
		aliveNodes := handler.GetAliveNodes(*DataNodes)

		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(&aliveNodes); err != nil {
			logger.ErrorLogger.Printf("Failed to encode and send the Helath query of datanodes, %v", err)
			return
		}

		logger.InfoLogger.Println("Health information of datanodes sent successfully.")

	default:
		logger.ErrorLogger.Printf("Request Verb : {%s} not found\n", requestVerb)
	}
}
