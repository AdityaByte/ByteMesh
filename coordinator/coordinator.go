package coordinator

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strings"

	"github.com/AdityaByte/bytemesh/models"
	"github.com/AdityaByte/bytemesh/utils"
)

const nameNode = ":9004"

func SendChunks(chunks *[]models.Chunk, filename string) error {

	if filename == "" {
		return fmt.Errorf("File name does not exists.")
	}

	fmt.Println("original file name :", filename)
	fileData := strings.Split(filename, ".")
	name := fileData[0]
	extension := fileData[1]

	connections, err := utils.CreateConnectionPool()
	if err != nil {
		return err
	}

	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	} ()

	// nodes := []string{node1, node2, node3}

	// connections := make([]net.Conn, len(nodes))

	// for i, node := range nodes {
	// 	conn, err := net.Dial("tcp", node)
	// 	if err != nil {
	// 		return fmt.Errorf("Failed to connect to %s: %v", node, err)
	// 	}

	// 	defer conn.Close()
	// 	connections[i] = conn
	// }

	fileLocation := make(map[string]string)

	// Sending chunks to nodes in round-robin fashion
	for i, chunk := range *chunks {
		nodeIndex := i % len(utils.Nodes) // It select the node index as per the round robin fashion.
		conn := connections[nodeIndex]

		chunkData := models.ChunkData{
			Filename: name,
			FileId:   chunk.Id,
			Data:     chunk.Data,
		}

		// err := sendChunkToNode(conn, &chunk, name)

		err := sendChunkToDataNode(conn, &chunkData)

		if err != nil {
			continue
		}

		if err != nil {
			return fmt.Errorf("error sending chunk %s to node %s: %v", chunk.Id, utils.Nodes[nodeIndex], err)
		}

		// location[fmt.Sprintf("Node%d", nodeIndex)] = chunk.Id
		fileLocation[chunk.Id] = fmt.Sprintf("Node%d", nodeIndex)
	}

	fmt.Println(fileLocation)

	metaData := models.MetaData{
		Filename:      name,
		FileExtension: extension,
		Location:      fileLocation,
	}

	fmt.Println(metaData)

	if err := sendMetaData(&metaData); err != nil {
		return err
	}

	return nil
}

func sendMetaData(md *models.MetaData) error {

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

func sendChunkToDataNode(conn net.Conn, chunkData *models.ChunkData) error {
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

func sendChunkToNode(conn net.Conn, chunk *models.Chunk, name string) error {
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

func FetchChunks(metaData *models.MetaData) (*[]byte, error) {

	mappingData := map[string]string{
		"Node0": ":9001",
		"Node1": ":9002",
		"Node2": ":9003",
	}

	location := metaData.Location

	var allData bytes.Buffer


	for key, value := range location {
		// here we get the key which is the chunk1 ok so we derive in which node is being stored so we make a connection to 
		// the particular node and share out the name of the chunk ok means its id which is the name and we get the data which was we being stored to
		// the bytes.Buffer and at the very after end we rename the file to the actual name and download it in the downloaded folder.
	
		// key -> chunk id
		// value -> in which node is been stored

		fmt.Println("value is", value)
		fmt.Println(mappingData[value])

		data, err := getChunkFromNode(metaData.Filename, key, mappingData[value])
		if err != nil {
			return nil, err
		}

		allData.Write(data)
	}

	sendingData := allData.Bytes()
	
	return &sendingData, nil
}

// dekh dude apne pass kya hai keys hai aur values hai ok toh us hisab se mai data ko fetch krunga ok...
// ek for each loop chalayenge apan keys ke liye

// logic 
// map[string]string :::::::::::------------------------------------------------------->
//  node0 -> chunk1 and node1 -> chunk2 
// chunk1 -> node0 
// chunk2 -> node1 
// chunk3 -> node2 like this so we have to change that all 

func getChunkFromNode(filename string, chunkId string, nodeAddr string) ([]byte, error) {

	fmt.Println("Node address we are passing is :", nodeAddr)

	conn, err := net.Dial("tcp", nodeAddr)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s : %v", nodeAddr, err)
	}

	writer := bufio.NewWriter(conn)
	writer.WriteString("GET" + "\n" + filename + "\n" + chunkId + "\n") // We have to manually add the newline character cause in go it doesn't add it automatically.
	writer.Flush() // For sending data immediately.

	reader := bufio.NewReader(conn)
	data := make([]byte, 30*1024)
	n, err := reader.Read(data)

	if err != nil {
		return nil, fmt.Errorf("Failed to read data %v", err)
	}

	return data[:n],nil
}