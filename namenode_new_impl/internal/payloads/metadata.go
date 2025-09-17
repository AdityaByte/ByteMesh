package payloads

type MetaData struct {
	Filename      string              `json: "filename"`
	FileExtension string              `json: "fileExtension"`
	ActualSize    float64             `json: "actualFileSize"`
	Location      map[string]string   `json: "location"`
	AliveNodes    RegisteredDataNodes `json:"aliveNodes"`
}
