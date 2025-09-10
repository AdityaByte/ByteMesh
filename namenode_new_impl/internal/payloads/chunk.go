package payloads

type Chunk struct {
	FileName string `json:"file_name"`
	FileId   string `json:"file_id`
	Data     []byte `json:"data"`
}
