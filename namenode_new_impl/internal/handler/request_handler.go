package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"path/filepath"
	"sort"
	"time"

	"github.com/AdityaByte/namenode/internal"
	"github.com/AdityaByte/namenode/internal/model"
	"github.com/AdityaByte/namenode/internal/payloads"
	"go.mongodb.org/mongo-driver/bson"
)

var aliveNodeThreshold = 30 // 30 seconds

/*
The main logic of fetching out the file data from different datanodes resides here.
The fault tolerant system logic also resides here.
*/
func HandleGetRequest(ctx context.Context, mongoRepo internal.MongoRepository, datanodes *payloads.RegisteredDataNodes, fullFileName string) ([]byte, error) {

	// Seperating the extension and the filename
	fileExtension := filepath.Ext(fullFileName)
	filename := fullFileName[:len(fullFileName)-len(fileExtension)]

	if fileExtension == "" {
		return nil, fmt.Errorf("No extension exists")
	} else if filename == "" {
		return nil, fmt.Errorf("No filename exists")
	}

	filters := bson.M{"filename": filename, "fileextension": fileExtension}
	data := mongoRepo.Collection.FindOne(ctx, filters)
	if data == nil {
		return nil, fmt.Errorf("No file exists of name: %s", fullFileName)
	}

	var metaData *model.MetaData
	if err := data.Decode(&metaData); err != nil {
		return nil, fmt.Errorf("Failed to decode the data, %v", err)
	}

	// Now we have to use the existing connection pool connection and fetch the data chunks.
	var keys []string

	for k := range metaData.Location {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var final_file_data []byte

	for _, key := range keys {
		node := metaData.Location[key]
		conn := mapDataNodeandCheckAlive(node, *datanodes)
		if conn == nil {
			return nil, fmt.Errorf("Node name: %s is not alive", node)
		}
		// Implement the else replication after some time.
		// If the connection is alive then we have to send the get request to the node.

		// Sending the request verb first.
		_, err := conn.Write([]byte("GET\n"))
		if err != nil {
			return nil, fmt.Errorf("Failed to write the data to the datanode, %v", err)
		}

		// Now we need to send the json payload.
		get_request_payload := payloads.GetRequest{
			FileName: filename,
			ChunkId:  key,
		}

		// Now we need to send that over the network.
		encoded_json_data, err := json.Marshal(get_request_payload)
		if err != nil {
			return nil, fmt.Errorf("Failed to encode the data to json, %v", err)
		}

		if _, err := conn.Write(encoded_json_data); err != nil {
			return nil, fmt.Errorf("Failed to send the data to the datanode, %v", err)
		}

		// Now we need to recieve the data.
		decoder := json.NewDecoder(conn)

		var chunk payloads.Chunk
		if err := decoder.Decode(&chunk); err != nil {
			return nil, fmt.Errorf("Failed to decode the json data, %v", err)
		}

		// Since we decode the data we have to just add the byte data
		final_file_data = append(final_file_data, chunk.Data...)
	}

	// If everything goes correctly we gets the data.

	// Now have to check the length of the data is correct or not.
	if len(final_file_data) != int(math.Round(metaData.ActualSize)) {
		return nil, fmt.Errorf("Corrupted data reciecved")
	}

	return final_file_data, nil
}

// This function ususally takes name of the DataNode.
// And it returns the connection and also checks that it is alive or not.
func mapDataNodeandCheckAlive(nodeName string, datanodes payloads.RegisteredDataNodes) (conn net.Conn) {
	for _, node := range datanodes.Nodes {
		if nodeName == node.Name {
			// Then we have to check that it is alive or not.
			currentTimeStamp := time.Now().Unix()
			if currentTimeStamp-node.TimeStamp <= int64(aliveNodeThreshold) {
				return node.Conn
			}
		}
	}
	return nil
}
