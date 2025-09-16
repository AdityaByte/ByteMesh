package payload

import "net"

// RegisteredDataNodes is mainly for deseralizing the gob object.
type RegisteredDataNodes struct {
	Nodes []DataNode `json:"nodes"`
}

type DataNode struct {
	Name      string   `json:"name"`
	Conn      net.Conn `json:"-"`
	Port      uint16   `json:"port"`
	TimeStamp int64    `json:"time_stamp"`
}
