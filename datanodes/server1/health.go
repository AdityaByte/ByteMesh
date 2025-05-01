package main

import (
	"bufio"
	"fmt"
	"net"
)

func Health(conn net.Conn, reader *bufio.Reader, writer *bufio.Writer) error {

	// request, err := reader.ReadString('\n')

	// print("Request type: ", request)

	// if err != nil {
	// 	return fmt.Errorf("ERROR: Failed to Read the request: %v", err)
	// }

	// if strings.TrimSpace(request) != "HEALTH" {
	// 	return fmt.Errorf("ERROR: The request is not a health request: %s", request)
	// }

	if _, err := writer.WriteString("OK\n"); err != nil {
		return fmt.Errorf("ERROR: Failed to send the OK response: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ERROR: Failed to Flush : %v", err)
	}

	return nil
}
