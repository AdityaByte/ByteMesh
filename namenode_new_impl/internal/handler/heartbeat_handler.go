package handler

import (
	"fmt"

	"github.com/AdityaByte/namenode/internal/payloads"
)

func HeartBeatHandler(datanodes *payloads.RegisteredDataNodes, heartbeat *payloads.HeartBeat) {
	for _, nodes := range datanodes.Nodes {
		if nodes.Name == heartbeat.NodeName {
			// If the name is same then we have to just update the timestamp.
			fmt.Println("Time stamp difference - ", heartbeat.TimeStamp-nodes.TimeStamp)
			nodes.TimeStamp = heartbeat.TimeStamp
		} else {
			fmt.Println("I am outside the if block")
		}
	}
}
