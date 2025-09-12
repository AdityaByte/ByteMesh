package handler

import (
	"time"

	"github.com/AdityaByte/namenode/internal/payloads"
)

// GetAliveNodes returns information about the currently alive nodes.
func GetAliveNodes(allNodes payloads.RegisteredDataNodes) payloads.RegisteredDataNodes {

	aliveNodeThreshold := 30
	currentTimeStamp := time.Now().Unix()

	var currentlyAliveNodes payloads.RegisteredDataNodes

	for _, node := range allNodes.Nodes {
		if currentTimeStamp - node.TimeStamp <= int64(aliveNodeThreshold) {
			currentlyAliveNodes.Nodes = append(currentlyAliveNodes.Nodes, node)
		}
	}

	return currentlyAliveNodes
}