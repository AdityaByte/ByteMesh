package model

type MetaData struct {
	Filename      string            `bson:"filename"`
	FileExtension string            `bson: "fileExtension"`
	ActualSize    float64           `bson:"actualFileSize"`
	Location      map[string]string `bson: "location"`
}
