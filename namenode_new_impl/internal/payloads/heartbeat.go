package payloads

type HeartBeat struct {
	NodeName  string `json:"node_name"`
	TimeStamp int64  `json:"timestamp"`
}
