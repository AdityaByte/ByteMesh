package logger

import (
	"log"
	"os"
)

const (
	green = "\033[32m"
	red = "\033[31m"
	reset = "\033[0m"
)

var (
	InfoLogger = log.New(os.Stdout, green+"INFO: "+reset, log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, red+"ERROR: "+reset, log.Ldate|log.Ltime|log.Lshortfile)
)