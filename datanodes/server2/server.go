package main

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
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
	log.Println("Data Node 2 is listening on:", s.listenAddr)
	s.acceptConnection()

	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		requestType, err := reader.ReadString('\n')
	
		if err != nil {
			if err == io.EOF {
				log.Println("Client Disconnected")
				return
			}
			log.Fatalf("Read error : %v", err)
			return
		}

		requestType = strings.TrimSpace(requestType)
		log.Println("Request type:", requestType)

		switch requestType {
		case "GET":
			if err := handleGetRequest(conn, reader, writer); err != nil {
				log.Fatalf("GET Failed: %v", err)
			}
		case "POST":
			if err := handlePostRequest(reader, writer); err != nil {
				log.Fatalf("POST Failed: %v", err)
			}
		default:
			writer.WriteString("Error: Invalid Request\n")
			writer.Flush()
			return
		}

	}
}

func handleGetRequest(conn net.Conn, reader *bufio.Reader, writer *bufio.Writer) error {

	filename, err := reader.ReadString('\n')
	
	if err != nil {
		return fmt.Errorf("Failed to read the filename %v", err)
	}

	log.Println("filename is", filename)

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

	log.Println("ChunkId is ", chunkId)

	data, err := getBytes(filename, chunkId)

	if err != nil {
		return err
	}

	// I am going to do a small change here like using the length prefixed protocol 
	// firstly we are sending the size of the chunk so that the client must read all the data as per the size.

	chunkSize := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, chunkSize); err != nil {
		return fmt.Errorf("Failed to send the chunk size: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("Error sending data %v", err)
	}

	log.Printf("Sent Chunk %s : %d bytes", chunkId, len(data))
	return nil
}

func getBytes(filename string, chunkId string) ([]byte, error) {
	
	data, err := os.ReadFile(fmt.Sprintf("storage/%s/%s", filename, chunkId))
	if err != nil {
		return nil, fmt.Errorf("Failed to read the file: %v", err)
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

	log.Printf("Saving chunk %s, Size %d bytes", recievedData.FileId, len(recievedData.Data))

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

func main() {

	const listenAddr = ":9002"

	server := NewServer(listenAddr)
	if err := server.Start(); err != nil {
		fmt.Println("Server failed to start", err)
		os.Exit(1)
	}
}
