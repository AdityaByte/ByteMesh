package coordinator

import (
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
	FileId string
	Data []byte
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
			return fmt.Errorf("error sending chunk %s to node %s: %v", chunk.Id, nodes[nodeIndex], err)
		}

		location[fmt.Sprintf("Node%d", nodeIndex)] = chunk.Id
	}

	fmt.Println(location)

	metaData := metaData{
		Filename: name,
		FileExtension: extension,
		Location: location,
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

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(md)

	if err != nil {
		fmt.Errorf("Error encoding the metadata", err)
	}

	fmt.Println("MetaData sent successfully to the Name node server")
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

// func sendChunkToNode(conn net.Conn, chunk *chunk.Chunk, name string) error {
// 	_, err := conn.Write([]byte(chunk.Id + "\n"))
// 	if err != nil {
// 		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
// 	}
// 	_, err = conn.Write(chunk.Data)
// 	if err != nil {
// 		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
// 	}

// 	fmt.Printf("Sent chunk %s to node\n", chunk.Id)
// 	return nil
// }
