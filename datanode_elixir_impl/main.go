package main

import (
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:51315")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("POST\n"))

	// sending another data packet.
	conn.Write([]byte(`{"file_name": "foo", "file_id": "file_1", "data": "..."}` + "\n"))

	time.Sleep(3 * time.Second)
}
