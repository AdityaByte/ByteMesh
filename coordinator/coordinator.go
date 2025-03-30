package coordinator

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"strings"

	"github.com/AdityaByte/bytemesh/chunk"
)

type metaData struct {
	Filename      string
	FileExtension string
	Location      map[string]string
}

type chunkData struct {
	Filename string
	FileId   string
	Data     []byte
}

const (
	node1    = ":9001"
	node2    = ":9002"
	node3    = ":9003"
	nameNode = ":9004"
)

func SendChunks(chunks *[]chunk.Chunk, filename string) error {

	if filename == "" {
		return fmt.Errorf("File name does not exists.")
	}

	fmt.Println("original file name :", filename)
	fileData := strings.Split(filename, ".")
	name := fileData[0]
	extension := fileData[1]

	nodes := []string{node1, node2, node3}

	connections := make([]net.Conn, len(nodes))

	location := make(map[string]string)

	for i, node := range nodes {
		conn, err := net.Dial("tcp", node)
		if err != nil {
			return fmt.Errorf("Failed to connect to %s: %v", node, err)
		}

		defer conn.Close()
		connections[i] = conn
	}

	// Sending chunks to nodes in round-robin fashion
	for i, chunk := range *chunks {
		nodeIndex := i % len(nodes) // It select the node index as per the round robin fashion.
		conn := connections[nodeIndex]

		chunkData := chunkData{
			Filename: name,
			FileId: chunk.Id,
			Data: chunk.Data,
		}

		// err := sendChunkToNode(conn, &chunk, name)

		err := chunkData.sendChunkToDataNode(conn)

		if err != nil {
			continue
		}

		if err != nil {
			return fmt.Errorf("error sending chunk %s to node %s: %v", chunk.Id, nodes[nodeIndex], err)
		}

		location[fmt.Sprintf("Node%d", nodeIndex)] = chunk.Id
	}

	fmt.Println(location)

	metaData := metaData{
		Filename:      name,
		FileExtension: extension,
		Location:      location,
	}

	fmt.Println(metaData)

	if err := metaData.sendMetaData(); err != nil {
		return err
	}

	return nil
}

func (md *metaData) sendMetaData() error {

	conn, err := net.Dial("tcp", nameNode)
	fmt.Println("meta data location:", md.Location)
	
	if err != nil {
		return fmt.Errorf("Error connecting to the Name node server")
	}
	
	defer conn.Close()
	// conn.Write([]byte("POST\n")) // This is causing some error - deadlock situation instead of that using bufio.NewWriter()

	writer := bufio.NewWriter(conn)
	_, err = writer.Write([]byte("POST\n"))
	
	if err != nil {
		return fmt.Errorf("Error sending post request to the server %w", err)
	}
	writer.Flush()

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(md); err != nil {
		return fmt.Errorf("Error encoding metadata: %v", err)
	}

	// ensuring that the server will get's the EOF (End of file) signal properly
	if err = conn.(*net.TCPConn).CloseWrite(); err != nil {
		return fmt.Errorf("Error closing the write side of connection: %v", err)
	}

	reader := bufio.NewReader(conn)
	serverResponse, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	serverResponse = strings.TrimSpace(serverResponse)
	if serverResponse != "200" {
		return fmt.Errorf("Response is not 200", serverResponse)
	}

	fmt.Println("Metadata saved successfully to the namenode")
	return nil
}

func (chunkData *chunkData) sendChunkToDataNode(conn net.Conn) error {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(chunkData)

	if err != nil {
		return fmt.Errorf("error sending chunks to data node", err)
	}

	fmt.Println("ChunkId", chunkData.FileId)
	fmt.Println("ChunkName", chunkData.Filename)
	fmt.Println("chunkData length", len(chunkData.Data))

	fmt.Printf("%s sent to the data node server successfully\n", chunkData.FileId)
	return nil
}

func sendChunkToNode(conn net.Conn, chunk *chunk.Chunk, name string) error {
	_, err := conn.Write([]byte(chunk.Id + "\n"))
	if err != nil {
		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
	}
	_, err = conn.Write(chunk.Data)
	if err != nil {
		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
	}

	fmt.Printf("Sent chunk %s to node\n", chunk.Id)
	return nil
}

func GetChunks(filename string) (*[]chunk.Chunk, error) {
	conn, err := net.Dial("tcp", nameNode)
	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to namenode", err)
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("File name is empty")
	}

	conn.Write([]byte("GET\n" + filename + "\n"))

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if response != "200" {
		return nil, fmt.Errorf("Response is not OK", response)
	}

	decoder := gob.NewDecoder(conn)
	var recievedData metaData
	err = decoder.Decode(&recievedData)

	if err != nil {
		return nil, fmt.Errorf("Error occured while decoding the data", err)
	}

	fmt.Println("metadata is :", recievedData)
	return nil, nil
}

// func (metaData *metaData) FetchChunks() (*[]chunk.Chunk, error) {
// 	mappingData := make(map[string]int)
// 	mappingData["node0"] = 0
// 	mappingData["node1"] = 1
// 	mappingData["node2"] = 2

// }
