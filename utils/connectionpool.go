package utils

import (
	"fmt"
	"net"
)

const (
	node1    = ":9001"
	node2    = ":9002"
	node3    = ":9003"
	nameNode = ":9004"
)

var Nodes = []string{node1, node2, node3}


func CreateConnectionPool() ([]net.Conn, error) {

	connections := make([]net.Conn, len(Nodes))

	for i, node := range Nodes {
		conn, err := net.Dial("tcp", node)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to %s", node)
		}

		connections[i] = conn
	}

	return connections, nil
}