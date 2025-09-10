package payloads

import (
	"net"
)

type RegisteredDataNodes struct {
	Nodes []DataNode `json:"nodes"`
}

type DataNode struct {
	Name      string   `json:"name"`
	Conn      net.Conn `json:"-"`
	Port      uint16   `json:"port"`
	TimeStamp int64    `json:"time_stamp"`
}


// Health nodes information bhejenga ki kon konse nodes active hai thik hai post request me 