package payloads

type RegisteredDataNodes struct {
	Nodes []DataNode `json:"nodes"`
}

type DataNode struct {
	Name      string   `json:"name"`
	Port      uint16   `json:"port"`
	TimeStamp int64    `json:"time_stamp"`
}
