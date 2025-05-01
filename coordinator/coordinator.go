package coordinator

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/AdityaByte/bytemesh/models"
	"github.com/AdityaByte/bytemesh/utils"
)

const nameNode = ":9004"

func SendChunks(chunks *[]models.Chunk, filename string, filesize float64) error {

	if filename == "" {
		return fmt.Errorf("File name does not exists.")
	}

	log.Println("Actual file name:", filename)
	fileData := strings.Split(filename, ".")
	name := fileData[0]
	log.Println("Prefix:", name)
	extension := fileData[1]

	// Before creating the connection we have to check the health
	// of the Datanodes and the Namenodes so that all are working correctly or not.

	errors := HealthChecker([]string{
		":9001",
		":9002",
		":9003",
		":9004",
	})

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		return fmt.Errorf("Something goes wrong one or more nodes are unhealthy.")
	} else {
		log.Println("All Nodes are healthy")
	}

	connections, err := utils.CreateConnectionPool()
	if err != nil {
		return err
	}

	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	fileLocation := make(map[string]string)

	// Sending chunks to nodes in round-robin fashion
	for i, chunk := range *chunks {

		log.Println("Iteration:", i, "ChunkId:", chunk.Id)

		nodeIndex := i % len(utils.Nodes) // It select the node index as per the round robin fashion.
		conn := connections[nodeIndex]
		if conn == nil {
			fmt.Println("connection is null")
		}

		log.Println("Node index for iteration i:", i, "is", nodeIndex)

		chunkData := models.ChunkData{
			Filename: name,
			FileId:   chunk.Id,
			Data:     chunk.Data,
		}

		err := sendChunkToDataNode(conn, &chunkData)

		if err != nil {
			log.Println("Error occured:", err)
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
		ActualSize:    filesize,
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

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Sending POST Request
	if _, err := writer.WriteString("POST" + "\n"); err != nil {
		return fmt.Errorf("Failed to send the post request to datanode %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Flush Failed: %v", err)
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(chunkData); err != nil {
		return fmt.Errorf("Failed to encode chunk : %v", err)
	}

	log.Println("ChunkId", chunkData.FileId)
	log.Println("ChunkName", chunkData.Filename)
	log.Println("chunkData length", len(chunkData.Data))

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Response", string(response))
		return fmt.Errorf("Failed to read the response : %v", err)
	}

	if strings.TrimSpace(response) != "OK" {
		return fmt.Errorf("Server error : %s", response)
	}

	return nil
}

func FetchChunks(metaData *models.MetaData) (*[]byte, error) {

	mappingData := map[string]string{
		"Node0": ":9001",
		"Node1": ":9002",
		"Node2": ":9003",
	}

	location := metaData.Location

	var finalData bytes.Buffer

	// Here firstly we need to sort the keys before passing them to the for range loop
	var keys []string
	for k := range location {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Now we passed the sorted keys to it.
	for _, key := range keys {
		value := location[key]

		log.Println("Key:", key, "Value:", value)
		log.Println("Chunk stored in:", mappingData[value])

		data, err := getChunkFromNode(metaData.Filename, key, mappingData[value])
		if err != nil {
			return nil, err
		}

		finalData.Write(data)
	}

	// This is the old code without the sorting logic.
	// for key, value := range location {
	// 	// here we get the key which is the chunk1 ok so we derive in which node is being stored so we make a connection to
	// 	// the particular node and share out the name of the chunk ok means its id which is the name and we get the data which was we being stored to
	// 	// the bytes.Buffer and at the very after end we rename the file to the actual name and download it in the downloaded folder.

	// 	// key -> chunk id
	// 	// value -> in which node is been stored

	// 	fmt.Println("value is", value)
	// 	fmt.Println(mappingData[value])

	// 	data, err := getChunkFromNode(metaData.Filename, key, mappingData[value])
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	finalData.Write(data)
	// }

	sendingData := finalData.Bytes()

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

	log.Println("Node address we are passing is :", nodeAddr)

	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s : %v", nodeAddr, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	_, err = writer.WriteString("GET\n" + filename + "\n" + chunkId + "\n") // We have to manually add the newline character cause in go it doesn't add it automatically.
	if err != nil {
		return nil, fmt.Errorf("Write Failed: %v", err)
	}
	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("Flush Failed: %v", err)
	}

	// reader := bufio.NewReader(conn)
	// data := make([]byte, 500*1024)
	// n, err := reader.Read(data)

	// var buf bytes.Buffer
	// _, err = io.Copy(&buf, conn)

	// Here firstly we have to read the chunksize

	// if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
	// 	return nil, fmt.Errorf("Failed to set the read deadline: %v", err)
	// }

	var chunkSize uint32
	if err := binary.Read(reader, binary.BigEndian, &chunkSize); err != nil { // Always read data from reader.
		return nil, fmt.Errorf("Failed to read the chunk size: %v", err)
	}

	chunkData := make([]byte, chunkSize)
	_, err = io.ReadFull(reader, chunkData)

	if err != nil {
		return nil, fmt.Errorf("Failed to read the chunk data: %v", err)
	}

	return chunkData, nil
}
