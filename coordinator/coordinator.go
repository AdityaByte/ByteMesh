package coordinator

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/models"
	"github.com/AdityaByte/bytemesh/payload"
	"github.com/AdityaByte/bytemesh/utils"
)

const namenode = ":9004"

// SendChunks send the chunk data and metadata to namenode.
func SendChunks(chunks *[]models.Chunk, filename string, filesize float64, username string) error {

	namenodeConnection, err := net.Dial("tcp", namenode)
	if err != nil {
		return fmt.Errorf("Failed to make the connection to the namenode server, %v", err)
	}
	defer namenodeConnection.Close()

	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("Empty filename recieved at sendchunks function.")
	}

	logger.InfoLogger.Println("Full filename:", filename)
	extension := path.Ext(filename)
	name := filename[:(len(filename) - len(extension))]

	// Fetching alive datanodes data from the namenode.
	// Sending the HEATLH verb.
	if _, err := namenodeConnection.Write([]byte("HEALTH\n")); err != nil {
		return fmt.Errorf("Failed to send the health request to the namenode server.")
	}

	var registrationPayload payload.RegisteredDataNodes
	gobDecoder := gob.NewDecoder(namenodeConnection)
	gobDecoder.Decode(&registrationPayload)

	// Since we got the alive nodes now we need to check they may not be empty if they we will throw an info message.
	if !isNodeHealthGood(registrationPayload) {
		return fmt.Errorf("No datanode is alive")
	}

	// If the health is good we need to make the connection.
	aliveNodeConnections, err := utils.CreateConnectionPool(registrationPayload)
	if err != nil {
		return err
	}

	defer func() {
		for i := range aliveNodeConnections.Nodes {
			node := &aliveNodeConnections.Nodes[i]
			if node.Conn != nil {
				node.Conn.Close()
			}
		}
	}()

	fileLocation := make(map[string]string)

	// Sending chunks to nodes in round-robin fashion
	for i, chunk := range *chunks {

		logger.InfoLogger.Println("Iteration:", i, "ChunkId:", chunk.Id)

		nodeIndex := i % len(aliveNodeConnections.Nodes) // It select the node index as per the round robin fashion.
		conn := aliveNodeConnections.Nodes[nodeIndex].Conn
		if conn == nil {
			logger.InfoLogger.Println("Nil Connection")
		}

		logger.InfoLogger.Println("Node index for iteration i:", i, "is", nodeIndex)

		chunkData := models.ChunkData{
			Filename: name,
			FileId:   chunk.Id,
			Data:     chunk.Data,
		}

		err := sendChunkToDataNode(conn, &chunkData)

		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}

		if err != nil {
			return fmt.Errorf("Failed to send chunk %s to node %s: %v", chunk.Id, aliveNodeConnections.Nodes[nodeIndex], err)
		}

		fileLocation[chunk.Id] = fmt.Sprintf("Node%d", nodeIndex)
	}

	fmt.Println(fileLocation)

	// Now before sending the chunks we can read the authenticated username.
	if username == "" {
		data, err := os.ReadFile("../.auth/.cred")
		if err != nil {
			return fmt.Errorf("ERROR: Failed to read the file")
		}

		// Else we gets the file data.
		username = string(data)
		if utils.CheckEmptyField(username) {
			return fmt.Errorf("ERROR: Fetched username is empty")
		}
	}

	metaData := models.MetaData{
		Owner:         username,
		Filename:      name,
		FileExtension: extension,
		UploadDate:    time.Now(),
		ActualSize:    filesize,
		Location:      fileLocation,
	}

	fmt.Println(metaData)

	if err := sendMetaData(namenodeConnection, &metaData); err != nil {
		return err
	}

	return nil
}

func sendMetaData(conn net.Conn, metadata *models.MetaData) error {

	logger.InfoLogger.Println("meta data location:", metadata.Location)

	writer := bufio.NewWriter(conn)
	if _, err := writer.Write([]byte("POST\n")); err != nil {
		return fmt.Errorf("Error sending post request to the server %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Failed to flush the post request to the namenode server, %v", err)
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("Error encoding metadata: %v", err)
	}

	// ensuring that the server will get's the EOF (End of file) signal properly
	if err := conn.(*net.TCPConn).CloseWrite(); err != nil {
		return fmt.Errorf("Error closing the write side of connection: %v", err)
	}

	reader := bufio.NewReader(conn)
	serverResponse, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	serverResponse = strings.TrimSpace(serverResponse)
	if serverResponse != "201" {
		return fmt.Errorf("Response is not 201 %s", serverResponse)
	}

	logger.InfoLogger.Println("Metadata saved successfully to the namenode")
	return nil
}

func sendChunkToDataNode(conn net.Conn, chunkData *models.ChunkData) error {

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Sending POST Request
	if _, err := writer.WriteString("POST\n"); err != nil {
		return fmt.Errorf("Failed to send the post request to datanode %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Flush Failed: %v", err)
	}

	jsonEncodedData, err := json.Marshal(chunkData)
	if err != nil {
		return fmt.Errorf("Failed to encode the json data, %v", err)
	}

	// Adding a delimetter at the end so that it can read the exact part.
	jsonEncodedData = append(jsonEncodedData, '\n')

	if _, err := conn.Write(jsonEncodedData); err != nil {
		return fmt.Errorf("Failed to send the chunkdata to datanode, %v", err)
	}

	logger.InfoLogger.Println("ChunkId", chunkData.FileId)
	logger.InfoLogger.Println("ChunkName", chunkData.Filename)
	logger.InfoLogger.Println("chunkData length", len(chunkData.Data))

	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Failed to read the response : %v", err)
	}

	if strings.TrimSpace(response) != "201" {
		return fmt.Errorf("Server error : %s", response)
	}

	return nil
}

// FetchChunks: Do a request to the namenode fetch the metadata then according to the metadata makes connection to the datanodes and fetch the chunk data.
func FetchChunks(filename string) (*[]byte, error) {

	// Firstly we had to make the connection to the namenode server for fetching the chunk data that have been stored.
	namenodeConnection, err := net.Dial("tcp", namenode)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the connection to namenode server, %v", err)
	}

	if _, err := namenodeConnection.Write([]byte("GET\n")); err != nil {
		return nil, fmt.Errorf("Failed to send the request verb to namenode server, %v", err)
	}

	// Once we have sent the request now we need to send the filename.
	namenodeConnection.Write([]byte(filename + "\n"))

	// Now I need to wait for the output for a definite period of time.
	var metadata payload.MetaData
	decoder := gob.NewDecoder(namenodeConnection)
	decoder.Decode(&metadata)

	// Now I need to check the alivenodes length should be greater than 0.
	// and also the other fields can't be empty.
	if strings.TrimSpace(metadata.Filename) == "" || strings.TrimSpace(metadata.FileExtension) == "" {
		return nil, fmt.Errorf("Empty data recieved from the namenode")
	} else if len(metadata.AliveNodes.Nodes) == 0 {
		return nil, fmt.Errorf("Failed to fetch the file since no alive nodes are found")
	} else if len(metadata.Location) == 0 {
		return nil, fmt.Errorf("Failed to get the location, empty location founded.")
	}

	// If everything goes correctly we need to fetch the chunks from the nodes.
	// Note: Port shouldn't be saved to database because it is dynamic.
	// This time it has been given the port too.

	var keys []string
	for k := range metadata.Location {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// NodeAddressLookup
	var nodeAddr map[string]payload.DataNode // Key - DataNodeName Value - Port at which it is running.
	for _, datanode := range metadata.AliveNodes.Nodes {
		nodeAddr[datanode.Name] = datanode
	}

	// Buffer for the final data.
	var finalData bytes.Buffer

	// Now we have to make the connection to the datanodes and fetch out the chunks.
	for _, key := range keys {
		datanodeName := metadata.Location[key]
		datanode := nodeAddr[datanodeName]
		if strings.TrimSpace(datanode.Name) == "" || datanode.Port == 0 {
			return nil, fmt.Errorf("Node name %s is not alive, try again later", datanodeName)
		}

		// Now we need to create the connection to the datanode and fetch out the chunk.
		data, err := getChunkFromNode(filename, key, fmt.Sprintf(":%d", datanode.Port))
		if err != nil {
			return &data, err
		}

		finalData.Write(data)
	}

	data := finalData.Bytes()

	return &data, nil

	// mappingData := map[string]string{
	// 	"Node0": ":9001",
	// 	"Node1": ":9002",
	// 	"Node2": ":9003",
	// }

	// location := metaData.Location

	// var finalData bytes.Buffer

	// // Here firstly we need to sort the keys before passing them to the for range loop
	// var keys []string
	// for k := range location {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)

	// // Now we passed the sorted keys to it.
	// for _, key := range keys {
	// 	value := location[key]

	// 	logger.InfoLogger.Println("Key:", key, "Value:", value)
	// 	logger.InfoLogger.Println("Chunk stored in:", mappingData[value])

	// 	data, err := getChunkFromNode(metaData.Filename, key, mappingData[value])
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	finalData.Write(data)
	// }

	// // This is the old code without the sorting logic.
	// // for key, value := range location {
	// // 	// here we get the key which is the chunk1 ok so we derive in which node is being stored so we make a connection to
	// // 	// the particular node and share out the name of the chunk ok means its id which is the name and we get the data which was we being stored to
	// // 	// the bytes.Buffer and at the very after end we rename the file to the actual name and download it in the downloaded folder.

	// // 	// key -> chunk id
	// // 	// value -> in which node is been stored

	// // 	fmt.Println("value is", value)
	// // 	fmt.Println(mappingData[value])

	// // 	data, err := getChunkFromNode(metaData.Filename, key, mappingData[value])
	// // 	if err != nil {
	// // 		return nil, err
	// // 	}

	// // 	finalData.Write(data)
	// // }

	// sendingData := finalData.Bytes()

	// return &sendingData, nil
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

	logger.InfoLogger.Printf("getChunkFromNode[INFO]: Node Address - %s\n", nodeAddr)

	datanodeConnection, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s : %v", nodeAddr, err)
	}
	defer datanodeConnection.Close()

	getRequestPayload := payload.GetRequest{
		Filename: filename,
		ChunkId:  chunkId,
	}

	jsonGetRequestPayload, err := json.Marshal(getRequestPayload)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode the get request payload to json, %v", err)
	}

	// Setting up the delimetter to the json encoded data.
	jsonGetRequestPayload = append(jsonGetRequestPayload, '\n')

	if _, err := datanodeConnection.Write(jsonGetRequestPayload); err != nil {
		return nil, fmt.Errorf("Failed to send the get request payload to the datanode, %v", err)
	}

	reader := bufio.NewReader(datanodeConnection)

	rawData, err := reader.ReadString('\n')
	var getResponse payload.Response
	if err := json.Unmarshal([]byte(rawData), &getResponse); err != nil {
		return nil, fmt.Errorf("Failed to decode the get response data, %v", err)
	}

	// Now we need to check the response type is success or failure.

	if strings.TrimSpace(getResponse.Type) != "SUCCESS" {
		logger.InfoLogger.Printf("Get response type: %s\n", getResponse.Type)
		return nil, fmt.Errorf(getResponse.Message.(string))
	}

	// Success Now we need to return the data.
	return getResponse.Message.([]byte), nil

	// reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)

	// _, err = writer.WriteString("GET\n" + filename + "\n" + chunkId + "\n") // We have to manually add the newline character cause in go it doesn't add it automatically.
	// if err != nil {
	// 	return nil, fmt.Errorf("Write Failed: %v", err)
	// }
	// if err := writer.Flush(); err != nil {
	// 	return nil, fmt.Errorf("Flush Failed: %v", err)
	// }

	// reader := bufio.NewReader(conn)
	// data := make([]byte, 500*1024)
	// n, err := reader.Read(data)

	// var buf bytes.Buffer
	// _, err = io.Copy(&buf, conn)

	// Here firstly we have to read the chunksize

	// if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
	// 	return nil, fmt.Errorf("Failed to set the read deadline: %v", err)
	// }

	// var chunkSize uint32
	// if err := binary.Read(reader, binary.BigEndian, &chunkSize); err != nil { // Always read data from reader.
	// 	return nil, fmt.Errorf("Failed to read the chunk size: %v", err)
	// }

	// chunkData := make([]byte, chunkSize)
	// _, err = io.ReadFull(reader, chunkData)

	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to read the chunk data: %v", err)
	// }

	// return chunkData, nil
}
