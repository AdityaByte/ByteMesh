package coordinator

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func createConnection(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to: %s, ERROR: %v", addr, err)
	}

	return conn, nil
}

func HealthChecker(addrs []string) []error {
	wg := sync.WaitGroup{} // Created an instance of a wait group.

	errorChan := make(chan error, len(addrs))

	for _, addr := range addrs {
		wg.Add(1) // Since we have to find the health of the 4 go routines.
		go func(addr string) {
			defer wg.Done()

			conn, err := createConnection(addr)
			if err != nil {
				errorChan <- err
				return
			}
			defer conn.Close()
			writer := bufio.NewWriter(conn)
			reader := bufio.NewReader(conn)
			writer.WriteString("HEALTH")
			if err := writer.Flush(); err != nil {
				errorChan <- fmt.Errorf("ERROR: Failed to flush to  %s", addr)
				return
			}
			response, err := reader.ReadString('\n')
			if err != nil {
				errorChan <- err
				return
			}

			if strings.TrimSpace(response) != "OK" {
				errorChan <- fmt.Errorf("Node %s sends a non-OK response: %s", addr, response)
				return
			} else {
				log.Printf("Health is OK: %s\n", addr)
			}
		}(addr)
	}

	wg.Wait()
	close(errorChan)

	var errors []error

	for err := range errorChan {
		errors = append(errors, err)
	}

	return errors
}

// I am just thinking about how we check the health ok this function runs in background and
