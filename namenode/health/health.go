package health

import (
	"bufio"
	"fmt"
	"net"
)

func Health(conn net.Conn, reader *bufio.Reader) error {

	writer := bufio.NewWriter(conn)

	if _, err := writer.WriteString("OK\n"); err != nil {
		return fmt.Errorf("ERROR: Failed to send the OK response: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ERROR: Failed to Flush : %v", err)
	}

	return nil
}
