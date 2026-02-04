package handler

import (
	"github.com/AdityaByte/namenode/internal/payloads"
	"github.com/AdityaByte/namenode/logger"
)

func HeartBeatHandler(datanodes *payloads.RegisteredDataNodes, heartbeat *payloads.HeartBeat) {
	for i, node := range datanodes.Nodes {
		if node.Name == heartbeat.NodeName {
			// If the name is same then we have to just update the timestamp.
			logger.InfoLogger.Println("Time stamp difference - ", heartbeat.TimeStamp-node.TimeStamp)
			datanodes.Nodes[i].TimeStamp = heartbeat.TimeStamp
		}
	}
}
