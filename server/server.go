package server

import (
	"fmt"
	"net"

	"github.com/AdityaByte/bytemesh/chunk"
)

const (
	node1 = ":9001"
	node2 = ":9002"
	node3 = ":9003"
)


func SendChunks(chunks *[]chunk.Chunk) error {

	nodes := []string{node1, node2, node3}

	connections := make([]net.Conn, len(nodes))

	for i, node := range nodes {
		conn, err := net.Dial("tcp", node)
		if err != nil {
			return fmt.Errorf("Failed to connect to %s: %v", node, err)
		}

		defer conn.Close()
		connections[i] = conn
	}

	// Sending chunks to nodes in round-robin fashion 
	for i, chunk := range *chunks {
		nodeIndex := i % len(nodes) // here 2 chunks is been created ok 
		conn := connections[nodeIndex]

		err := sendChunkToNode(conn, &chunk)

		if err != nil {
			return fmt.Errorf("error sending chunk %s to node %s: %v", chunk.Id, nodes[nodeIndex], err)
		}
	}

	return nil
}

func sendChunkToNode(conn net.Conn, chunk *chunk.Chunk) error {
	_, err := conn.Write([]byte(chunk.Id + "\n"))
	if err != nil {
		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
	}
	_, err = conn.Write(chunk.Data)
	if err != nil {
		return fmt.Errorf("Failed to send chunk %s: %v", chunk.Id, err)
	}

	fmt.Printf("Sent chunk %s to node\n", chunk.Id)
	return nil
}
