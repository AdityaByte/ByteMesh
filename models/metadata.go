package models

import "time"

type MetaData struct {
	Owner         string
	Filename      string
	FileExtension string
	UploadDate    time.Time
	ActualSize    float64
	Location      map[string]string
}
