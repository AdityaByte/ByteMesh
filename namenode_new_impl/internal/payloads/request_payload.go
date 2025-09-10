package payloads

type GetRequest struct {
	FileName string `json:"file_name"`
	ChunkId  string `json:"chunk_id"`
}