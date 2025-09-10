package handler

import "github.com/AdityaByte/namenode/internal/payloads"

func HeartBeatHandler(datanodes *payloads.RegisteredDataNodes, heartbeat *payloads.HeartBeat) {
	for _, nodes := range datanodes.Nodes {
		if nodes.Name == heartbeat.NodeName {
			// If the name is same then we have to just update the timestamp.
			nodes.TimeStamp = heartbeat.TimeStamp
		}
	}
}
