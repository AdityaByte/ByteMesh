package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AdityaByte/namenode/internal"
	"github.com/AdityaByte/namenode/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("No .env file found.")
	}
}

func main() {
	addr := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	server := internal.NewServer(addr)
	if err := server.Start(); err != nil {
		logger.ErrorLogger.Fatalf("Server failed to start %v", err)
		os.Exit(1)
	}
}
