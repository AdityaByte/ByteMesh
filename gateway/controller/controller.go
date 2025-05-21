package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AdityaByte/bytemesh/client"
	"github.com/AdityaByte/bytemesh/coordinator"
	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
	"github.com/AdityaByte/bytemesh/utils"
)

// It handles the upload request ok.
func UploadController(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Request {Upload} Recieved")
	// Here we need to fetch the file which was sent by the client.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(20 << 30)
	if err != nil {
		http.Error(w, "Failed to parse the file "+err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error retrieving file from form "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Here i need to create the file at the local storage folder ok.
	if err := os.MkdirAll("../../storage", os.ModePerm); err != nil {
		http.Error(w, "Failed to create the folder "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now i need to save the file
	localFile, err := os.Create("../../storage/" + fileHeader.Filename)
	if err != nil {
		http.Error(w, "Failed to create the file "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(localFile, file)
	if err != nil {
		http.Error(w, "Failed to copy the content "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now i need to call the function ok.
	// middleware.CreateChunk(localFile)

	chunks, filename, filesize, err := middleware.CreateChunk(localFile)
	if err != nil {
		utils.RemoveFile(localFile.Name())
		logger.ErrorLogger.Fatalf("%v\n", err)
	}

	for i, chunk := range *chunks {
		logger.InfoLogger.Println("Iteration:", i, "Chunk ID:", chunk.Id, "Data Length:", len(chunk.Data))
	}

	if err := coordinator.SendChunks(chunks, filename, filesize); err != nil {
		utils.RemoveFile(localFile.Name())
		logger.ErrorLogger.Fatalf("%v\n", err)
	}

	if err := os.Remove(localFile.Name()); err != nil {
		logger.ErrorLogger.Println("Failed to remove the file from the storage folder:", err)
	}

	// If everthing goes correctly i need to return a response to the client.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully")

}

// Download controller for downloading the file.
func DownloadController(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Request {Download} Recieved")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If i gets the request i need to deserialize the details ok.
	// The client sends the request with the username and the filename he had to fetch.
	filename := r.URL.Query().Get("filename")
	_ = r.URL.Query().Get("user")

	filename = strings.TrimSpace(filename)

	//Once the data has been deserailized we need to fetch it ok.
	if err := client.Download(filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If it fetched out the file we just need to read it from the download folder and sent back to the client.

	filePath := "../../storage/" + filename
	// Here we need to fetch the file.
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed not found "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Here we have to set the headers
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to copy the data "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Here we need to print out the message.
	logger.InfoLogger.Println("File sent successfully")

}

// So what we have to sent to the client.
// Filename, Upload Time just these things.
type Data struct {
	Filename   string    `json: "filename"`
	UploadDate time.Time `json: "uploaddate"`
	Size       float64   `json: "size"`
}

func FetchController(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Request {FetchAll} Recieved")
	// If the request is of fetching the metadata like all files the user could have in that case.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Here what i need to do i need to call the client function before that i need to fetch the query parameter.
	user := r.URL.Query().Get("user")

	response, err := client.RetriveUserFiles(strings.TrimSpace(user))
	// Now i need to send the data to the client ok.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Here i need to send the data to the client.
	// Data could be in the json.
	var allData []Data

	for _, data := range *response {
		allData = append(allData, Data{
			Filename: data.Filename+"."+data.FileExtension,
			UploadDate: data.UploadDate,
			Size: data.ActualSize,
		})
	}

	// Once it has to be done we have to send this to the client.
	w.Header().Set("Content-Type", "application/json")

	// Here we need to send it using the json encoder.
	if err := json.NewEncoder(w).Encode(&allData); err != nil {
		http.Error(w, "Failed to encode the data "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.InfoLogger.Println("Data sent successfully")
}
