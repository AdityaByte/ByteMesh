// package utils

// import (
// 	"fmt"
// 	"net"
// )

// const (
// 	node1 = ":9001"
// 	node2 = ":9002"
// 	node3 = ":9003"
// )

// var Nodes = []string{node1, node2, node3}

// func CreateConnectionPool() ([]net.Conn, error) {

// 	connections := make([]net.Conn, len(Nodes))

// 	for i, node := range Nodes {
// 		conn, err := net.Dial("tcp", node)
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to connect to %s", node)
// 		}

// 		connections[i] = conn
// 	}

// 	return connections, nil
// }

package utils

import (
	"fmt"
	"net"

	"github.com/AdityaByte/bytemesh/payload"
)

func CreateConnectionPool(aliveNodes payload.RegisteredDataNodes) (payload.RegisteredDataNodes, error) {
	// Have to create the connection to each and every alive node and dump the connection object to the
	// existed struct.
	for i := range aliveNodes.Nodes {
		node := &aliveNodes.Nodes[i]
		// Although assuming currently all the nodes are running to the localhost.
		// Moreover we will implement that host logic after some time.
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", node.Port))
		if err != nil {
			return payload.RegisteredDataNodes{}, fmt.Errorf("DataNode running on port(%d) is refusing to connect, %v", node.Port, err)
		}
		node.Conn = conn
	}

	// When everything goes correctly we just need to forward the aliveNodes with the connection.
	return aliveNodes, nil
}