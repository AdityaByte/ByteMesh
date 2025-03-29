package model

type MetaData struct {
	Filename      string            `bson:"filename"`
	FileExtension string            `bson: "fileExtension"`
	Location      map[string]string `bson: "location"`
}
