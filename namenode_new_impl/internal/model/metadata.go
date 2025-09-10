package model

import "time"

type MetaData struct {
	Owner         string            `bson: "owner"`
	Filename      string            `bson: "filename"`
	FileExtension string            `bson: "fileExtension"`
	UploadDate    time.Time         `bson: "uploadDate"`
	ActualSize    float64           `bson: "actualFileSize"`
	Location      map[string]string `bson: "location"`
}
