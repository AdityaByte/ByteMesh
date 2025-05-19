package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/AdityaByte/bytemesh/auth/middleware"
	"github.com/AdityaByte/bytemesh/auth/service"
	"github.com/AdityaByte/bytemesh/datanodes/server1/logger"
)

func main() {

	http.HandleFunc("/signup", middleware.EnableCORS(service.Signup))
	http.HandleFunc("/login", middleware.EnableCORS(service.Login))
	http.HandleFunc("/validate", middleware.EnableCORS(service.ValidateToken))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not Allowed")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Health is OK")
	})

	logger.InfoLogger.Println("Starting Auth server")

	// When the Auth server starts we need to save the process id.
	if err := os.MkdirAll("../.auth/", os.ModePerm); err != nil {
		logger.ErrorLogger.Fatalf("ERROR: Failed to create the directory %v", err)
	}

	// Now we need to save the process id to the folder
	pidFile, err := os.Create("../.auth/.pid")
	if err != nil {
		logger.ErrorLogger.Fatalf("ERROR: Failed to create the file %v", err)
	}

	logger.InfoLogger.Printf("Process Id %d and its type %T", os.Getpid(), os.Getpid())

	// pidStr := string(os.Getpid()) //That's the wrong way to convert the pid to string
	pidStr := strconv.Itoa(os.Getpid())
	fmt.Println("String pid: ", pidStr)

	if _, err := pidFile.WriteString(pidStr); err != nil {
		logger.ErrorLogger.Fatalf("ERROR: Failed to write the process id %v", err)
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		if err := os.Remove("../.auth/.pid"); err != nil {
			logger.ErrorLogger.Fatalf("ERROR: Failed to remove the file %v", err)
		}
		logger.ErrorLogger.Fatalf("ERROR: Failed to start the Auth server at %s , %v\n", ":8080", err)
	}

}
