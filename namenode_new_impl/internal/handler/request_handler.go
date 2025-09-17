package handler

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/AdityaByte/namenode/internal/database"
	"github.com/AdityaByte/namenode/internal/model"
	"github.com/AdityaByte/namenode/internal/payloads"
	"go.mongodb.org/mongo-driver/bson"
)

var aliveNodeThreshold = 30 // 30 seconds

/*
The main logic of fetching out the file data from different datanodes resides here.
The fault tolerant system logic also resides here.
*/
func HandleGetRequest(ctx context.Context, mongoRepo database.MongoRepository, datanodes *payloads.RegisteredDataNodes, fullFileName string) (*payloads.MetaData, error) {

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

	var metadata model.MetaData
	if err := data.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("Failed to decode the data, %v", err)
	}

	// Here we need to check the nodes availability too.
	aliveDataNodes := checkNodeAvailiability(metadata, *datanodes)

	payload_metadata := payloads.MetaData{
		Filename:      metadata.Filename,
		FileExtension: metadata.FileExtension,
		ActualSize:    metadata.ActualSize,
		Location:      metadata.Location,
		AliveNodes:    aliveDataNodes,
	}

	return &payload_metadata, nil
}

func checkNodeAvailiability(metadata model.MetaData, allNodes payloads.RegisteredDataNodes) payloads.RegisteredDataNodes {
	// Just need to check the nodes that do exists in the chunk location they are alive or not.
	var aliveNodes payloads.RegisteredDataNodes

	// Node for fast lookup.
	nodeMap := make(map[string]payloads.DataNode)
	for _, node := range allNodes.Nodes {
		nodeMap[node.Name] = node
	}

	// seen tracking duplicates.
	seen := make(map[string]bool)

	currentTimeStamp := time.Now().Unix()

	for _, value := range metadata.Location {
		if seen[value] {
			continue
		}

		seen[value] = true

		if node, exists := nodeMap[value]; exists {
			if currentTimeStamp-node.TimeStamp <= int64(aliveNodeThreshold) {
				aliveNodes.Nodes = append(aliveNodes.Nodes, node)
			}
		}
	}

	return aliveNodes
}

func HandlePostRequest(ctx context.Context, mongoRepo *database.MongoRepository, metadata model.MetaData) error {
	if err := emptyFieldsChecker(metadata); err != nil {
		return err
	}

	if _, err := mongoRepo.Collection.InsertOne(ctx, metadata); err != nil {
		return fmt.Errorf("Failed to insert metadata to the database, %v", err)
	}

	return nil
}

func emptyFieldsChecker(metadata model.MetaData) error {
	if strings.TrimSpace(metadata.Filename) == "" {
		return fmt.Errorf("ERROR: Filename of metadata is empty")
	} else if strings.TrimSpace(metadata.FileExtension) == "" {
		return fmt.Errorf("ERROR: FileExtension of metadata is empty")
	} else {
		for key, value := range metadata.Location {
			if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
				return fmt.Errorf("ERROR: Location fields cannot be empty")
			}
		}
	}
	return nil
}

// checkNodeAvailiablity Takes pre-existed all nodes and check at which time the last heartbeat signal reaches and returns
// a list of alive nodes.
// func checkNodeAvailiability(metaData model.MetaData, allNodes payloads.RegisteredDataNodes) (payloads.RegisteredDataNodes) {
// 	// Just need to check the nodes that do exists in the chunk location they are alive or not.
// 	var aliveNodes payloads.RegisteredDataNodes

// 	var tempNodeData []string

// 	for _, value := range metaData.Location {
// 		if !contains(tempNodeData, value) {
// 			// Now we have to check that the node is alive or not.
// 			for _, nodeInAliveNode := range allNodes.Nodes {
// 				if value == nodeInAliveNode.Name {
// 					currentTimeStamp := time.Now().Unix()
// 					if currentTimeStamp - nodeInAliveNode.TimeStamp <= int64(aliveNodeThreshold) {
// 						aliveNodes.Nodes = append(aliveNodes.Nodes, nodeInAliveNode)
// 					}
// 				}
// 			}
// 			tempNodeData = append(tempNodeData, value)
// 		}
// 	}

// 	return aliveNodes
// }

// func contains(tempNodeData []string, value string) bool {

// 	if len(tempNodeData) == 0 {
// 		return false
// 	}

// 	for _, node := range tempNodeData {
// 		if node == value {
// 			return true
// 		}
// 	}
// 	return false
// }

// This function ususally takes name of the DataNode.
// And it returns the connection and also checks that it is alive or not.
// func mapDataNodeandCheckAlive(nodeName string, datanodes payloads.RegisteredDataNodes) (payloads.RegisteredDataNodes) {
// 	var aliveNodes []payloads.DataNode
// 	for _, node := range datanodes.Nodes {
// 		if nodeName == node.Name {
// 			// Then we have to check that it is alive or not.
// 			currentTimeStamp := time.Now().Unix()
// 			if currentTimeStamp-node.TimeStamp <= int64(aliveNodeThreshold) {
// 				aliveNodes = append(aliveNodes, payloads.DataNode{
// 					Name: nodeName,
// 				})
// 			}
// 		}
// 	}
// 	return nil
// }
