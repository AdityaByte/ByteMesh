package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/AdityaByte/bytemesh/client"
	"github.com/AdityaByte/bytemesh/coordinator"
	"github.com/AdityaByte/bytemesh/logger"
	"github.com/AdityaByte/bytemesh/middleware"
	"github.com/AdityaByte/bytemesh/utils"
	"github.com/joho/godotenv"
)

// The main purpose of the distributed file storage system is to allow users to
// upload files to the cloud and retrieve them when needed.
//
// One of the key advantages of a distributed file storage system is its fault tolerance,
// which is achieved through its architecture.
//
// 1. **NameNode**: Responsible for handling metadata, such as file locations and storage management.
// 2. **DataNode**: We currently have three DataNodes where file chunks are stored.
// 3. **Client**: Can perform two main operations—uploading files or downloading them.
// 4. **Middleware**: Manages file uploads, downloads, and metadata retrieval.
// 5. **Coordinator**: Sends requests to servers, fetches data, and passes it to the middleware for further processing.

// Firstly we have to check the health of the AuthServer
func isAuthServerRunning() bool {
	_, err := http.Get("http://localhost:8080/health")
	return err == nil
}

func startAuthServer() error {
	cmd := exec.Command("go", "run", "../auth/cmd/server.go")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("ERROR: Failed to start the Auth Server %v", err)
	}

	logger.InfoLogger.Println("Auth Server Started (Child process PID:", cmd.Process.Pid, ")")
	time.Sleep(2 * time.Second) // Waiting till the server is fully boots up.
	return nil
}

func stopAuthServer() error {
	data, err := os.ReadFile("../.auth/.pid")
	if err != nil {
		return fmt.Errorf("ERROR: Failed to read the file %v", err)
	}

	// Since we read the file we have to parse it to the integer
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("ERROR: Invalid PID found %v", err)
	}

	// Now we have to find the process and kill it.
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to find the process %v", err)
	}

	// Killing the process.
	if err := process.Kill(); err != nil {
		return fmt.Errorf("ERROR: Failed to kill the process %v", err)
	}

	// Removing the file.
	// Right now we are not removing the file.
	// if err := os.Remove("../../.auth/.pid"); err != nil {
	// 	return fmt.Errorf("ERROR: Failed to remove the process pid file %v", err)
	// }

	logger.InfoLogger.Println("Auth server stopped")
	return nil
}

// Initializing function which initialize all the environment variables by
// loading out the .env file.

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		logger.ErrorLogger.Fatalf("No .env file found")
	}
}

func defineCredentialsFlags(fs *flag.FlagSet) (*string, *string) {
	username := fs.String("username", "", "Username")
	password := fs.String("password", "", "Password")
	return username, password
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "signup":

		signupCmd := flag.NewFlagSet("signup", flag.ExitOnError)
		signupUsername, signupPassword := defineCredentialsFlags(signupCmd)
		signupCmd.Parse(os.Args[2:])

		logger.InfoLogger.Println(*signupUsername, *signupPassword)

		if err := client.SignUp(*signupUsername, *signupPassword); err != nil {
			logger.ErrorLogger.Fatalf("%v", err)
		}

	case "login":
		loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
		loginUsername, loginPassword := defineCredentialsFlags(loginCmd)

		loginCmd.Parse(os.Args[2:])

		if err := client.LogIn(*loginUsername, *loginPassword); err != nil {
			logger.ErrorLogger.Fatalf("%v", err)
		}

		logger.InfoLogger.Println("Login successful")

	case "auth":
		if len(os.Args) > 3 {
			logger.ErrorLogger.Fatalf("Invalid command")
		}

		switch strings.TrimSpace(os.Args[2]) {
		case "start":
			if !isAuthServerRunning() {
				err := startAuthServer()
				if err != nil {
					logger.ErrorLogger.Fatalf(err.Error())
				}
			} else {
				logger.InfoLogger.Println("Server is already running..")
			}
		case "stop":
			if isAuthServerRunning() {
				err := stopAuthServer()
				if err != nil {
					logger.ErrorLogger.Fatalf(err.Error())
				}
			} else {
				logger.InfoLogger.Println("Auth server is not running")
			}

		default:
			logger.ErrorLogger.Fatalf("Invalid command %s", os.Args[2])
		}

	default:

		const (
			version = "1.0.0"
			author  = "@AdityaByte"
		)

		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage : %s\n", "\nStarting the Authentication Server:\n \tgo run . auth start\nStopping the Authentication Server: \n\tgo run . auth stop\nSigning Up: \n\tgo run . signup -username <username> -password <password>\nLogging In: \n\tgo run . login -username <username> -password <password>\nUploading the file and Downloading the file: \n\tgo run . -upload \"FileLocation\" -download \"FileName\"")
			fmt.Fprintf(flag.CommandLine.Output(), "Version : %s\n", version)
			fmt.Fprintf(flag.CommandLine.Output(), "Author : %s\n", author)
			fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
			flag.PrintDefaults()
		}

		// Here we need to validate if the user is doing the upload or it is doing the download.
		fileLocation := flag.String("upload", "", "Location of the file to Upload")
		fileName := flag.String("download", "", "Name of the file to download")

		flag.Parse()

		if err := client.ValidateToken(); err != nil {
			logger.ErrorLogger.Fatalln(err)
		}

		if *fileLocation != "" {
			file, err := client.Upload(*fileLocation)
			if err != nil {
				utils.RemoveFile(file.Name()) // If something fails out we remove the file from the local folder.
				logger.ErrorLogger.Fatalf("%v", err)
			}

			chunks, filename, filesize, err := middleware.CreateChunk(file)
			if err != nil {
				utils.RemoveFile(file.Name())
				logger.ErrorLogger.Fatalf("%v", err)
			}

			for i, chunk := range *chunks {
				logger.InfoLogger.Println("Iteration:", i, "Chunk ID:", chunk.Id, "Data Length:", len(chunk.Data))
			}

			if err := coordinator.SendChunks(chunks, filename, filesize); err != nil {
				utils.RemoveFile(file.Name())
				logger.ErrorLogger.Fatalf("%v", err)
			}

			if err := os.Remove(fmt.Sprintf("../storage/%s", filename)); err != nil {
				logger.ErrorLogger.Println("Failed to remove the file from the storage folder:", err)
			}
		}

		if *fileName != "" {
			if err := client.Download(*fileName); err != nil {
				logger.ErrorLogger.Fatalf("%v", err)
			}
		}

	}
}
