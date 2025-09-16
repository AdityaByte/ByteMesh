package payload

type GetRequest struct {
	Filename string `json:"file_name"`
	ChunkId  string `json:"chunk_id"`
}
