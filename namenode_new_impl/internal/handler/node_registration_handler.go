package handler

import (
	"net"

	"github.com/AdityaByte/namenode/internal/payloads"
)

func NodeRegistrationHandler(conn net.Conn, node *payloads.DataNode, datanodes *payloads.RegisteredDataNodes) {
	for _, datanode := range datanodes.Nodes {
		if node.Name == datanode.Name {
			datanode = *node
		} else {
			datanodes.Nodes = append(datanodes.Nodes, *node)
		}
	}
}
