package main

import (
	"bufio"
	"fmt"
	"net"
)

func Health(conn net.Conn, reader *bufio.Reader, writer *bufio.Writer) error {

	if _, err := writer.WriteString("OK\n"); err != nil {
		return fmt.Errorf("ERROR: Failed to send the OK response: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ERROR: Failed to Flush : %v", err)
	}

	return nil
}
