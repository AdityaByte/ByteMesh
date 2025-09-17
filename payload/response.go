package payload

type Response struct {
	Type string      `json:"type"`
	Message interface{} `json:"data"`
}
