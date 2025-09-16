package handler

import (
	"net"

	"github.com/AdityaByte/namenode/internal/payloads"
)

func NodeRegistrationHandler(conn net.Conn, node *payloads.DataNode, datanodes *payloads.RegisteredDataNodes) {
	if datanodes == nil {
		return
	}
	for i, datanode := range datanodes.Nodes {
		if node.Name == datanode.Name {
			datanodes.Nodes[i] = *node
			return
		}
	}

	datanodes.Nodes = append(datanodes.Nodes,  *node)
}
