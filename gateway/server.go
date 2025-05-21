package main

import (
	"net/http"

	"github.com/AdityaByte/bytemesh/gateway/controller"
	"github.com/AdityaByte/bytemesh/gateway/middleware"
	"github.com/AdityaByte/bytemesh/logger"
)

// This is the gateway for the clients web interface
// Its the http server handling routes upload and download and getfiles.

func main() {
	http.HandleFunc("/upload", middleware.EnableCORS(controller.UploadController))
	http.HandleFunc("/download", middleware.EnableCORS(controller.DownloadController))
	http.HandleFunc("/fetchall", middleware.EnableCORS(controller.FetchController))

	logger.InfoLogger.Println("Starting server at :4444")
	if err := http.ListenAndServe(":4444", nil); err != nil {
		logger.ErrorLogger.Fatalf("ERROR: Failed to start the server, %v", err)
	}
}
